package models

import "time"

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Tags        []Tag      `json:"tags,omitempty"`
	Links       []Link     `json:"links,omitempty"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Display string `json:"display"`
	Type    string `json:"type"`
}
