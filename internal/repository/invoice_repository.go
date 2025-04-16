package repository

import (
	"database/sql"
	"gateway/internal/domain"
)

// InvoiceRepository é uma estrutura que contém uma conexão com o banco de dados
// Por que usar uma estrutura?
// Usar uma estrutura permite encapsular a lógica de acesso ao banco de dados, facilitando a manutenção e a reutilização do código.
type InvoiceRepository struct {
	db *sql.DB
}

// NewInvoiceRepository cria uma nova instância de InvoiceRepository
// Por que usar um construtor?
// Um construtor como este garante que a estrutura seja inicializada corretamente com uma conexão de banco de dados válida.
func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

// Save salva uma nova fatura no banco de dados
// Por que ter um método Save?
// Ter um método Save centraliza a lógica de inserção de dados, garantindo que todas as faturas sejam salvas de maneira consistente.
func (r *InvoiceRepository) Save(invoice *domain.Invoice) error {
	// Executa uma instrução SQL para inserir uma nova fatura na tabela 'invoices'
	_, err := r.db.Exec(
		"INSERT INTO invoices(id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		invoice.ID, invoice.AccountID, invoice.Amount, invoice.Status, invoice.Description, invoice.PaymentType, invoice.CardLastDigits, invoice.CreatedAt, invoice.UpdatedAt,
	)
	if err != nil {
		return err // Retorna um erro se a execução falhar
	}

	return nil // Retorna nil se a execução for bem-sucedida
}

// FindByID busca uma fatura pelo ID
// Por que ter um método FindByID?
// Este método permite buscar uma fatura específica de forma eficiente, usando o ID como chave única.
func (r *InvoiceRepository) FindByID(id string) (*domain.Invoice, error) {
	var invoice domain.Invoice

	// Executa uma consulta SQL para buscar uma fatura pelo ID
	err := r.db.QueryRow(`SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at,
	updated_at FROM invoices WHERE id = $1
	`, id).Scan(
		&invoice.ID,
		&invoice.AccountID,
		&invoice.Amount,
		&invoice.Status,
		&invoice.Description,
		&invoice.PaymentType,
		&invoice.CardLastDigits,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, domain.ErrInvoiceNotFound // Retorna um erro se a fatura não for encontrada
	}

	return &invoice, nil // Retorna a fatura encontrada
}

// FindByAccountID busca todas as faturas de um determinado accountID
// Por que ter um método FindByAccountID?
// Este método é útil para recuperar todas as faturas associadas a uma conta específica, o que pode ser necessário para relatórios ou análises.
func (r *InvoiceRepository) FindByAccountID(accountID string) ([]*domain.Invoice, error) {
	// Executa uma consulta SQL para buscar todas as faturas de um accountID específico
	rows, err := r.db.Query(`SELECT id, account_id, amount, status, description, payment_type,
	card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1
	`, accountID)
	if err != nil {
		return nil, err // Retorna um erro se a consulta falhar
	}
	defer rows.Close() // Garante que as linhas sejam fechadas após o uso

	var invoices []*domain.Invoice

	// Itera sobre as linhas retornadas pela consulta
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.AccountID,
			&invoice.Amount,
			&invoice.Status,
		)
		if err != nil {
			return nil, err // Retorna um erro se a leitura das linhas falhar
		}
		invoices = append(invoices, &invoice) // Adiciona a fatura à lista de faturas
	}

	return invoices, nil // Retorna a lista de faturas
}

// UpdateStatus atualiza o status de uma fatura
// Por que ter um método UpdateStatus?
// Atualizar o status de uma fatura é uma operação comum, e ter um método dedicado a isso garante que a lógica de atualização seja consistente e fácil de manter.
func (r *InvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	// Executa uma instrução SQL para atualizar o status de uma fatura
	rows, err := r.db.Exec(`UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3
	`, invoice.Status, invoice.UpdatedAt, invoice.ID)
	if err != nil {
		return err // Retorna um erro se a execução falhar
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err // Retorna um erro se a verificação das linhas afetadas falhar
	}

	if rowsAffected == 0 {
		return domain.ErrInvoiceNotFound // Retorna um erro se nenhuma linha foi afetada (fatura não encontrada)
	}

	return nil // Retorna nil se a execução for bem-sucedida
}
