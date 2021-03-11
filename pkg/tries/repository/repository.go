package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/tries"
)

type TriesRepo interface {
	Insert(ctx context.Context, try try.Try) error
	List(ctx context.Context, limit, offset uint32, sortField string) ([]try.Try, error)
	ListByUser(ctx context.Context, userID uint64, limit, offset uint32, sortField string) ([]try.Try, error)
}

type Repo struct {
	db *sql.DB
}

func NewTriesRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo Repo) Insert(ctx context.Context, try try.Try) error {
	var lastInsertId int64
	var isExist = fmt.Errorf(`pq: duplicate key value violates unique constraint "constraint_name"`)
	row := repo.db.QueryRow("INSERT INTO tasks(user_id, solution, status, description, task_id, language_id) VALUES($1, $2, $3, $4, $5, $6) returning id;",
		try.UserID, try.Solution, try.Status, try.TaskID, try.LanguageID)
	err := row.Scan(&lastInsertId)
	if err != nil {
		if err.Error() == isExist.Error() {
			return fmt.Errorf("try with id: %d already exists", lastInsertId)
		}
		return err
	}
	return nil
}

func (repo Repo) List(ctx context.Context, limit, offset uint32, sortField string) ([]try.Try, error) {
	rows, err := repo.db.Query("SELECT * FROM tries ORDER BY $1 DESC LIMIT $2 OFFSET $3", sortField, limit, offset)
	if err != nil {
		return []try.Try{}, fmt.Errorf("query tries list failed: %s", err)
	}
	tries := make([]try.Try, 0)
	for rows.Next() {
		try := try.Try{}
		err = rows.Scan(&try)
		if err != nil {
			return nil, fmt.Errorf("read rows tries list failed: %s", err)
		}
		tries = append(tries, try)
	}
	return tries, nil
}

func (repo Repo) ListByUser(ctx context.Context, userID uint64, limit, offset uint32, sortField string) ([]try.Try, error) {
	rows, err := repo.db.Query("SELECT * FROM tries WHERE user_id = $1 ORDER BY $2 DESC LIMIT $3 OFFSET $4",
		userID, sortField, limit, offset)
	if err != nil {
		return []try.Try{}, fmt.Errorf("query tries list failed: %s", err)
	}
	tries := make([]try.Try, 0)

	for rows.Next() {
		try := try.Try{}
		err = rows.Scan(&try)
		if err != nil {
			return nil, fmt.Errorf("read rows tries list failed: %s", err)
		}
		tries = append(tries, try)
	}
	return tries, nil
}
