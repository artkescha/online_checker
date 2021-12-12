package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/task"
)

//go:generate mockgen -destination=./repository_mock.go -package=repository . TaskRepo

type TaskRepo interface {
	Insert(ctx context.Context, task task.Task) (*task.Task, error)
	List(ctx context.Context, limit, offset uint32, sortField string) ([]task.Task, error)
	ListByUser(ctx context.Context, userID uint64, limit, offset uint32, sortField string) ([]task.Task, error)
	GetByID(ctx context.Context, id int) (*task.Task, error)
	Update(ctx context.Context, task *task.Task) (*task.Task, error)
	Delete(ctx context.Context, id uint32) (bool, error)
}

type Repo struct {
	db *sql.DB
}

func NewTasksRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo Repo) Insert(ctx context.Context, task task.Task) (*task.Task, error) {
	var lastInsertId int64
	//TODO!!!!!!!!!
	var isExist = fmt.Errorf(`pq: duplicate key value violates unique constraint "constraint_name"`)
	row := repo.db.QueryRow("INSERT INTO tasks(title, description, test_path) VALUES($1,$2, $3) returning id;", task.Title, task.Description, "C:/")
	err := row.Scan(&lastInsertId)
	if err != nil {
		if err.Error() == isExist.Error() {
			return nil, fmt.Errorf("task with id: %d already exists", lastInsertId)
		}
		return nil, err
	}
	task.ID = int(lastInsertId)
	return &task, nil
}

func (repo Repo) List(ctx context.Context, limit, offset uint32, sortField string) ([]task.Task, error) {
	rows, err := repo.db.Query("SELECT * FROM tasks ORDER BY $1 LIMIT $2 OFFSET $3", sortField, limit, offset)
	if err != nil {
		return []task.Task{}, fmt.Errorf("query tasks list failed: %s", err)
	}
	tasks := make([]task.Task, 0)

	for rows.Next() {
		task := task.Task{}
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Created, &task.TestsPath)
		if err != nil {
			return nil, fmt.Errorf("read rows tasks list failed: %s", err)
		}
		task.Created = task.Created.Local()
		tasks = append(tasks, task)
	}
	return tasks, nil
}

//TODO implement later
func (repo Repo) ListByUser(ctx context.Context, userID uint64, limit, offset uint32, sortField string) ([]task.Task, error) {
	panic("implement list by user")
	//return nil, nil
}

func (repo Repo) GetByID(ctx context.Context, id int) (*task.Task, error) {
	task := task.Task{}
	//TODO!!!!!!!
	row := repo.db.QueryRow("SELECT id, title, description FROM tasks where id=$1", id)
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	return &task, nil
}

func (repo Repo) Update(ctx context.Context, task *task.Task) (*task.Task, error) {
	repo.db.QueryRow("UPDATE tasks SET title = $1, description = $2 where id=$3", task.Title, task.Description, task.ID)
	//TODO!!!!!!!!!!!!!!!
	return nil, nil
}

func (repo Repo) Delete(ctx context.Context, id uint32) (bool, error) {
	//TODO!!!!!!!!!!!!!!!!!
	repo.db.QueryRow("DELETE FROM tasks where id = $1", id)
	return true, nil
}
