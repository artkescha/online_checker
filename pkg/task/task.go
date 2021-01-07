package task

import "time"

type Task struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"body"`
	Created       time.Time `json:"created"`
	TestsPaths    string    `json:"tests"`
}
