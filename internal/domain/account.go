package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Account representa uma conta com suas informações e saldo protegido para acessos concorrentes
type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex
	CreatedAt time.Time
	UpdatedAt time.Time
}

// generateAPIKey gera uma chave API segura usando crypto/rand
func generateAPIKey() string {
	// Usa crypto/rand para garantir chaves API seguras
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// NewAccount cria uma conta com ID único, API Key segura e timestamps iniciais
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
func (a *Account) AddBalance(amount float64) {
	// Mutex garante exclusão mútua no acesso ao saldo
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	a.UpdatedAt = time.Now()
}