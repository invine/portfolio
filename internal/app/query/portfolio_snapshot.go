package query

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/invine/portfolio/internal/domain/portfolio"
)

type PortfolioHandler struct {
	readModel PortfolioReadModel
}

type PortfolioReadModel interface {
	GetPortfolio(ctx context.Context, userID, id uuid.UUID) (*portfolio.Portfolio, error)
}

type Portfolio struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Date   time.Time
}

func NewPortfolioHandler(readModel PortfolioReadModel) (*PortfolioHandler, error) {
	if readModel == nil {
		return nil, fmt.Errorf("empty readModel")
	}
	return &PortfolioHandler{readModel: readModel}, nil

}

func (h PortfolioHandler) Handle(ctx context.Context, query Portfolio) (*portfolio.Snapshot, error) {
	p, err := h.readModel.GetPortfolio(ctx, query.UserID, query.ID)
	if err != nil {
		return nil, fmt.Errorf("can't get portfolio %s: %w", query.ID.String(), err)
	}
	return p.Snapshot(query.Date), nil
}
