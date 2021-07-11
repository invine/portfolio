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
	key := []byte(getenv("JWT_KEY", "34$FtGVP*8Uzhp"))

	db_conn := fmt.Sprintf("%s/db.sqlite3", db_path)

	ur, err := repos.NewUserRepo(db_conn)
	if err != nil {
		panic(err)
	}

	pr, err := repos.NewPortfolioRepo(db_conn)
	if err != nil {
		panic(err)
	}

	// av := alphavantage.NewAlphaVantage(API_KEY)
	y := yahooapi.NewYahooAPI()

	s := api.NewServer(ur, pr, y, key)
	s.InitializeRoutes()

	log.Printf("Starting server on %s...", port)

	log.Fatal(s.ListenAndServe(fmt.Sprintf(":%s", port)))
}
