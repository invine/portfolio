package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/domain/user"
)

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) (*UserService, error) {
	if repo == nil {
		return nil, fmt.Errorf("missing repository")
	}

	u := &UserService{
		repo: repo,
	}
	return u, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, login, password, name string) error {
	if login == "" {
		login = email
	}

	u, err := user.NewUser(uuid.New(), email, login, password, name)
	if err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	if err := s.repo.CreateUser(ctx, u); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}

	return nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, loginOrEmail, password string) (*user.User, error) {
	u, err := s.repo.GetUserByLoginOrEmail(ctx, loginOrEmail)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if err := u.PasswordMatch(password); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return u, nil
}

func (s *UserService) ChangeUserEmail(ctx context.Context, id uuid.UUID, email string) error {
	err := s.repo.UpdateUser(ctx, id, func(u *user.User) error {
		return u.ChangeEmail(email)
	})

	if err != nil {
		return fmt.Errorf("change email for user id = %s: %w", id, err)
	}

	return nil
}

func (s *UserService) ChangeUserPassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	err := s.repo.UpdateUser(ctx, id, func(u *user.User) error {
		return u.ChangePassword(oldPassword, newPassword)
	})

	if err != nil {
		return fmt.Errorf("can't change password for user id = %s: %w", id, err)
	}

	return nil
}

func (s *UserService) ChangeUserName(ctx context.Context, id uuid.UUID, name string) error {
	err := s.repo.UpdateUser(ctx, id, func(u *user.User) error {
		return u.ChangeEmail(name)
	})

	if err != nil {
		return fmt.Errorf("change name for user id = %s: %w", id, err)
	}

	return nil
}
