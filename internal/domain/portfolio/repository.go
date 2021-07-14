package portfolio

import (
	"context"

	"github.com/google/uuid"
)

type PortfolioRepository interface {
	CreatePortfolio(ctx context.Context, p *Portfolio) error
	ListPortfolios(ctx context.Context, userID uuid.UUID) ([]*Portfolio, error)
	GetPortfolio(ctx context.Context, id uuid.UUID) (*Portfolio, error)
	UpdatePortfolio(ctx context.Context, id uuid.UUID, updateFn func(p *Portfolio) (*Portfolio, error)) error
	DeletePortfolio(ctx context.Context, id uuid.UUID) error
}
