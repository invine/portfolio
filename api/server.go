package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/invine/Portfolio/internal/models"
)

const CONCURENT_THREADS = 4

type UserRepo interface {
	CreateUser(email, login, password, name string) error
	Authenticate(login, password string) (uuid.UUID, error)
	Close() error
}

type PortfolioRepo interface {
	Read() (*models.Portfolio, error)
	Update(t models.Transaction) error
	Close() error
}

type PriceAPI interface {
	Price(symbol string) (float32, error)
	PriceHistorical(symbol string, t time.Time) (float32, error)
}

type Server struct {
	r   *chi.Mux
	pr  PortfolioRepo
	api PriceAPI
	ur  UserRepo
	key []byte
}

func NewServer(ur UserRepo, pr PortfolioRepo, api PriceAPI, key []byte) *Server {
	s := &Server{
		r:   chi.NewRouter(),
		pr:  pr,
		api: api,
		ur:  ur,
		key: key,
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

	out := make(chan asset)
	go func() {
		for n, a := range p.Assets() {
			out <- asset{
				Name:   n,
				Amount: a,
				Price:  0,
			}
		}
		close(out)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(CONCURENT_THREADS)
	resC := make(chan asset)

	go func() {
		for i := 0; i < CONCURENT_THREADS; i++ {
			go func(in <-chan asset) {
				defer wg.Done()
				for a := range in {
					a.Price, err = s.api.Price(a.Name)
					if err != nil {
						log.Printf("read portfolio: %v", err)
					}
					resC <- a
				}
			}(out)
		}
		wg.Wait()
		close(resC)
	}()

	// for k, v := range p.Assets() {
	// 	p, err := s.api.Price(k)
	// 	if err != nil {
	// 		log.Printf("read portfolio: %v", err)
	// 	}

	// 	assets = append(assets, asset{
	// 		Name:   k,
	// 		Amount: v,
	// 		Price:  p,
	// 	})
	// }
	for a := range resC {
		assets = append(assets, a)
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

func (s *Server) ReadPriceHandler(rw http.ResponseWriter, r *http.Request) {
	price, err := s.api.Price(chi.URLParam(r, "symbol"))
	if err != nil {
		log.Printf("get price: %v", err)
		rw.WriteHeader(500)
		return
	}

	if _, err := rw.Write([]byte(fmt.Sprintf("%f", price))); err != nil {
		log.Printf("get price: %v", err)
	}
}

func (s *Server) ReadPriceHistoricHandler(rw http.ResponseWriter, r *http.Request) {
	timeString := chi.URLParam(r, "date")
	t, err := time.Parse("20060102", timeString)
	if err != nil {
		log.Printf("historical price can't parse time %s", timeString)
		rw.WriteHeader(400)
		return
	}
	price, err := s.api.PriceHistorical(chi.URLParam(r, "symbol"), t)
	if err != nil {
		log.Printf("historical price: %v", err)
		rw.WriteHeader(500)
		return
	}

	if _, err := rw.Write([]byte(fmt.Sprintf("%f", price))); err != nil {
		log.Printf("historical price: %v", err)
	}
}

func (s *Server) UserSignInHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	type credintials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	c := new(credintials)
	if err := json.Unmarshal(bytes, c); err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(400)
		return
	}

	id, err := s.ur.Authenticate(c.Login, c.Password)
	if err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(401)
		return
	}

	// TODO add jwt generation
	log.Println(id.String())

	// Create the Claims
	claims := UserClaims{
		id.String(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 15000,
			Issuer:    "portfolio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(s.key)
	if err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(500)
		return
	}

	if _, err := rw.Write([]byte(ss)); err != nil {
		log.Printf("failed sign in: %v", err)
	}
}

func (s *Server) UserSignUpHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("update portfolio: %v", err)
		rw.WriteHeader(400)
		return
	}

	type user struct {
		Email    string `json:"email"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	u := new(user)
	if err := json.Unmarshal(bytes, u); err != nil {
		log.Printf("failed sign up: %v", err)
		rw.WriteHeader(400)
		return
	}

	if err := s.ur.CreateUser(u.Email, u.Login, u.Password, u.Name); err != nil {
		log.Printf("failed sign up: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
}
