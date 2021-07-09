package api

import "github.com/go-chi/cors"

func (s *Server) InitializeRoutes() {
	s.r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// s.r.Get("/s/{symbol}", s.ReadAssetHandler)

	// s.r.Route("/u", func(r chi.Router) {
	s.r.Get("/p", s.ReadPortfolioHandler)
	s.r.Post("/t", s.UpdatePortfolioHandler)
	s.r.Get("/q/{symbol}", s.ReadPriceHandler)
	s.r.Get("/q/{symbol}/{date}", s.ReadPriceHistoricHandler)
	// })
}
