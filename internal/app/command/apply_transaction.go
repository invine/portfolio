package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/domain/portfolio"
)

type ApplyTransaction struct {
	UserID      uuid.UUID
	PortfolioID uuid.UUID
	Transaction *portfolio.Transaction
}

type ApplyTransactionHandler struct {
	repo portfolio.PortfolioRepository
}

func NewApplyTransactionHandler(repo portfolio.PortfolioRepository) (*ApplyTransactionHandler, error) {
	if repo == nil {
		return nil, fmt.Errorf("portfolio repo can't be empty")
	}
	return &ApplyTransactionHandler{repo: repo}, nil
}

func (h ApplyTransactionHandler) Handle(ctx context.Context, cmd ApplyTransaction) error {
	return h.repo.UpdatePortfolio(
		ctx,
		cmd.UserID,
		cmd.PortfolioID,
		func(p *portfolio.Portfolio) error {
			if err := p.ApplyTransaction(cmd.Transaction); err != nil {
				return fmt.Errorf("can't apply transaction to portfolio %s: %w", cmd.PortfolioID.String(), err)
			}
			return nil
		})
}
