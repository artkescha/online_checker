package repository

import (
	"database/sql"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/user"
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

	//TODO fix hardCode RoleID - user
	roleID := 2
	row := repo.db.QueryRow("INSERT INTO users(name, password, role_id) VALUES($1,$2, $3) returning id;", login, password, roleID)
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
		RoleID:   roleID,
	}

	return user, nil
}

func (repo *Repo) GetUserByLogin(login string) (*user.User, error) {
	user := &user.User{}
	if err := repo.db.Ping(); err != nil {
		return nil, err
	}
	row := repo.db.QueryRow("SELECT * FROM users where name=$1", login)

	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("get user by login %s failed, reason: %s", login, err)
	}
	return user, nil
}
