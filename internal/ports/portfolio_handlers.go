package ports

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/app/command"
	"github.com/invine/Portfolio/internal/app/query"
	"github.com/invine/Portfolio/internal/domain/portfolio"
)

type assetModel struct {
	Asset    string `json:"asset"`
	Quantity int    `json:"quantity"`
}

type portfolioModel struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Assets  []assetModel `json:"assets"`
	Balance float64      `json:"balance"`
}

func (s *Server) ListPortfoliosHandler(rw http.ResponseWriter, r *http.Request) {
	u, err := UserFromCtx(r.Context())
	if err != nil {
		log.Printf("list portfolios: %v", err)
		rw.WriteHeader(400)
		return
	}

	// TODO add pagination
	ps, err := s.app.Queries.AllPortfolios.Handle(
		r.Context(),
		query.AllPortfolios{
			UserID: u.ID,
		},
	)
	if err != nil {
		log.Printf("list portfolios: %v", err)
		rw.WriteHeader(400)
		return
	}

	var pms []portfolioModel
	for _, p := range ps {
		pms = append(pms, portfolioToPortfolioModel(p, time.Time{}))
	}
	bytes, err := json.Marshal(pms)
	if err != nil {
		log.Printf("list portfolios: %v", err)
		rw.WriteHeader(500)
		return
	}
	if _, err := rw.Write(bytes); err != nil {
		log.Printf("list portfolios: %v", err)
	}
}

func (s *Server) AddPortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	u, err := UserFromCtx(r.Context())
	if err != nil {
		log.Printf("create portfolio: %v", err)
		rw.WriteHeader(401)
		return
	}

	var pm portfolioModel
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("create portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}
	err = json.Unmarshal(bytes, &pm)
	if err != nil {
		log.Printf("create portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	err = s.app.Commands.CreatePortfolio.Handle(
		r.Context(),
		command.CreatePortfolio{
			ID:     uuid.New(),
			UserID: u.ID,
			Name:   pm.Name,
		},
	)
	if err != nil {
		log.Printf("create portfolio: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(201)
}

func (s *Server) GetPortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	u, err := UserFromCtx(r.Context())
	if err != nil {
		log.Printf("get portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("get portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	timeString := r.URL.Query().Get("date")
	var t time.Time
	if timeString == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("20060102", timeString)
		if err != nil {
			log.Printf("get portfolio: can't parse date %s: %v", timeString, err)
			rw.WriteHeader(400)
			return
		}
	}

	p, err := s.app.Queries.Portfolio.Handle(
		r.Context(),
		query.Portfolio{
			UserID: u.ID,
			ID:     id,
		},
	)
	if err != nil {
		log.Printf("get portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	pm := portfolioToPortfolioModel(p, t)
	bytes, err := json.Marshal(pm)
	if err != nil {
		log.Printf("get portfolio: %v", err)
		rw.WriteHeader(500)
		return
	}
	if _, err := rw.Write(bytes); err != nil {
		log.Printf("get portfolio: %v", err)
	}
}

func (s *Server) UpdatePortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	var pm portfolioModel
	err = json.Unmarshal(bytes, &pm)
	if err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	log.Printf("update portfolio: %v", pm)
	//TODO add command

	rw.WriteHeader(501)
}

func (s *Server) DeletePortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	// TODO implement
	rw.WriteHeader(501)
}

func (s *Server) AddTransactionHandler(rw http.ResponseWriter, r *http.Request) {
	type transactionModel struct {
		Symbol string    `json:"symbol"`
		Amount int       `json:"amount"`
		Date   time.Time `json:"date"`
		Price  float64   `json:"price"`
	}

	u, err := UserFromCtx(r.Context())
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(400)
		return
	}

	var trm transactionModel
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(400)
		return
	}
	err = json.Unmarshal(bytes, &trm)
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(400)
		return
	}

	portfolioID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(400)
		return
	}

	tr, err := portfolio.NewTransaction(
		uuid.New(),
		trm.Date,
		trm.Symbol,
		trm.Amount,
		trm.Price,
	)
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(400)
		return
	}

	err = s.app.Commands.ApplyTransaction.Handle(
		r.Context(),
		command.ApplyTransaction{
			UserID:      u.ID,
			PortfolioID: portfolioID,
			Transaction: tr,
		},
	)
	if err != nil {
		log.Printf("add transaction: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(201)
}

func portfolioToPortfolioModel(p *portfolio.Portfolio, date time.Time) portfolioModel {
	pm := portfolioModel{
		ID:   p.ID().String(),
		Name: p.Name(),
	}
	if !date.IsZero() {
		assets, balance := p.Snapshot(date)
		pm.Assets = assetsToAssetsModel(assets)
		pm.Balance = balance
	}
	return pm
}

func assetsToAssetsModel(assets portfolio.Assets) []assetModel {
	res := []assetModel{}
	for k, v := range assets {
		res = append(res, assetModel{
			Asset:    k,
			Quantity: v,
		})
	}
	return res
}
