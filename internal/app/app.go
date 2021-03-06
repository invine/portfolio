package app

import (
	"github.com/invine/portfolio/internal/app/command"
	"github.com/invine/portfolio/internal/app/query"
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
	AllPortfolios   query.AllPortfoliosHandler
	AllTransactions query.AllTransactionsHandler
	Portfolio       query.PortfolioHandler
}
