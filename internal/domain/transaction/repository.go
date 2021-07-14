package transaction

import (
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, t *Transaction) error
	ListTransactions(ctx context.Context, userID, portfolioID uuid.UUID) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, updateFn func(t *Transaction) (*Transaction, error)) error
}
