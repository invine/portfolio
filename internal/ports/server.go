package ports

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/invine/portfolio/internal/app"
)

type Server struct {
	r       *chi.Mux
	app     app.Application
	userSvc *app.UserService
	key     []byte
}

func NewServer(userSvc *app.UserService, app app.Application, key []byte) *Server {
	s := &Server{
		r:       chi.NewRouter(),
		app:     app,
		userSvc: userSvc,
		key:     key,
	}
	return s
}

func (s *Server) ListenAndServe(p string) error {
	return http.ListenAndServe(p, s.r)
}
