package portfolio

import (
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, t *Transaction) error
	ListTransactions(ctx context.Context, userID, portfolioID uuid.UUID) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, updateFn func(t *Transaction) error) error
	DeleteTransactions(ctx context.Context, ids []uuid.UUID) error
}
