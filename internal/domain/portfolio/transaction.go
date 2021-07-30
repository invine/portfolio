package portfolio

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	id       uuid.UUID
	date     time.Time
	asset    string
	quantity int
	price    float64
}

func NewTransaction(id uuid.UUID, date time.Time, asset string, quantity int, price float64) (*Transaction, error) {
	t := &Transaction{}

	if err := t.setID(id); err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}
	if err := t.setDate(date); err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}
	if err := t.setAsset(asset); err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}
	if err := t.setQuantity(quantity); err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}
	if err := t.setPrice(price); err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}

	return t, nil
}

func (t *Transaction) UpdateTransaction(date time.Time, asset string, quantity int, price float64) error {
	if err := t.setDate(date); err != nil {
		return fmt.Errorf("can't update transaction: %w", err)
	}
	if err := t.setAsset(asset); err != nil {
		return fmt.Errorf("can't update transaction: %w", err)
	}
	if err := t.setQuantity(quantity); err != nil {
		return fmt.Errorf("can't update transaction: %w", err)
	}
	if err := t.setPrice(price); err != nil {
		return fmt.Errorf("can't update transaction: %w", err)
	}
	return nil
}

func (t *Transaction) setID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id can't be empty")
	}
	t.id = id
	return nil
}

func (t *Transaction) setDate(date time.Time) error {
	if date.IsZero() {
		return fmt.Errorf("date can't be empty")
	}
	t.date = date
	return nil
}

func (t *Transaction) setAsset(asset string) error {
	if asset == "" {
		return fmt.Errorf("asset can't be empty")
	}
	t.asset = asset
	return nil
}

func (t *Transaction) setQuantity(quantity int) error {
	if quantity == 0 {
		return fmt.Errorf("quantity can't be zero")
	}
	t.quantity = quantity
	return nil
}

func (t *Transaction) setPrice(price float64) error {
	if price < 0 {
		return fmt.Errorf("price can't be negative")
	}
	t.price = price
	return nil
}

func (t *Transaction) Asset() string {
	return t.asset
}

func (t *Transaction) Price() float64 {
	return t.price
}

func (t *Transaction) Quantity() int {
	return t.quantity
}

func (t *Transaction) Date() time.Time {
	return t.date
}

func (t *Transaction) ID() uuid.UUID {
	return t.id
}
