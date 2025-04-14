package server

import (
	"gateway/internal/service"
	"gateway/internal/web/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	server *http.Server
	accountService *service.AccountService
	port 	string
}

func NewServer (accountService *service.AccountService, port string) *Server{
	return &Server{
		router: chi.NewRouter(),
		accountService: accountService,
		port: port,
	}
}

func (s *Server) ConfigureRoutes(){
	accountHanler := handlers.NewAccountHandler(s.accountService)
	
	s.router.Post("/accounts", accountHanler.Create)
	s.router.Get("/accounts", accountHanler.Get)
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}