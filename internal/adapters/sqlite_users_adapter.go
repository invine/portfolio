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

type userModel struct {
	ID           uuid.UUID
	Email        string
	Login        string
	PasswordHash string
	Name         string
}

func NewSQLiteUsersRepository(db *sql.DB) (*SQLiteUsersRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database required")
	}

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS users
        (
            id text not null primary key,
            email text,
            login text not null,
            password text not null,
            name text,
            unique(email, login)
        );
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
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

	um := userToUserModel(u)

	// TODO: evaluate if it's really necessary to have all this getters
	sql := "insert into users (id, email, login, password, name) values ($1, $2, $3, $4, $5)"
	if _, err := tx.ExecContext(ctx, sql, um.ID, um.Email, um.Login, um.PasswordHash, um.Name); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}

func (r *SQLiteUsersRepository) GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (*user.User, error) {
	um, err := r.getUserByLoginOrEmail(ctx, r.db, loginOrEmail, false)
	if err != nil {
		return nil, fmt.Errorf("wrong user parameters: %w", err)
	}

	// TODO evaluate if 's ok to have separate user constructor for DB or it should be decided in use case
	u, err := user.NewUserFromDB(um.ID, um.Email, um.Login, um.PasswordHash, um.Name)
	if err != nil {
		return nil, fmt.Errorf("wrong user parameters: %w", err)
	}

	return u, nil
}

func (r *SQLiteUsersRepository) UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(u *user.User) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}
	defer tx.Rollback()

	um, err := r.getUserByID(ctx, tx, id, true)
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	// TODO evaluate if 's ok to have separate user constructor for DB or it should be decided in use case
	u, err := user.NewUserFromDB(um.ID, um.Email, um.Login, um.PasswordHash, um.Name)
	if err != nil {
		return fmt.Errorf("wrong user parameters: %w", err)
	}

	if err := updateFn(u); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}
	um = userToUserModel(u)

	sql := "update users set email = $2, login = $3, password = $4, name = $5 where id = $1"
	if _, err := tx.ExecContext(ctx, sql, um.ID, um.Email, um.Login, um.PasswordHash, um.Name); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	return nil
}

func (r *SQLiteUsersRepository) getUserByID(ctx context.Context, db rowQuerier, id uuid.UUID, forUpdate bool) (*userModel, error) {
	sql := "select email, login, password, name from users where id = $1"
	if forUpdate {
		sql += " for update"
	}
	row := db.QueryRowContext(ctx, sql, id)

	var email, login, passwordHash, name string
	if err := row.Scan(&email, &login, &passwordHash, &name); err != nil {
		return nil, fmt.Errorf("user not foud: %w", err)
	}

	um := &userModel{
		ID:           id,
		Email:        email,
		Login:        login,
		PasswordHash: passwordHash,
		Name:         name,
	}

	return um, nil
}

func (r *SQLiteUsersRepository) getUserByLoginOrEmail(ctx context.Context, db rowQuerier, loginOrEmail string, forUpdate bool) (*userModel, error) {
	sql := "select id, email, login, password, name from users where (email = $1 or login = $1)"
	if forUpdate {
		sql += " for update"
	}
	row := db.QueryRowContext(ctx, sql, loginOrEmail)

	var id uuid.UUID
	var email, login, passwordHash, name string
	if err := row.Scan(&id, &email, &login, &passwordHash, &name); err != nil {
		return nil, fmt.Errorf("user not foud: %w", err)
	}

	u := &userModel{
		ID:           id,
		Email:        email,
		Login:        login,
		PasswordHash: passwordHash,
		Name:         name,
	}

	return u, nil
}

// TODO: evaluate if it's really necessary to have all this getters
func userToUserModel(u *user.User) *userModel {
	return &userModel{
		ID:           u.ID(),
		Email:        u.Email(),
		Login:        u.Login(),
		PasswordHash: u.Hash(),
		Name:         u.Name(),
	}
}
