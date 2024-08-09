package db

import (
	"database/sql"
	"time"

	"github.com/AleksBekker/project-api/models"
	"github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) (*Database, error) {
	db, err := sql.Open("mysql", dsn)
	return &Database{db: db}, err
}

func FromEnv(lookupEnv func(string) (string, bool)) (*Database, error) {
	return New(getDsn(lookupEnv))
}

func (database *Database) Close() error {
	return database.db.Close()
}

func (database *Database) GetNamesLimited(limit, offset int) ([]string, error) {
	rows, err := database.db.Query("SELECT name FROM Projects LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

func (database *Database) GetProjectsLimited(limit, offset int) ([]models.Project, error) {
	rows, err := database.db.Query(projectQueryStem+" LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	for rows.Next() {
		p := models.Project{}
		if err := scanProjectQuery(rows, &p); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (database *Database) GetProject(id int) (models.Project, error) {
	row := database.db.QueryRow(projectQueryStem+" WHERE id = ?", id)
	p := models.Project{}
	err := scanProjectQuery(row, &p)
	return p, err
}

func (database *Database) GetTags(projectID int) ([]models.Tag, error) {
	rows, err := database.db.Query(`
        SELECT t.tag_id, t.name 
        FROM Tags t RIGHT JOIN ProjectTags pt ON pt.tag_id = t.tag_id
        RIGHT JOIN Projects p ON p.project_id = pt.project_id
        WHERE p.project_id=?`, projectID)
	if err != nil {
		return nil, err
	}

	var tags []models.Tag
	for rows.Next() {
		tag := models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (database *Database) GetLinks(projectID int) ([]models.Link, error) {
    return []models.Link{{ID: projectID, URL: "some/url", Display: "neat link", Type: "experimental"}}, nil // TODO: implement me
}

const projectQueryStem = "SELECT project_id, name, description, start_date, end_date, status, priority FROM Projects"

type sqlScannable interface{ Scan(dest ...any) error }

func scanProjectQuery[R sqlScannable](row R, p *models.Project) error {
	var endDate *time.Time
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &endDate, &p.Status, &p.Priority)

	if endDate != nil {
		p.EndDate = endDate
	}

	return err
}

func getDsn(lookupEnv func(string) (string, bool)) string {
	cfg := mysql.NewConfig()
	cfg.ParseTime = true

	envs := []struct {
		key string
		ptr *string
	}{
		{"DB_USER", &cfg.User},
		{"DB_PASSWORD", &cfg.Passwd},
		{"DB_ADDRESS", &cfg.Addr},
		{"DB_DATABASE", &cfg.DBName},
	}

	for _, pair := range envs {
		val, _ := lookupEnv(pair.key)
		*pair.ptr = val
	}

	return cfg.FormatDSN()
}
