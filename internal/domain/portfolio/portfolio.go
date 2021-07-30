package portfolio

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Assets map[string]int

type Portfolio struct {
	id           uuid.UUID
	userID       uuid.UUID
	name         string
	transactions []*Transaction
}

func NewPortfolio(id, userID uuid.UUID, name string) (*Portfolio, error) {
	p := &Portfolio{
		id:     id,
		userID: userID,
	}
	if err := p.setName(name); err != nil {
		return nil, fmt.Errorf("can't create portfolio: %w", err)
	}

	return p, nil
}

func NewPortfolioWithTransactions(id, userID uuid.UUID, name string, transactions []*Transaction) (*Portfolio, error) {
	p := &Portfolio{
		id:           id,
		userID:       userID,
		transactions: transactions,
	}
	p.setName(name)
	if err := p.setName(name); err != nil {
		return nil, fmt.Errorf("can't create portfolio: %w", err)
	}
	return p, nil
}

func (p *Portfolio) Snapshot(date time.Time) (Assets, float64) {
	assets := map[string]int{}
	balance := 0.0
	for _, t := range p.transactions {
		if !t.date.After(date) {
			assets[t.Asset()] += t.Quantity()
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

func (p *Portfolio) setName(name string) error {
	if name == "" {
		return fmt.Errorf("name is mandatory")
	}
	p.name = name
	return nil
}
