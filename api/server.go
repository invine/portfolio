package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/invine/Portfolio/internal/models"
)

type PortfolioRepo interface {
	Read() (*models.Portfolio, error)
	Update(t models.Transaction) error
	Close() error
}

type PriceAPI interface {
	Price(symbol string) (float32, error)
}

type Server struct {
	r   *chi.Mux
	pr  PortfolioRepo
	api PriceAPI
}

func NewServer(pr PortfolioRepo, api PriceAPI) *Server {
	s := &Server{
		r:   chi.NewRouter(),
		pr:  pr,
		api: api,
	}
	return s
}

func (s *Server) ListenAndServe(p string) error {
	return http.ListenAndServe(p, s.r)
}

func (s *Server) ReadPortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	p, err := s.pr.Read()
	if err != nil {
		log.Printf("read portfolio: %v", err)
		rw.WriteHeader(500)
		return
	}

	type asset struct {
		Name   string  `json:"symbol"`
		Amount int     `json:"amount"`
		Price  float32 `json:"price"`
	}

	assets := []asset{}
	p.Assets()
	for k, v := range p.Assets() {
		price, err := s.api.Price(k)
		if err != nil {
			log.Printf("read portfolio: %v", err)
			price = 0
		}

		assets = append(assets, asset{
			Name:   k,
			Amount: v,
			Price:  price,
		})
	}

	bytes, err := json.Marshal(assets)
	if err != nil {
		log.Printf("read portfolio: %v", err)
		rw.WriteHeader(500)
		return
	}
	if _, err := rw.Write(bytes); err != nil {
		log.Printf("read portfolio: %v", err)
	}
}

func (s *Server) UpdatePortfolioHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	type transaction struct {
		Symbol string `json:"symbol"`
		Amount int    `json:"amount"`
	}

	t := new(transaction)
	if err := json.Unmarshal(bytes, t); err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	if err := s.pr.Update(models.Transaction{t.Symbol: t.Amount}); err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
}
