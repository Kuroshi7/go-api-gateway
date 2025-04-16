package server

import (
	"gateway/internal/service"
	"gateway/internal/web/handlers"
	"gateway/internal/web/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Server representa o servidor HTTP
// Por que usar uma estrutura?
// Encapsula a configuração e execução do servidor, facilitando a manutenção.
type Server struct {
	router         *chi.Mux
	server         *http.Server
	accountService *service.AccountService
	invoiceService *service.InvoiceService
	port           string
}

// NewServer cria uma nova instância de Server
// Por que usar um construtor?
// Garante que o servidor seja inicializado corretamente com os serviços necessários.
func NewServer(accountService *service.AccountService, invoiceService *service.InvoiceService, port string) *Server {
	return &Server{
		router:         chi.NewRouter(),
		accountService: accountService,
		invoiceService: invoiceService,
		port:           port,
	}
}

// ConfigureRoutes configura as rotas do servidor
// Por que configurar rotas separadamente?
// Facilita a manutenção e a adição de novas rotas.
func (s *Server) ConfigureRoutes() {
	accountHanler := handlers.NewAccountHandler(s.accountService)
	invoiceHandler := handlers.NewInvoiceHandler(s.invoiceService)
	authMiddleware := middleware.NewAuthMiddleware(s.accountService)

	s.router.Post("/accounts", accountHanler.Create)
	s.router.Get("/accounts", accountHanler.Get)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate) // Middleware de autenticação

		s.router.Post("/invoice", invoiceHandler.Create)
		s.router.Get("/invoice", invoiceHandler.ListByAccount)
		s.router.Get("/invoice/{id}", invoiceHandler.GetByID)
	})
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
