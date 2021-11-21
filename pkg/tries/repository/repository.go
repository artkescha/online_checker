package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/tries"
)

//go:generate mockgen -destination=./repository_mock.go -package=repository . TriesRepo

type TriesRepo interface {
	Insert(ctx context.Context, try try.Try) error
	List(ctx context.Context, limit, offset uint32) ([]try.Try, error)
	ListByUser(ctx context.Context, userID uint64, limit, offset uint32) ([]try.Try, error)
	GetByID(ctx context.Context, id uint64) (try.Try, error)
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

func (repo Repo) List(ctx context.Context, limit, offset uint32) ([]try.Try, error) {
	rows, err := repo.db.Query("SELECT * FROM attempts ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return []try.Try{}, fmt.Errorf("query tries list failed: %s", err)
	}
	tries := make([]try.Try, 0)
	for rows.Next() {
		try := try.Try{}
		err = rows.Scan(&try.ID, &try.UserID, &try.Solution, &try.Status,
			&try.Description, &try.Created, &try.TaskID, &try.LanguageID)
		if err != nil {
			return nil, fmt.Errorf("read rows tries list failed: %s", err)
		}
		tries = append(tries, try)
	}
	return tries, nil
}

func (repo Repo) ListByUser(ctx context.Context, userID uint64, limit, offset uint32) ([]try.Try, error) {
	rows, err := repo.db.Query("SELECT * FROM attempts WHERE user_id = $1 ORDER BY attempts.created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset)
	if err != nil {
		return []try.Try{}, fmt.Errorf("query tries list failed: %s", err)
	}
	tries := make([]try.Try, 0)

	for rows.Next() {
		try := try.Try{}
		err = rows.Scan(&try.ID, &try.UserID, &try.Solution, &try.Status,
			&try.Description, &try.Created, &try.TaskID, &try.LanguageID)
		if err != nil {
			return nil, fmt.Errorf("read rows tries list failed: %s", err)
		}
		tries = append(tries, try)
	}
	return tries, nil
}

func (repo Repo) GetByID(ctx context.Context, id uint64) (try.Try, error) {
	rows, err := repo.db.Query("SELECT id, user_id, solution, status, task_id, language_id FROM attempts WHERE id = $1", id)
	if err != nil {
		return try.Try{}, fmt.Errorf("query try with id %d failed: %s", id, err)
	}
	try_ := try.Try{}
	if !rows.Next() {
		return try.Try{}, fmt.Errorf("query try with id %d response is empty", id)
	}

	err = rows.Scan(&try_.ID, &try_.UserID, &try_.Solution, &try_.Status, &try_.TaskID, &try_.LanguageID)
	if err != nil {
		return try.Try{}, fmt.Errorf("scanning response failed: %s", err)
	}
	return try_, nil

}
