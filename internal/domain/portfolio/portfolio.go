package portfolio

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Assets map[string]int

type Snapshot struct {
	ID      uuid.UUID
	Name    string
	Assets  Assets
	Balance float64
}

type Portfolio struct {
	id           uuid.UUID
	userID       uuid.UUID
	name         string
	transactions []*Transaction
}

func NewPortfolio(id, userID uuid.UUID, name string, transactions []*Transaction) (*Portfolio, error) {
	p := &Portfolio{}
	if err := p.setID(id); err != nil {
		return nil, fmt.Errorf("can't create portfolio: %w", err)
	}
	if err := p.setUserID(userID); err != nil {
		return nil, fmt.Errorf("can't create portfolio: %w", err)
	}
	if err := p.setName(name); err != nil {
		return nil, fmt.Errorf("can't create portfolio: %w", err)
	}
	p.transactions = transactions
	if transactions == nil {
		p.transactions = []*Transaction{}
	}
	return p, nil
}

func (p *Portfolio) Snapshot(date time.Time) *Snapshot {
	assets := map[string]int{}
	balance := 0.0
	for _, t := range p.transactions {
		if !t.date.After(date) {
			assets[t.Asset()] += t.Quantity()
			balance -= t.price * float64(t.quantity)
		}
	}
	return &Snapshot{
		ID:      p.ID(),
		Name:    p.Name(),
		Assets:  assets,
		Balance: balance,
	}
}

func (p *Portfolio) ApplyTransaction(t *Transaction) error {
	if p.Snapshot(t.Date()).Assets[t.Asset()]+t.Quantity() < 0 {
		return fmt.Errorf("can't apply transaction: asset quantity can't be less than zero")
	}
	p.transactions = append(p.transactions, t)
	return nil
}

func (p *Portfolio) RenamePortfolio(name string) error {
	p.setName(name)
	return nil
}

func (p *Portfolio) setID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is mandatory")
	}
	p.id = id
	return nil
}

func (p *Portfolio) setUserID(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf("id is mandatory")
	}
	p.userID = userID
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
