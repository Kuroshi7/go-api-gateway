package dto

import (
	"gateway/internal/domain"
	"time"
)

// CreateAccountInput representa dados para criação de conta
type CreateAccountInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AccountOutput representa dados da conta nas respostas da API
type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	APIKey    string    `json:"api_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToAccount converte CreateAccountInput para domain.Account
func ToAccount(input CreateAccountInput) *domain.Account {
	return domain.NewAccount(input.Name, input.Email)
}

// FromAccount converte domain.Account para AccountOutput
func FromAccount(account *domain.Account) AccountOutput {
	return AccountOutput{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Balance:   account.Balance,
		APIKey:    account.APIKey,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}