package user

import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password string
	login    string
}

func NewUser(id uuid.UUID, email, login, password, name string) (*User, error) {
	u := User{
		id: id,
	}

	if err := u.ChangePassword(password); err != nil {
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
	if _, err := mail.ParseAddress(u.email); u.email != "" && err != nil {
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

func (u *User) ChangePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password can't be empty")
	}

	u.password = password
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

func (u *User) Password() string {
	return u.password
}

func (u *User) Name() string {
	return u.name
}
