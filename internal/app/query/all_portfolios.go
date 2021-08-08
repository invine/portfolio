package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/portfolio/internal/domain/portfolio"
)

type AllPortfoliosHandler struct {
	readModel AllPortfoliosReadModel
}

type AllPortfoliosReadModel interface {
	GetAllPortfolios(ctx context.Context, userID uuid.UUID) ([]*portfolio.Portfolio, error)
}

type AllPortfolios struct {
	UserID uuid.UUID
}

func NewAllPortfoliosHandler(readModel AllPortfoliosReadModel) (*AllPortfoliosHandler, error) {
	if readModel == nil {
		return nil, fmt.Errorf("empty readModel")
	}
	return &AllPortfoliosHandler{readModel: readModel}, nil

}

func (h AllPortfoliosHandler) Handle(ctx context.Context, query AllPortfolios) ([]*portfolio.Portfolio, error) {
	return h.readModel.GetAllPortfolios(ctx, query.UserID)
}
