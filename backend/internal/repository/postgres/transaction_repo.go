package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// Transaction represents a payment transaction
type Transaction struct {
	ID                string          `json:"id"`
	TenantID          string          `json:"tenant_id"`
	UserID            string          `json:"user_id"`
	CourseID          *string         `json:"course_id,omitempty"`
	Amount            float64         `json:"amount"`
	Currency          string          `json:"currency"`
	Status            string          `json:"status"` // pending, success, failed, refunded
	PaymentGateway    string          `json:"payment_gateway"`
	PaymentGatewayRef *string         `json:"payment_gateway_ref,omitempty"`
	PaymentMethod     *string         `json:"payment_method,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// TransactionRepository handles transaction data access
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id string) (*Transaction, error) {
	query := `
		SELECT id, tenant_id, user_id, course_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions WHERE id = $1
	`
	
	var tx Transaction
	var courseID, gatewayRef, method sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID, &tx.TenantID, &tx.UserID, &courseID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.PaymentGateway, &gatewayRef, &method, &tx.Metadata,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if courseID.Valid {
		tx.CourseID = &courseID.String
	}
	if gatewayRef.Valid {
		tx.PaymentGatewayRef = &gatewayRef.String
	}
	if method.Valid {
		tx.PaymentMethod = &method.String
	}
	
	return &tx, nil
}

// Create inserts a new transaction
func (r *TransactionRepository) Create(tx *Transaction) error {
	query := `
		INSERT INTO transactions (tenant_id, user_id, course_id, amount, currency, 
		                          status, payment_gateway, payment_gateway_ref, 
		                          payment_method, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	
	now := time.Now()
	tx.CreatedAt = now
	tx.UpdatedAt = now
	
	if tx.Status == "" {
		tx.Status = "pending"
	}
	if tx.Currency == "" {
		tx.Currency = "IDR"
	}
	
	return r.db.QueryRow(query,
		tx.TenantID, tx.UserID, tx.CourseID, tx.Amount, tx.Currency,
		tx.Status, tx.PaymentGateway, tx.PaymentGatewayRef,
		tx.PaymentMethod, tx.Metadata, tx.CreatedAt, tx.UpdatedAt,
	).Scan(&tx.ID)
}

// UpdateStatus updates the status of a transaction
func (r *TransactionRepository) UpdateStatus(id, status string) error {
	query := `UPDATE transactions SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, time.Now())
	return err
}

// UpdateGatewayRef updates the payment gateway reference
func (r *TransactionRepository) UpdateGatewayRef(id, gatewayRef string) error {
	query := `UPDATE transactions SET payment_gateway_ref = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, gatewayRef, time.Now())
	return err
}

// ListByTenant retrieves transactions for a tenant
func (r *TransactionRepository) ListByTenant(tenantID string, limit, offset int) ([]*Transaction, error) {
	if tenantID == "" || tenantID == "default" {
		query := `
			SELECT id, tenant_id, user_id, course_id, amount, currency, status,
			       payment_gateway, payment_gateway_ref, payment_method, metadata,
			       created_at, updated_at
			FROM transactions 
			WHERE tenant_id IS NULL
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		return r.scanTransactions(query, limit, offset)
	}
	
	query := `
		SELECT id, tenant_id, user_id, course_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions 
		WHERE tenant_id = $1 OR tenant_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	return r.scanTransactions(query, tenantID, limit, offset)
}

// ListByUser retrieves transactions for a user
func (r *TransactionRepository) ListByUser(userID string, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, tenant_id, user_id, course_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	return r.scanTransactions(query, userID, limit, offset)
}

// ListByStatus retrieves transactions by status
func (r *TransactionRepository) ListByStatus(tenantID, status string, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id, tenant_id, user_id, course_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions 
		WHERE (tenant_id = $1 OR tenant_id IS NULL) AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	
	rows, err := r.db.Query(query, tenantID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTransactionRows(rows)
}

// CountByTenant returns total transaction count for a tenant
func (r *TransactionRepository) CountByTenant(tenantID string) (int, error) {
	var query string
	var count int
	var err error
	
	if tenantID == "" || tenantID == "default" {
		query = `SELECT COUNT(*) FROM transactions WHERE tenant_id IS NULL`
		err = r.db.QueryRow(query).Scan(&count)
	} else {
		query = `SELECT COUNT(*) FROM transactions WHERE tenant_id = $1 OR tenant_id IS NULL`
		err = r.db.QueryRow(query, tenantID).Scan(&count)
	}
	return count, err
}

// SumByTenant returns total revenue for a tenant
func (r *TransactionRepository) SumByTenant(tenantID string) (float64, error) {
	var query string
	var sum float64
	var err error
	
	if tenantID == "" || tenantID == "default" {
		query = `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE tenant_id IS NULL AND status = 'success'`
		err = r.db.QueryRow(query).Scan(&sum)
	} else {
		query = `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE (tenant_id = $1 OR tenant_id IS NULL) AND status = 'success'`
		err = r.db.QueryRow(query, tenantID).Scan(&sum)
	}
	return sum, err
}

// Helper function to scan transactions
func (r *TransactionRepository) scanTransactions(query string, args ...interface{}) ([]*Transaction, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTransactionRows(rows)
}

func (r *TransactionRepository) scanTransactionRows(rows *sql.Rows) ([]*Transaction, error) {
	var transactions []*Transaction
	for rows.Next() {
		var tx Transaction
		var courseID, gatewayRef, method sql.NullString
		
		err := rows.Scan(
			&tx.ID, &tx.TenantID, &tx.UserID, &courseID, &tx.Amount, &tx.Currency,
			&tx.Status, &tx.PaymentGateway, &gatewayRef, &method, &tx.Metadata,
			&tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if courseID.Valid {
			tx.CourseID = &courseID.String
		}
		if gatewayRef.Valid {
			tx.PaymentGatewayRef = &gatewayRef.String
		}
		if method.Valid {
			tx.PaymentMethod = &method.String
		}
		
		transactions = append(transactions, &tx)
	}
	
	return transactions, nil
}
