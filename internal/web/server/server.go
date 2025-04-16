package server

import (
	"gateway/internal/service"
	"gateway/internal/web/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server representa o servidor HTTP
// Por que usar uma estrutura?
// Encapsula a configuração e execução do servidor, facilitando a manutenção.
type Server struct {
	router         *chi.Mux
	server         *http.Server
	accountService *service.AccountService
	port           string
}

// NewServer cria uma nova instância de Server
// Por que usar um construtor?
// Garante que o servidor seja inicializado corretamente com os serviços necessários.
func NewServer(accountService *service.AccountService, port string) *Server {
	return &Server{
		router:         chi.NewRouter(),
		accountService: accountService,
		port:           port,
	}
}

// ConfigureRoutes configura as rotas do servidor
// Por que configurar rotas separadamente?
// Facilita a manutenção e a adição de novas rotas.
func (s *Server) ConfigureRoutes() {
	accountHanler := handlers.NewAccountHandler(s.accountService)

	s.router.Post("/accounts", accountHanler.Create)
	s.router.Get("/accounts", accountHanler.Get)
}

// Start inicia o servidor HTTP
// Por que ter um método Start?
// Encapsula a lógica de inicialização do servidor, facilitando a reutilização.
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}
