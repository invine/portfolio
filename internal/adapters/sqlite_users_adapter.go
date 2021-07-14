package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/domain/user"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUsersRepository struct {
	db *sql.DB
}

type rowQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func NewSQLiteUsersRepository(db *sql.DB) (*SQLiteUsersRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database required")
	}

	sqlStmt := `
	create table users (id text not null primary key, email text, login text not null, password text not null, name text, unique(email, login));
	`
	_, err := db.Exec(sqlStmt)
	if err != nil && err.Error() != "table users already exists" {
		return nil, fmt.Errorf("can't insert table: %w", err)
	}

	r := &SQLiteUsersRepository{db: db}
	return r, nil
}

func (r *SQLiteUsersRepository) CreateUser(ctx context.Context, u *user.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	defer tx.Rollback()

	// TODO: evaluate if it's really necessary to have all this getters
	sql := "insert into users (id, email, login, password, name) values ($1, $2, $3, $4, $5)"
	if _, err := tx.ExecContext(ctx, sql, u.ID(), u.Email(), u.Login(), u.Password(), u.Name()); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}

func (r *SQLiteUsersRepository) GetUserByCredentials(ctx context.Context, login, passowrd string) (*user.User, error) {
	return r.getUserByCredentials(ctx, r.db, login, passowrd, false)
}

func (r *SQLiteUsersRepository) UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(u *user.User) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}
	defer tx.Rollback()

	u, err := r.getUserByID(ctx, tx, id, true)
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	if err := updateFn(u); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	sql := "update users set email = $2, login = $3, password = $4, name = $5 where id = $1"
	if _, err := tx.ExecContext(ctx, sql, u.ID(), u.Email(), u.Login(), u.Password(), u.Name()); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	return nil
}

func (r *SQLiteUsersRepository) getUserByID(ctx context.Context, db rowQuerier, id uuid.UUID, forUpdate bool) (*user.User, error) {
	sql := "select email, login, password, name from users where id = $1"
	if forUpdate {
		sql += " for update"
	}
	row := db.QueryRowContext(ctx, sql, id)

	var email, login, password, name string
	if err := row.Scan(&email, &login, &password, &name); err != nil {
		return nil, fmt.Errorf("user not foud: %w", err)
	}

	u, err := user.NewUser(id, email, login, password, name)
	if err != nil {
		return nil, fmt.Errorf("wrong user parameters: %w", err)
	}

	return u, nil
}

func (r *SQLiteUsersRepository) getUserByCredentials(ctx context.Context, db rowQuerier, login, password string, forUpdate bool) (*user.User, error) {
	sql := "select id, email, login, name from users where (login = $1 or email = $1) and password = $2"
	if forUpdate {
		sql += " for update"
	}
	row := db.QueryRowContext(ctx, sql, login, password)

	var id uuid.UUID
	var email, loginDB, name string
	if err := row.Scan(&id, &email, &loginDB, &name); err != nil {
		return nil, fmt.Errorf("user not foud: %w", err)
	}

	u, err := user.NewUser(id, email, loginDB, password, name)
	if err != nil {
		return nil, fmt.Errorf("wrong user parameters: %w", err)
	}

	return u, nil
}
