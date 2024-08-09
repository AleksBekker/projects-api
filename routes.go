package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AleksBekker/project-api/database"
)

func handleGetProjects(db *db.Database, l *log.Logger) http.Handler {
	const (
		maxLimit      = 100
		defaultLimit  = 20
		defaultOffset = 0
	)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			limit, err := forceQuery(r, "limit", defaultLimit, strconv.Atoi, func(x int) bool { return x > 0 && x <= maxLimit })
			if writeQueryError(err, w, l) {
				return
			}

			offset, err := forceQuery(r, "offset", defaultOffset, strconv.Atoi, func(x int) bool { return x >= 0 })
			if writeQueryError(err, w, l) {
				return
			}

			projects, err := db.GetProjectsLimited(limit, offset)
			if err != nil {
				l.Printf("internal server error: %+s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for idx := range projects {
                project := &projects[idx]
				if tags, err := db.GetTags(project.ID); err != nil {
                    l.Printf("internal server error: %+s\n", err)
                    w.WriteHeader(http.StatusInternalServerError)
                    return
				} else {
					project.Tags = tags
				}

                if links, err := db.GetLinks(project.ID); err != nil {
                    l.Printf("internal server error: %+s\n", err)
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                } else {
                    project.Links = links
                }
			}

			err = encode(w, http.StatusOK, projects)

			if err != nil {
				l.Printf("internal server error: %+s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		},
	)
}

type transformer[T any] func(string) (T, error)
type validator[T any] func(T) bool

func forceQuery[T any](r *http.Request, key string, default_ T, transformer transformer[T], validator validator[T]) (T, error) {
	boxed, err := query(r, key, default_, transformer)

	if err == nil && !validator(boxed) {
		return boxed, &invalidQueryError{Key: key}
	}

	return boxed, err
}

func query[T any](r *http.Request, key string, default_ T, transform transformer[T]) (T, error) {
	var value string

	if value = r.URL.Query().Get(key); value == "" {
		return default_, nil
	}

	return transform(value)
}

func atoi(str string) (*int, error) {
	val, err := strconv.Atoi(str)
	return &val, err
}

type invalidQueryError struct {
	Key string
}

// writeQueryError writes an error description to a response and returns true if err is an error,
// returns false otherwise
func writeQueryError(err error, w http.ResponseWriter, l *log.Logger) bool {
	var qerr *invalidQueryError
	if errors.As(err, &qerr) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, qerr.UserError())
	} else if err != nil {
		l.Printf("internal server error: %+s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err != nil
}

func (err *invalidQueryError) Error() string {
	return err.UserError()
}

func (err *invalidQueryError) UserError() string {
	return "invalid query " + err.Key
}
