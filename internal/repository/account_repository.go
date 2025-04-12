package repository

import (
	"database/sql"
	"time"
	"gateway/internal/domain"
	
)

// AccountRepository implementa operações de persistência para Account
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository cria um novo repositório de contas
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Save persiste uma nova conta no banco de dados
// Retorna erro se houver falha na inserção
func (r *AccountRepository) Save(account *domain.Account) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// FindByAPIKey busca uma conta pelo API Key
// Retorna ErrAccountNotFound se não encontrada
func (r *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE api_key = $1
	`, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt
	return &account, nil
}

// FindByID busca uma conta pelo ID
// Retorna ErrAccountNotFound se não encontrada
func (r *AccountRepository) FindByID(id string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt
	return &account, nil
}

// UpdateBalance atualiza o saldo da conta usando SELECT FOR UPDATE para consistência em acessos concorrentes
// Retorna ErrAccountNotFound se a conta não existir
func (r *AccountRepository) UpdateBalance(account *domain.Account) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// SELECT FOR UPDATE previne race conditions no saldo
	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`,
		account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        UPDATE accounts
        SET balance = $1, updated_at = $2
        WHERE id = $3
    `, account.Balance, time.Now(), account.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}