package portfolio

import "github.com/google/uuid"

type Portfolio struct {
	id      uuid.UUID
	name    string
	assets  map[string]int
	balance float64
}
