package handlers

import (
	"encoding/json"
	"net/http"

	"gateway/internal/dto"
	"gateway/internal/service"
)

// AccountHandler processa requisições HTTP relacionadas a contas
// Por que usar um handler?
// Separa a lógica de processamento de requisições HTTP, facilitando a manutenção e testes.
type AccountHandler struct {
	accountService *service.AccountService
}

// NewAccountHandler cria um novo handler de contas
// Por que usar um construtor?
// Garante que o handler seja inicializado corretamente com o serviço necessário.
func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// Create processa POST /accounts
// Por que ter um método Create?
// Encapsula a lógica de criação de contas via HTTP, garantindo consistência.
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAccountInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.accountService.CreateAccount(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// Get processa GET /accounts
// Por que ter um método Get?
// Encapsula a lógica de busca de contas via HTTP, garantindo consistência.
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API Key is required", http.StatusUnauthorized)
		return
	}

	output, err := h.accountService.FindByAPIKey(apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
