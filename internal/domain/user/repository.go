package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByCredentials(ctx context.Context, loginOrEmail string, validateFunc func(u User) error) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(u *User) error) error
}
