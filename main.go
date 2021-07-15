package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/invine/Portfolio/archive/api"
	"github.com/invine/Portfolio/archive/repos"
	"github.com/invine/Portfolio/archive/yahooapi"
	"github.com/invine/Portfolio/internal/adapters"
	"github.com/invine/Portfolio/internal/app"
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

	pr, err := repos.NewPortfolioRepo(db_conn)
	if err != nil {
		panic(err)
	}

	// av := alphavantage.NewAlphaVantage(API_KEY)
	y := yahooapi.NewYahooAPI()

	s := api.NewServer(userService, pr, y, key)
	s.InitializeRoutes()

	log.Printf("Starting server on %s...", port)

	log.Fatal(s.ListenAndServe(fmt.Sprintf(":%s", port)))
}
