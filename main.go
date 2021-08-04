package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/invine/Portfolio/internal/adapters"
	"github.com/invine/Portfolio/internal/app"
	"github.com/invine/Portfolio/internal/app/command"
	"github.com/invine/Portfolio/internal/app/query"
	"github.com/invine/Portfolio/internal/ports"
	_ "github.com/mattn/go-sqlite3"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	db_path := getenv("DB_PATH", ".")
	port := getenv("PORT", "3001")
	key := []byte(getenv("JWT_KEY", "34$FtGVP*8Uzhp"))

	db_conn := fmt.Sprintf("%s/db.sqlite3", db_path)
	db, err := sql.Open("sqlite3", db_conn)
	if err != nil {
		panic(err)
	}
	userRepo, err := adapters.NewSQLiteUsersRepository(db)
	if err != nil {
		panic(err)
	}
	userService, err := app.NewUserService(userRepo)
	if err != nil {
		panic(err)
	}

	portfolioRepo, err := adapters.NewSQLitePortfolioRepository(db)
	if err != nil {
		panic(err)
	}

	applyTransactionHandler, err := command.NewApplyTransactionHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}
	createPortfolioHandler, err := command.NewCreatePortfolioHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}
	deletePortfolioHandler, err := command.NewDeletePortfolioHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}
	allPortfoliosHandler, err := query.NewAllPortfoliosHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}
	AllTransactionsHandler, err := query.NewAllTransactionsHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}
	portfolioSnapshotHandler, err := query.NewPortfolioHandler(portfolioRepo)
	if err != nil {
		panic(err)
	}

	app := app.Application{
		Commands: app.Commands{
			ApplyTransaction: *applyTransactionHandler,
			CreatePortfolio:  *createPortfolioHandler,
			DeletePortfolio:  *deletePortfolioHandler,
		},
		Queries: app.Queries{
			AllPortfolios:   *allPortfoliosHandler,
			AllTransactions: *AllTransactionsHandler,
			Portfolio:       *portfolioSnapshotHandler,
		},
	}

	s := ports.NewServer(userService, app, key)
	s.InitializeRoutes()

	log.Printf("Starting server on %s...", port)

	log.Fatal(s.ListenAndServe(fmt.Sprintf(":%s", port)))
}
