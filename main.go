package main

import (
	"fmt"
	"log"
	"os"

	"github.com/invine/Portfolio/api"
	"github.com/invine/Portfolio/internal/repos"
	"github.com/invine/Portfolio/internal/yahooapi"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// const API_KEY = "2UVV3YNX8IY2LNF6"

func main() {
	db_path := getenv("DB_PATH", ".")
	port := getenv("PORT", "3001")

	pr, err := repos.NewPortfolioRepo(fmt.Sprintf("%s/db.sqlite3", db_path))
	if err != nil {
		panic(err)
	}

	// av := alphavantage.NewAlphaVantage(API_KEY)
	y := yahooapi.NewYahooAPI()

	s := api.NewServer(pr, y)
	s.InitializeRoutes()

	log.Printf("Starting server on %s...", port)

	log.Fatal(s.ListenAndServe(fmt.Sprintf(":%s", port)))
}
