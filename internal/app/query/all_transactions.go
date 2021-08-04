package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/domain/portfolio"
)

type AllTransactionsHandler struct {
	readModel AllTransactionsReadModel
}

type AllTransactionsReadModel interface {
	GetAllTransactions(ctx context.Context, userID, PortfolioId uuid.UUID) ([]*portfolio.Transaction, error)
}

type AllTransactions struct {
	UserID      uuid.UUID
	PortfolioID uuid.UUID
}

func NewAllTransactionsHandler(readModel AllTransactionsReadModel) (*AllTransactionsHandler, error) {
	if readModel == nil {
		return nil, fmt.Errorf("empty readModel")
	}
	return &AllTransactionsHandler{readModel: readModel}, nil

}

func (h AllTransactionsHandler) Handle(ctx context.Context, query AllTransactions) ([]*portfolio.Transaction, error) {
	return h.readModel.GetAllTransactions(ctx, query.UserID, query.PortfolioID)
}
