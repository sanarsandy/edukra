package payment

import (
	"context"
	"time"
)

// PaymentProvider defines the interface for payment gateway providers
type PaymentProvider interface {
	// GetName returns the provider name
	GetName() string
	
	// CreateTransaction creates a new payment transaction and returns payment details
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	
	// HandleNotification processes payment gateway webhook/callback
	HandleNotification(ctx context.Context, data []byte) (*NotificationResult, error)
	
	// VerifySignature verifies the signature of a notification
	VerifySignature(data map[string]interface{}) bool
	
	// GetTransactionStatus gets the current status of a transaction
	GetTransactionStatus(ctx context.Context, orderID string) (*TransactionStatus, error)
}

// CreateTransactionRequest holds the data needed to create a payment
type CreateTransactionRequest struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	CustomerPhone string  `json:"customer_phone,omitempty"`
	ItemName      string  `json:"item_name"`
	ItemID        string  `json:"item_id"`
	ItemCategory  string  `json:"item_category,omitempty"`
	PaymentMethod string  `json:"payment_method,omitempty"` // For Duitku: BC, M2, I1, etc.
	CallbackURL   string  `json:"callback_url,omitempty"`
	ReturnURL     string  `json:"return_url,omitempty"`
}

// CreateTransactionResponse holds the payment gateway response
type CreateTransactionResponse struct {
	OrderID    string     `json:"order_id"`
	SnapToken  string     `json:"snap_token,omitempty"`  // For Midtrans Snap
	PaymentURL string     `json:"payment_url,omitempty"` // Redirect URL
	ExpiredAt  *time.Time `json:"expired_at,omitempty"`
	RawResponse interface{} `json:"raw_response,omitempty"`
}

// NotificationResult holds the processed notification data
type NotificationResult struct {
	OrderID           string     `json:"order_id"`
	TransactionID     string     `json:"transaction_id,omitempty"`
	TransactionStatus string     `json:"transaction_status"`
	PaymentType       string     `json:"payment_type,omitempty"`
	FraudStatus       string     `json:"fraud_status,omitempty"`
	GrossAmount       float64    `json:"gross_amount,omitempty"`
	TransactionTime   *time.Time `json:"transaction_time,omitempty"`
	SettlementTime    *time.Time `json:"settlement_time,omitempty"`
	
	// Derived status for our system
	IsSuccess bool `json:"is_success"`
	IsPending bool `json:"is_pending"`
	IsFailed  bool `json:"is_failed"`
}

// TransactionStatus holds the current transaction status from gateway
type TransactionStatus struct {
	OrderID           string  `json:"order_id"`
	TransactionID     string  `json:"transaction_id,omitempty"`
	TransactionStatus string  `json:"transaction_status"`
	PaymentType       string  `json:"payment_type,omitempty"`
	FraudStatus       string  `json:"fraud_status,omitempty"`
	GrossAmount       float64 `json:"gross_amount,omitempty"`
}

// ProviderFactory creates payment providers based on name
type ProviderFactory struct {
	providers map[string]PaymentProvider
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]PaymentProvider),
	}
}

// Register adds a provider to the factory
func (f *ProviderFactory) Register(provider PaymentProvider) {
	f.providers[provider.GetName()] = provider
}

// Get retrieves a provider by name
func (f *ProviderFactory) Get(name string) PaymentProvider {
	return f.providers[name]
}

// TransactionStatusMapper maps gateway status to our internal status
func MapMidtransStatus(status, fraudStatus string) string {
	switch status {
	case "capture":
		if fraudStatus == "accept" {
			return "settlement"
		}
		return "capture"
	case "settlement":
		return "settlement"
	case "pending":
		return "pending"
	case "deny":
		return "deny"
	case "cancel":
		return "cancel"
	case "expire":
		return "expire"
	case "failure":
		return "failure"
	case "refund", "partial_refund":
		return "refund"
	default:
		return status
	}
}

// IsSuccessStatus checks if the status indicates successful payment
func IsSuccessStatus(status string) bool {
	return status == "settlement" || status == "capture"
}

// IsPendingStatus checks if the status indicates pending payment
func IsPendingStatus(status string) bool {
	return status == "pending"
}

// IsFailedStatus checks if the status indicates failed payment
func IsFailedStatus(status string) bool {
	return status == "deny" || status == "cancel" || status == "expire" || status == "failure"
}
