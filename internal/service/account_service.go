package service

import (
	"gateway/internal/domain"
	"gateway/internal/dto"
)

// AccountService implementa a lógica de negócios para operações com Account
// Por que usar um serviço?
// Separa a lógica de negócios da lógica de persistência, facilitando a manutenção e testes.
type AccountService struct {
	repository domain.AccountRepository
}

// NewAccountService cria um novo serviço de contas
// Por que usar um construtor?
// Garante que o serviço seja inicializado corretamente com um repositório válido.
func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

// CreateAccount cria uma nova conta e valida duplicidade de API Key
// Por que validar duplicidade?
// Garante que cada conta tenha uma chave API única, evitando conflitos.
func (s *AccountService) CreateAccount(input dto.CreateAccountInput) (*dto.AccountOutput, error) {
	account := dto.ToAccount(input)

	// Verifica duplicidade de API Key antes da criação
	existingAccount, err := s.repository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}
	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey
	}

	err = s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

// UpdateBalance atualiza o saldo de uma conta de forma thread-safe
// Por que thread-safe?
// Garante que o saldo seja atualizado corretamente mesmo com acessos concorrentes.
func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return &output, nil
}

// FindByAPIKey busca uma conta pelo API Key
// Por que ter um método FindByAPIKey?
// Permite buscar uma conta de forma eficiente usando a chave API como identificador único.
func (s *AccountService) FindByAPIKey(apiKey string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return &output, nil
}

// FindByID busca uma conta pelo ID
// Por que ter um método FindByID?
// Permite buscar uma conta específica usando o ID como chave única.
func (s *AccountService) FindByID(id string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return &output, nil
}
