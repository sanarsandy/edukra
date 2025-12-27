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
	OrderID           *string         `json:"order_id,omitempty"`
	Amount            float64         `json:"amount"`
	GrossAmount       *float64        `json:"gross_amount,omitempty"`
	Currency          string          `json:"currency"`
	Status            string          `json:"status"` // pending, settlement, capture, deny, cancel, expire, failure, refund
	PaymentGateway    string          `json:"payment_gateway"`
	PaymentGatewayRef *string         `json:"payment_gateway_ref,omitempty"`
	PaymentMethod     *string         `json:"payment_method,omitempty"`
	PaymentType       *string         `json:"payment_type,omitempty"`
	SnapToken         *string         `json:"snap_token,omitempty"`
	PaymentURL        *string         `json:"payment_url,omitempty"`
	FraudStatus       *string         `json:"fraud_status,omitempty"`
	TransactionTime   *time.Time      `json:"transaction_time,omitempty"`
	SettlementTime    *time.Time      `json:"settlement_time,omitempty"`
	ExpiredAt         *time.Time      `json:"expired_at,omitempty"`
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
		SELECT id::text, tenant_id::text, user_id::text, course_id::text, order_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions WHERE id = $1
	`
	
	var tx Transaction
	var tenantID, courseID, orderID, gatewayRef, method sql.NullString
	var metadata []byte
	
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID, &tenantID, &tx.UserID, &courseID, &orderID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.PaymentGateway, &gatewayRef, &method, &metadata,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tenantID.Valid {
		tx.TenantID = tenantID.String
	}
	if courseID.Valid {
		tx.CourseID = &courseID.String
	}
	if orderID.Valid {
		tx.OrderID = &orderID.String
	}
	if gatewayRef.Valid {
		tx.PaymentGatewayRef = &gatewayRef.String
	}
	if method.Valid {
		tx.PaymentMethod = &method.String
	}
	if metadata != nil {
		tx.Metadata = metadata
	}
	
	return &tx, nil
}

// Create inserts a new transaction
func (r *TransactionRepository) Create(tx *Transaction) error {
	query := `
		INSERT INTO transactions (tenant_id, user_id, course_id, order_id, amount, currency, 
		                          status, payment_gateway, payment_gateway_ref, 
		                          payment_method, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
	
	// Handle empty TenantID - pass NULL to DB instead of empty string
	var tenantID interface{}
	if tx.TenantID == "" {
		tenantID = nil
	} else {
		tenantID = tx.TenantID
	}
	
	// Handle empty Metadata - pass NULL to DB instead of empty/invalid json
	var metadata interface{}
	if len(tx.Metadata) == 0 {
		metadata = nil
	} else {
		metadata = tx.Metadata
	}
	
	return r.db.QueryRow(query,
		tenantID, tx.UserID, tx.CourseID, tx.OrderID, tx.Amount, tx.Currency,
		tx.Status, tx.PaymentGateway, tx.PaymentGatewayRef,
		tx.PaymentMethod, metadata, tx.CreatedAt, tx.UpdatedAt,
	).Scan(&tx.ID)
}

// UpdateStatus updates the status of a transaction
func (r *TransactionRepository) UpdateStatus(id, status string) error {
	query := `UPDATE transactions SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, time.Now())
	return err
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(id string) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateGatewayRef updates the payment gateway reference
func (r *TransactionRepository) UpdateGatewayRef(id, gatewayRef string) error {
	query := `UPDATE transactions SET payment_gateway_ref = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, gatewayRef, time.Now())
	return err
}

// GetByOrderID retrieves a transaction by order_id (used for payment gateway callbacks)
func (r *TransactionRepository) GetByOrderID(orderID string) (*Transaction, error) {
	query := `
		SELECT id, tenant_id, user_id, course_id, order_id, amount, gross_amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, payment_type, snap_token,
		       payment_url, fraud_status, transaction_time, settlement_time, expired_at,
		       metadata, created_at, updated_at
		FROM transactions WHERE order_id = $1
	`
	
	var tx Transaction
	var tenantID, courseID, orderIDVal, gatewayRef, method, paymentType, snapToken, paymentURL, fraudStatus sql.NullString
	var grossAmount sql.NullFloat64
	var transactionTime, settlementTime, expiredAt sql.NullTime
	var metadata []byte
	
	err := r.db.QueryRow(query, orderID).Scan(
		&tx.ID, &tenantID, &tx.UserID, &courseID, &orderIDVal, &tx.Amount, &grossAmount, &tx.Currency,
		&tx.Status, &tx.PaymentGateway, &gatewayRef, &method, &paymentType, &snapToken,
		&paymentURL, &fraudStatus, &transactionTime, &settlementTime, &expiredAt,
		&metadata, &tx.CreatedAt, &tx.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		log.Printf("[GetByOrderID] No rows found for order_id: %s", orderID)
		return nil, nil
	}
	if err != nil {
		log.Printf("[GetByOrderID] Scan error for order_id %s: %v", orderID, err)
		return nil, err
	}
	
	if tenantID.Valid {
		tx.TenantID = tenantID.String
	}
	if courseID.Valid {
		tx.CourseID = &courseID.String
	}
	if orderIDVal.Valid {
		tx.OrderID = &orderIDVal.String
	}
	if grossAmount.Valid {
		tx.GrossAmount = &grossAmount.Float64
	}
	if gatewayRef.Valid {
		tx.PaymentGatewayRef = &gatewayRef.String
	}
	if method.Valid {
		tx.PaymentMethod = &method.String
	}
	if paymentType.Valid {
		tx.PaymentType = &paymentType.String
	}
	if snapToken.Valid {
		tx.SnapToken = &snapToken.String
	}
	if paymentURL.Valid {
		tx.PaymentURL = &paymentURL.String
	}
	if fraudStatus.Valid {
		tx.FraudStatus = &fraudStatus.String
	}
	if transactionTime.Valid {
		tx.TransactionTime = &transactionTime.Time
	}
	if settlementTime.Valid {
		tx.SettlementTime = &settlementTime.Time
	}
	if expiredAt.Valid {
		tx.ExpiredAt = &expiredAt.Time
	}
	
	return &tx, nil
}

// UpdateFromCallback updates transaction with payment gateway callback data
func (r *TransactionRepository) UpdateFromCallback(orderID string, status string, paymentType *string, fraudStatus *string, transactionTime *time.Time, settlementTime *time.Time, paymentGatewayRef *string) error {
	query := `
		UPDATE transactions SET 
			status = $2,
			payment_type = $3,
			fraud_status = $4,
			transaction_time = $5,
			settlement_time = $6,
			payment_gateway_ref = $7,
			updated_at = $8
		WHERE order_id = $1
	`
	_, err := r.db.Exec(query, orderID, status, paymentType, fraudStatus, transactionTime, settlementTime, paymentGatewayRef, time.Now())
	return err
}

// UpdateSnapToken updates the Midtrans snap token
func (r *TransactionRepository) UpdateSnapToken(id, snapToken, paymentURL string, expiredAt time.Time) error {
	query := `UPDATE transactions SET snap_token = $2, payment_url = $3, expired_at = $4, updated_at = $5 WHERE id = $1`
	_, err := r.db.Exec(query, id, snapToken, paymentURL, expiredAt, time.Now())
	return err
}

// GetPendingByUserAndCourse checks if user has pending transaction for a course
func (r *TransactionRepository) GetPendingByUserAndCourse(userID, courseID string) (*Transaction, error) {
	query := `
		SELECT id, snap_token, payment_url, expired_at FROM transactions 
		WHERE user_id = $1 AND course_id = $2 AND status = 'pending'
		AND (expired_at IS NULL OR expired_at > NOW())
		ORDER BY created_at DESC LIMIT 1
	`
	
	var tx Transaction
	var snapToken, paymentURL sql.NullString
	var expiredAt sql.NullTime
	
	err := r.db.QueryRow(query, userID, courseID).Scan(&tx.ID, &snapToken, &paymentURL, &expiredAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if snapToken.Valid {
		tx.SnapToken = &snapToken.String
	}
	if paymentURL.Valid {
		tx.PaymentURL = &paymentURL.String
	}
	if expiredAt.Valid {
		tx.ExpiredAt = &expiredAt.Time
	}
	
	return &tx, nil
}

// CancelPendingByUserAndCourse cancels existing pending transactions for a user and course
func (r *TransactionRepository) CancelPendingByUserAndCourse(userID, courseID string) error {
	query := `
		UPDATE transactions 
		SET status = 'cancelled', updated_at = NOW() 
		WHERE user_id = $1 AND course_id = $2 AND status = 'pending'
	`
	_, err := r.db.Exec(query, userID, courseID)
	return err
}

// ListByTenant retrieves transactions for a tenant
func (r *TransactionRepository) ListByTenant(tenantID string, limit, offset int) ([]*Transaction, error) {
	if tenantID == "" || tenantID == "default" {
		query := `
			SELECT id::text, tenant_id::text, user_id::text, course_id::text, order_id, amount, currency, status,
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
		SELECT id::text, tenant_id::text, user_id::text, course_id::text, order_id, amount, currency, status,
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
		SELECT id::text, tenant_id::text, user_id::text, course_id::text, order_id, amount, currency, status,
		       payment_gateway, payment_gateway_ref, payment_method, metadata,
		       created_at, updated_at
		FROM transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	return r.scanTransactions(query, userID, limit, offset)
}

// CountByUser returns total transaction count for a user
func (r *TransactionRepository) CountByUser(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// ListByStatus retrieves transactions by status
func (r *TransactionRepository) ListByStatus(tenantID, status string, limit, offset int) ([]*Transaction, error) {
	query := `
		SELECT id::text, tenant_id::text, user_id::text, course_id::text, order_id, amount, currency, status,
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
		var tenantID, courseID, orderID, gatewayRef, method sql.NullString
		var metadata []byte
		
		err := rows.Scan(
			&tx.ID, &tenantID, &tx.UserID, &courseID, &orderID, &tx.Amount, &tx.Currency,
			&tx.Status, &tx.PaymentGateway, &gatewayRef, &method, &metadata,
			&tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if tenantID.Valid {
			tx.TenantID = tenantID.String
		}
		if courseID.Valid {
			tx.CourseID = &courseID.String
		}
		if orderID.Valid {
			tx.OrderID = &orderID.String
		}
		if gatewayRef.Valid {
			tx.PaymentGatewayRef = &gatewayRef.String
		}
		if method.Valid {
			tx.PaymentMethod = &method.String
		}
		if metadata != nil {
			tx.Metadata = metadata
		}
		
		transactions = append(transactions, &tx)
	}
	
	return transactions, nil
}
