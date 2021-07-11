package repos

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(DBName string) (*UserRepo, error) {
	// os.Remove(DBName)

	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	sqlStmt := `
	create table users (id text not null primary key, email text, login text, password text, name text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil && err.Error() != "table users already exists" {
		return nil, fmt.Errorf("can't insert table: %w", err)
	}

	return &UserRepo{db: db}, nil
}

func (p *UserRepo) Close() error {
	return p.db.Close()

}

func (ur *UserRepo) Authenticate(login, password string) (uuid.UUID, error) {
	rows := ur.db.QueryRow("select id from users where (login = $1 OR email =$1) and password = $2", login, password)
	var id uuid.UUID
	err := rows.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("not authorized: %w", err)
	}

	return id, nil
}

func (ur *UserRepo) CreateUser(email, login, password, name string) error {
	id := uuid.New()
	_, err := models.NewUser(id, email, login, password, name)

	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	tx, err := ur.db.Begin()
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	stmt, err := tx.Prepare("insert into users(id, email, login, password, name) values(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, email, login, password, name)
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}
