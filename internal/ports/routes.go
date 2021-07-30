package ports

import (
	"github.com/go-chi/cors"
)

func (s *Server) InitializeRoutes() {
	s.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.r.With(s.AuthenticateMiddleware).Get("/portfolio", s.ListPortfoliosHandler)
	s.r.With(s.AuthenticateMiddleware).Post("/portfolio", s.AddPortfolioHandler)
	s.r.With(s.AuthenticateMiddleware).Get("/portfolio/{id}", s.GetPortfolioHandler)
	s.r.With(s.AuthenticateMiddleware).Post("/portfolio/{id}", s.UpdatePortfolioHandler)
	s.r.With(s.AuthenticateMiddleware).Delete("/portfolio/{id}", s.DeletePortfolioHandler)
	s.r.With(s.AuthenticateMiddleware).Post("/portfolio/{id}/transaction", s.AddTransactionHandler)
	s.r.Post("/signin", s.UserSignInHandler)
	s.r.Post("/signup", s.UserSignUpHandler)
}
