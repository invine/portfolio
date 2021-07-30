package app

import (
	"github.com/invine/Portfolio/internal/app/command"
	"github.com/invine/Portfolio/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ApplyTransaction command.ApplyTransactionHandler
	CreatePortfolio  command.CreatePortfolioHandler
	DeletePortfolio  command.DeletePortfolioHandler
}

type Queries struct {
	AllPortfolios query.AllPortfoliosHandler
	Portfolio     query.PortfolioHandler
}
