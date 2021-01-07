package repository

import (
	"database/sql"
	"fmt"
	"gitlab.com/artkescha/grader/online_checker/pkg/user"
)

//go:generate mockgen -destination=./repository_mock.go -package=repository . UserRepo

type UserRepo interface {
	GetUserByLogin(username string) (*user.User, error)
	Insert(login string, password string) (*user.User, error)
}

type Repo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo *Repo) Insert(login string, password string) (*user.User, error) {
	var lastInsertId int64
	var isExist = fmt.Errorf(`pq: duplicate key value violates unique constraint "constraint_name"`)

	row := repo.db.QueryRow("INSERT INTO users(name, password, role_id) VALUES($1,$2, $3) returning id;", login, password, 2)
	err := row.Scan(&lastInsertId)
	if err != nil {
		if err.Error() == isExist.Error() {
			return nil, fmt.Errorf("already exists")
		}
		return nil, err
	}
	user := &user.User{
		ID:       lastInsertId,
		Name:     login,
		Password: password,
	}

	return user, nil
}

func (repo Repo) GetUserByLogin(username string) (*user.User, error) {
	user := &user.User{}
	if err := repo.db.Ping(); err != nil {
		return nil, err
	}
	row := repo.db.QueryRow("SELECT * FROM users where name=$1", username)

	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
