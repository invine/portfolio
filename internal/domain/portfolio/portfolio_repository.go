package portfolio

import (
	"context"

	"github.com/google/uuid"
)

type PortfolioRepository interface {
	CreatePortfolio(ctx context.Context, p *Portfolio) error
	GetAllPortfolios(ctx context.Context, userID uuid.UUID) ([]*Portfolio, error)
	GetPortfolio(ctx context.Context, userID, id uuid.UUID) (*Portfolio, error)
	UpdatePortfolio(ctx context.Context, userID, id uuid.UUID, updateFn func(p *Portfolio) error) error
	DeletePortfolio(ctx context.Context, userID, id uuid.UUID) error
}
