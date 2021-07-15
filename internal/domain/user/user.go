package user

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password string
	login    string
	hash     string
}

func NewUser(id uuid.UUID, email, login, password, name string) (*User, error) {
	u := User{
		id: id,
	}

	if err := u.setPassword(password); err != nil {
		return nil, err
	}

	if err := u.ChangeEmail(email); err != nil {
		return nil, err
	}

	if err := u.setLogin(login); err != nil {
		return nil, err
	}

	if err := u.ChangeName(name); err != nil {
		return nil, err
	}

	return &u, nil
}

func NewUserFromDB(id uuid.UUID, email, login, passwordHash, name string) (*User, error) {
	u := User{
		id: id,
	}

	if err := u.setHash(passwordHash); err != nil {
		return nil, err
	}

	if err := u.ChangeEmail(email); err != nil {
		return nil, err
	}

	if err := u.setLogin(login); err != nil {
		return nil, err
	}

	if err := u.ChangeName(name); err != nil {
		return nil, err
	}

	return &u, nil
}

func (u *User) ChangeEmail(email string) error {
	if _, err := mail.ParseAddress(email); email != "" && err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	u.email = email
	return nil
}

func (u *User) setLogin(login string) error {
	if login == "" {
		return fmt.Errorf("login can't be empty")
	}

	u.login = login
	return nil
}

func (u *User) setPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password can't be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	u.password = password
	u.hash = string(hash)
	return nil
}

func (u *User) setHash(hash string) error {
	if hash == "" {
		return fmt.Errorf("hash can't be empty")
	}
	u.hash = string(hash)
	return nil
}

func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if err := u.PasswordMatch(oldPassword); err != nil {
		return err
	}
	if err := u.setPassword(newPassword); err != nil {
		return err
	}
	return nil
}

func (u *User) ChangeName(name string) error {
	u.name = name
	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Hash() string {
	return u.hash
}

func (u *User) Name() string {
	return u.name
}

func (u *User) PasswordMatch(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.hash), []byte(password)); err != nil {
		return err
	}
	return nil
}
