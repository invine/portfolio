package models

import (
	"fmt"
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
	if email == "" && login == "" {
		return nil, fmt.Errorf("both email and login can't be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password can't be empty")
	}

	u := User{}
	u.id = id
	u.email = email
	u.login = login
	u.password = password
	u.name = name

	return &u, nil
}

func NewUserWithEmail(email, password string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email can't be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password can't be empty")
	}
	return &User{id: uuid.New(), email: email, password: password}, nil
}

func NewUserWithLogin(login, password string) (*User, error) {
	if login == "" {
		return nil, fmt.Errorf("login can't be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password can't be empty")
	}
	return &User{id: uuid.New(), login: login, password: password}, nil
}

func (u *User) Update(email, login, password, name string) error {
	if email == "" && login == "" {
		return fmt.Errorf("both email and login can't be empty")
	}

	if password == "" {
		return fmt.Errorf("password can't be empty")
	}

	u.email = email
	u.login = login
	u.password = password
	u.name = name

	return nil
}
