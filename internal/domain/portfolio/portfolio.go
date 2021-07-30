package portfolio

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	asset    string
	quantity int
}

type Portfolio struct {
	id           uuid.UUID
	userID       uuid.UUID
	name         string
	transactions []*Transaction
}

func NewPortfolio(id, userID uuid.UUID, name string) (*Portfolio, error) {
	// TODO add validation that name is not nil
	p := &Portfolio{
		id:     id,
		userID: userID,
		name:   name,
	}

	return p, nil
}

func NewPortfolioWithTransactions(id, userID uuid.UUID, name string, transactions []*Transaction) (*Portfolio, error) {
	return &Portfolio{
		id:           id,
		userID:       userID,
		name:         name,
		transactions: transactions,
	}, nil
}

func (p *Portfolio) Snapshot(date time.Time) ([]Asset, float64) {
	assets := []Asset{}
	balance := 0.0
	for _, t := range p.transactions {
		if !t.date.After(date) {
			assets = append(assets, Asset{asset: t.asset, quantity: t.quantity})
			balance -= t.price * float64(t.quantity)
		}
	}
	return assets, balance
}

func (p *Portfolio) ApplyTransaction(t *Transaction) error {
	p.transactions = append(p.transactions, t)
	return nil
}

func (p *Portfolio) RenamePortfolio(name string) error {
	p.name = name
	return nil
}

func (p *Portfolio) ID() uuid.UUID {
	return p.id
}

func (p *Portfolio) UserID() uuid.UUID {
	return p.userID
}

func (p *Portfolio) Transactions() []*Transaction {
	return p.transactions
}

func (p *Portfolio) Name() string {
	return p.name
}

func (a *Asset) Asset() string {
	return a.asset
}

func (a *Asset) Quantity() int {
	return a.quantity
}
