package payment

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// MidtransProvider implements PaymentProvider for Midtrans
type MidtransProvider struct {
	serverKey    string
	clientKey    string
	isProduction bool
}

// MidtransConfig holds Midtrans configuration
type MidtransConfig struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

// NewMidtransProvider creates a new Midtrans provider
func NewMidtransProvider(cfg MidtransConfig) *MidtransProvider {
	return &MidtransProvider{
		serverKey:    cfg.ServerKey,
		clientKey:    cfg.ClientKey,
		isProduction: cfg.IsProduction,
	}
}

// GetName returns the provider name
func (m *MidtransProvider) GetName() string {
	return "midtrans"
}

// getBaseURL returns the appropriate base URL based on environment
func (m *MidtransProvider) getBaseURL() string {
	if m.isProduction {
		return "https://app.midtrans.com"
	}
	return "https://app.sandbox.midtrans.com"
}

// getAPIURL returns the appropriate API URL based on environment
func (m *MidtransProvider) getAPIURL() string {
	if m.isProduction {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

// CreateTransaction creates a payment transaction using Midtrans Snap
func (m *MidtransProvider) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	// Build Snap request payload
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": int64(req.Amount),
		},
		"customer_details": map[string]interface{}{
			"first_name": req.CustomerName,
			"email":      req.CustomerEmail,
		},
		"item_details": []map[string]interface{}{
			{
				"id":       req.ItemID,
				"price":    int64(req.Amount),
				"quantity": 1,
				"name":     truncateString(req.ItemName, 50),
				"category": req.ItemCategory,
			},
		},
	}

	// Add phone if provided
	if req.CustomerPhone != "" {
		payload["customer_details"].(map[string]interface{})["phone"] = req.CustomerPhone
	}

	// Add callbacks if provided
	if req.CallbackURL != "" || req.ReturnURL != "" {
		callbacks := make(map[string]string)
		if req.ReturnURL != "" {
			callbacks["finish"] = req.ReturnURL
			callbacks["unfinish"] = req.ReturnURL
			callbacks["error"] = req.ReturnURL
		}
		payload["callbacks"] = callbacks
	}

	// Marshal payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request to Snap API
	snapURL := m.getBaseURL() + "/snap/v1/transactions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", snapURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(m.serverKey, "")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[Midtrans] Snap API response: %s", string(body))

	// Parse response
	var snapResp struct {
		Token       string   `json:"token"`
		RedirectURL string   `json:"redirect_url"`
		ErrorMessages []string `json:"error_messages,omitempty"`
	}
	if err := json.Unmarshal(body, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(snapResp.ErrorMessages) > 0 {
		return nil, fmt.Errorf("midtrans error: %v", snapResp.ErrorMessages)
	}

	if snapResp.Token == "" {
		return nil, fmt.Errorf("no snap token in response")
	}

	// Snap token expires in 24 hours
	expiredAt := time.Now().Add(24 * time.Hour)

	return &CreateTransactionResponse{
		OrderID:    req.OrderID,
		SnapToken:  snapResp.Token,
		PaymentURL: snapResp.RedirectURL,
		ExpiredAt:  &expiredAt,
		RawResponse: snapResp,
	}, nil
}

// HandleNotification processes Midtrans webhook notification
func (m *MidtransProvider) HandleNotification(ctx context.Context, data []byte) (*NotificationResult, error) {
	var notification map[string]interface{}
	if err := json.Unmarshal(data, &notification); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	// Verify signature
	if !m.VerifySignature(notification) {
		return nil, fmt.Errorf("invalid signature")
	}

	orderID, _ := notification["order_id"].(string)
	transactionID, _ := notification["transaction_id"].(string)
	transactionStatus, _ := notification["transaction_status"].(string)
	paymentType, _ := notification["payment_type"].(string)
	fraudStatus, _ := notification["fraud_status"].(string)
	grossAmountStr, _ := notification["gross_amount"].(string)
	transactionTimeStr, _ := notification["transaction_time"].(string)
	settlementTimeStr, _ := notification["settlement_time"].(string)

	// Parse gross amount
	var grossAmount float64
	if grossAmountStr != "" {
		grossAmount, _ = strconv.ParseFloat(grossAmountStr, 64)
	}

	// Parse transaction time
	var transactionTime *time.Time
	if transactionTimeStr != "" {
		t, err := time.Parse("2006-01-02 15:04:05", transactionTimeStr)
		if err == nil {
			transactionTime = &t
		}
	}

	// Parse settlement time
	var settlementTime *time.Time
	if settlementTimeStr != "" {
		t, err := time.Parse("2006-01-02 15:04:05", settlementTimeStr)
		if err == nil {
			settlementTime = &t
		}
	}

	// Map status
	mappedStatus := MapMidtransStatus(transactionStatus, fraudStatus)

	return &NotificationResult{
		OrderID:           orderID,
		TransactionID:     transactionID,
		TransactionStatus: mappedStatus,
		PaymentType:       paymentType,
		FraudStatus:       fraudStatus,
		GrossAmount:       grossAmount,
		TransactionTime:   transactionTime,
		SettlementTime:    settlementTime,
		IsSuccess:         IsSuccessStatus(mappedStatus),
		IsPending:         IsPendingStatus(mappedStatus),
		IsFailed:          IsFailedStatus(mappedStatus),
	}, nil
}

// VerifySignature verifies Midtrans notification signature
func (m *MidtransProvider) VerifySignature(data map[string]interface{}) bool {
	orderID, _ := data["order_id"].(string)
	statusCode, _ := data["status_code"].(string)
	grossAmount, _ := data["gross_amount"].(string)
	signatureKey, _ := data["signature_key"].(string)

	if orderID == "" || statusCode == "" || grossAmount == "" || signatureKey == "" {
		log.Printf("[Midtrans] Missing fields for signature verification")
		return false
	}

	// Calculate expected signature: SHA512(order_id + status_code + gross_amount + server_key)
	rawSignature := orderID + statusCode + grossAmount + m.serverKey
	hash := sha512.Sum512([]byte(rawSignature))
	expectedSignature := hex.EncodeToString(hash[:])

	if expectedSignature != signatureKey {
		log.Printf("[Midtrans] Signature mismatch: expected %s, got %s", expectedSignature, signatureKey)
		return false
	}

	return true
}

// GetTransactionStatus gets the current status of a transaction from Midtrans
func (m *MidtransProvider) GetTransactionStatus(ctx context.Context, orderID string) (*TransactionStatus, error) {
	url := m.getAPIURL() + "/v2/" + orderID + "/status"
	
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.SetBasicAuth(m.serverKey, "")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var statusResp struct {
		OrderID           string `json:"order_id"`
		TransactionID     string `json:"transaction_id"`
		TransactionStatus string `json:"transaction_status"`
		PaymentType       string `json:"payment_type"`
		FraudStatus       string `json:"fraud_status"`
		GrossAmount       string `json:"gross_amount"`
		StatusCode        string `json:"status_code"`
		StatusMessage     string `json:"status_message"`
	}

	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	grossAmount, _ := strconv.ParseFloat(statusResp.GrossAmount, 64)

	return &TransactionStatus{
		OrderID:           statusResp.OrderID,
		TransactionID:     statusResp.TransactionID,
		TransactionStatus: MapMidtransStatus(statusResp.TransactionStatus, statusResp.FraudStatus),
		PaymentType:       statusResp.PaymentType,
		FraudStatus:       statusResp.FraudStatus,
		GrossAmount:       grossAmount,
	}, nil
}

// GetClientKey returns the client key for frontend use
func (m *MidtransProvider) GetClientKey() string {
	return m.clientKey
}

// IsProduction returns whether this is production mode
func (m *MidtransProvider) IsProduction() bool {
	return m.isProduction
}

// truncateString truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
