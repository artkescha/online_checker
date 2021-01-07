package repository

import (
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/artkescha/grader/online_checker/pkg/user"
)

func TestRepo_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		db       *sql.DB
		login    string
		password string
		wantRows *sqlmock.Rows
		want     *user.User
		wantErr  bool
	}{
		{
			name:     "test without error",
			db:       db,
			login:    "artem",
			password: "12345678",
			wantRows: sqlmock.NewRows([]string{"id"}).AddRow(1),
			want: &user.User{
				ID:       1,
				Name:     "artem",
				Password: "12345678",
			},
			wantErr: false,
		},

		{
			name:     "test with error",
			db:       db,
			login:    "artem",
			password: "12345678",
			wantRows: sqlmock.NewRows([]string{}),
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		const query = `INSERT INTO`
		mock.ExpectQuery(query).
			WithArgs(tt.login, tt.password).
			WillReturnRows(tt.wantRows)
		repo := &Repo{
			db: tt.db,
		}

		got, err := repo.Insert(tt.login, tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("Repo.Insert() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Repo.Insert() = %v, want %v", got, tt.want)
		}
	}
}

func TestRepo_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name     string
		db       *sql.DB
		login    string
		wantRows *sqlmock.Rows
		want     *user.User
		wantErr  bool
	}{
		{
			name:     "test without error",
			db:       db,
			login:    "artem",
			wantRows: sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(1, "artem", user.GetMD5Password("12345678")),
			want: &user.User{
				ID:       1,
				Name:     "artem",
				Password: user.GetMD5Password("12345678"),
			},
			wantErr: false,
		},
		{
			name:     "test with error",
			db:       db,
			login:    "artem",
			wantRows: sqlmock.NewRows([]string{}),
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		const query = `SELECT * FROM users where name`
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(tt.login).
			WillReturnRows(tt.wantRows)

		repo := &Repo{
			db: tt.db,
		}

		got, err := repo.GetUserByLogin(tt.login)
		if (err != nil) != tt.wantErr {
			t.Errorf("Repo.GetUserByLogin(%s) error = %v, wantErr %v", tt.login, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Repo.GetUserByLogin(%s) = %v, want %v", tt.login, got, tt.want)
		}
	}
}
