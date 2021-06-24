package main

import (
	"log"

	"github.com/invine/Portfolio/api"
	"github.com/invine/Portfolio/internal/repos"
	"github.com/invine/Portfolio/internal/yahooapi"
)

const API_KEY = "2UVV3YNX8IY2LNF6"

func main() {
	pr, err := repos.NewPortfolioRepo("./db.sqlite3")
	if err != nil {
		panic(err)
	}

	// av := alphavantage.NewAlphaVantage(API_KEY)
	y := yahooapi.NewYahooAPI()

	s := api.NewServer(pr, y)
	s.InitializeRoutes()

	log.Fatal(s.ListenAndServe(":3000"))
}
