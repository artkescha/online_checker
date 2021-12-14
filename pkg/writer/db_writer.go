package writer

import (
	"database/sql"
	try "github.com/artkescha/checker/online_checker/pkg/tries"
)

type DBWriter struct {
	db *sql.DB
}

func NewDBWriter(db *sql.DB) *DBWriter {
	return &DBWriter{
		db: db,
	}
}

func (w DBWriter) Write(one try.Try) (int64, error) {
	var lastInsertId int64
	row := w.db.QueryRow("INSERT INTO attempts(user_id, solution, status, description, task_id, language_id) "+
		"VALUES($1, $2, $3, $4, $5, $6) returning id;",
		one.UserID, one.Solution, one.Status, one.Description, one.TaskID, one.LanguageID)
	err := row.Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}
