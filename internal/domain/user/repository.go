package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByCredentials(ctx context.Context, login, passowrd string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(u *User) error) error
}
