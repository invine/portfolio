package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/domain/portfolio"
)

type CreatePortfolio struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

type CreatePortfolioHandler struct {
	repo portfolio.PortfolioRepository
}

func NewCreatePortfolioHandler(repo portfolio.PortfolioRepository) (*CreatePortfolioHandler, error) {
	if repo == nil {
		return nil, fmt.Errorf("repo can't be empty")
	}

	return &CreatePortfolioHandler{repo: repo}, nil
}

func (h CreatePortfolioHandler) Handle(ctx context.Context, cmd CreatePortfolio) error {
	p, err := portfolio.NewPortfolio(cmd.ID, cmd.UserID, cmd.Name, nil)
	if err != nil {
		return fmt.Errorf("can't create portfolio %s: %w", cmd.Name, err)
	}

	if err := h.repo.CreatePortfolio(ctx, p); err != nil {
		return fmt.Errorf("can't create portfolio %s: %w", cmd.Name, err)
	}

	return nil
}
