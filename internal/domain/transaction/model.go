package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	id          uuid.UUID
	userID      uuid.UUID
	portfolioID uuid.UUID
	date        time.Time
	asset       string
	quantity    int
	price       float64
}
