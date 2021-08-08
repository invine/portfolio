package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/portfolio/internal/domain/portfolio"
)

type DeletePortfolio struct {
	UserID      uuid.UUID
	PortfolioID uuid.UUID
}

type DeletePortfolioHandler struct {
	repo portfolio.PortfolioRepository
}

func NewDeletePortfolioHandler(repo portfolio.PortfolioRepository) (*DeletePortfolioHandler, error) {
	if repo == nil {
		return nil, fmt.Errorf("portfolio repo can't be empty")
	}
	return &DeletePortfolioHandler{repo: repo}, nil
}

func (h DeletePortfolioHandler) Handle(ctx context.Context, cmd DeletePortfolio) error {
	return h.repo.DeletePortfolio(ctx, cmd.UserID, cmd.PortfolioID)
}
