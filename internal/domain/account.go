package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Account representa uma conta com suas informações e saldo protegido para acessos concorrentes
// Por que usar uma estrutura?
// Encapsula dados relacionados, facilitando o gerenciamento e a manipulação.
type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex // Mutex para garantir exclusão mútua no acesso ao saldo
	CreatedAt time.Time
	UpdatedAt time.Time
}

// generateAPIKey gera uma chave API segura usando crypto/rand
// Por que usar crypto/rand?
// Garante que as chaves API sejam seguras e imprevisíveis.
func generateAPIKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// NewAccount cria uma conta com ID único, API Key segura e timestamps iniciais
// Por que usar um construtor?
// Garante que a conta seja inicializada corretamente com dados válidos.
func NewAccount(name, email string) *Account {
	account := &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Balance:   0,
		APIKey:    generateAPIKey(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

// AddBalance modifica o saldo da conta de forma thread-safe
// Por que thread-safe?
// Garante que o saldo seja atualizado corretamente mesmo com acessos concorrentes.
func (a *Account) AddBalance(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	a.UpdatedAt = time.Now()
}
