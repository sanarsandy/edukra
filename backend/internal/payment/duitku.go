package payment

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DuitkuConfig holds configuration for Duitku payment provider
type DuitkuConfig struct {
	MerchantCode string
	MerchantKey  string
	IsProduction bool
}

// DuitkuProvider implements PaymentProvider for Duitku
type DuitkuProvider struct {
	merchantCode string
	merchantKey  string
	isProduction bool
	httpClient   *http.Client
}

// NewDuitkuProvider creates a new Duitku payment provider
func NewDuitkuProvider(config DuitkuConfig) *DuitkuProvider {
	return &DuitkuProvider{
		merchantCode: config.MerchantCode,
		merchantKey:  config.MerchantKey,
		isProduction: config.IsProduction,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName returns the provider name
func (p *DuitkuProvider) GetName() string {
	return "duitku"
}

// getBaseURL returns the base API URL based on environment (API v2)
func (p *DuitkuProvider) getBaseURL() string {
	if p.isProduction {
		return "https://passport.duitku.com/webapi/api/merchant"
	}
	return "https://sandbox.duitku.com/webapi/api/merchant"
}

// generatePopSignature creates SHA256 signature for Duitku Pop API
// Format: SHA256(merchantCode + timestamp + merchantKey)
func (p *DuitkuProvider) generatePopSignature(timestamp int64) string {
	signatureString := fmt.Sprintf("%s%d%s",
		p.merchantCode,
		timestamp,
		p.merchantKey,
	)
	hash := sha256.Sum256([]byte(signatureString))
	return hex.EncodeToString(hash[:])
}

// generateSignatureForCreate creates MD5 signature for Duitku transaction creation
// Format: MD5(merchantCode + merchantOrderId + amount + merchantKey)
func (p *DuitkuProvider) generateSignatureForCreate(merchantOrderID string, amount int64) string {
	signatureString := fmt.Sprintf("%s%s%d%s",
		p.merchantCode,
		merchantOrderID,
		amount,
		p.merchantKey,
	)
	hash := md5.Sum([]byte(signatureString))
	return hex.EncodeToString(hash[:])
}

// generateSignatureForCallback creates MD5 signature for Duitku callback verification
// Format: MD5(merchantCode + amount + merchantOrderId + merchantKey)
func (p *DuitkuProvider) generateSignatureForCallback(merchantOrderID string, amount int64) string {
	signatureString := fmt.Sprintf("%s%d%s%s",
		p.merchantCode,
		amount,
		merchantOrderID,
		p.merchantKey,
	)
	hash := md5.Sum([]byte(signatureString))
	return hex.EncodeToString(hash[:])
}

// generateSignature is an alias for generateSignatureForCallback (for backward compatibility)
func (p *DuitkuProvider) generateSignature(merchantOrderID string, amount int64) string {
	return p.generateSignatureForCallback(merchantOrderID, amount)
}

// DuitkuInquiryRequest represents the transaction request to Duitku API v2
type DuitkuInquiryRequest struct {
	MerchantCode     string `json:"merchantCode"`
	PaymentAmount    int64  `json:"paymentAmount"`
	PaymentMethod    string `json:"paymentMethod"` // Payment method code (e.g., "SQ" for Shopeepay QRIS)
	MerchantOrderID  string `json:"merchantOrderId"`
	ProductDetails   string `json:"productDetails"`
	Email            string `json:"email"`
	CustomerVaName   string `json:"customerVaName"`
	CallbackURL      string `json:"callbackUrl"`
	ReturnURL        string `json:"returnUrl"`
	Signature        string `json:"signature"`
	ExpiryPeriod     int    `json:"expiryPeriod"` // in minutes
}

// DuitkuInquiryResponse represents response from Duitku API v2
type DuitkuInquiryResponse struct {
	MerchantCode  string `json:"merchantCode"`
	Reference     string `json:"reference"`
	PaymentURL    string `json:"paymentUrl"`
	VaNumber      string `json:"vaNumber,omitempty"`
	QrString      string `json:"qrString,omitempty"` // QR code string for QRIS
	Amount        string `json:"amount"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	Message       string `json:"Message"` // Error message field
}

// CreateTransaction creates a new payment transaction with Duitku API v2
func (p *DuitkuProvider) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	// Convert amount to integer (Duitku uses whole number)
	amount := int64(req.Amount)

	// Use payment method from frontend request, default to SQ (Shopeepay QRIS) if not specified
	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "SQ" // Default to Shopeepay QRIS
	}

	// Generate MD5 signature for API v2: md5(merchantCode + merchantOrderId + paymentAmount + apiKey)
	signature := p.generateSignatureForCreate(req.OrderID, amount)

	// Build request for API v2
	inquiryReq := DuitkuInquiryRequest{
		MerchantCode:    p.merchantCode,
		PaymentAmount:   amount,
		PaymentMethod:   paymentMethod,
		MerchantOrderID: req.OrderID,
		ProductDetails:  req.ItemName,
		Email:           req.CustomerEmail,
		CustomerVaName:  req.CustomerName,
		CallbackURL:     req.CallbackURL,
		ReturnURL:       req.ReturnURL,
		Signature:       signature,
		ExpiryPeriod:    1440, // 24 hours in minutes
	}

	// Marshal to JSON
	jsonBody, err := json.Marshal(inquiryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[Duitku v2] Request: %s", string(jsonBody))

	// Create HTTP request to API v2 inquiry endpoint
	url := p.getBaseURL() + "/v2/inquiry"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for API v2
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[Duitku v2] Response: %s", string(body))

	// Parse response
	var inquiryResp DuitkuInquiryResponse
	if err := json.Unmarshal(body, &inquiryResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	// Check for error
	if inquiryResp.Message != "" {
		return nil, fmt.Errorf("duitku error: %s", inquiryResp.Message)
	}
	if inquiryResp.StatusCode != "" && inquiryResp.StatusCode != "00" {
		return nil, fmt.Errorf("duitku error: %s - %s", inquiryResp.StatusCode, inquiryResp.StatusMessage)
	}

	// Set expiry time
	expiredAt := time.Now().Add(24 * time.Hour)

	return &CreateTransactionResponse{
		OrderID:    req.OrderID,
		PaymentURL: inquiryResp.PaymentURL,
		SnapToken:  inquiryResp.Reference,
		ExpiredAt:  &expiredAt,
	}, nil
}

// PaymentMethodInfo represents a single payment method from Duitku
type PaymentMethodInfo struct {
	PaymentMethod string `json:"paymentMethod"`
	PaymentName   string `json:"paymentName"`
	PaymentImage  string `json:"paymentImage"`
	TotalFee      string `json:"totalFee"`
}

// GetPaymentMethodsResponse represents the response from Get Payment Methods API
type GetPaymentMethodsResponse struct {
	PaymentFee      []PaymentMethodInfo `json:"paymentFee"`
	ResponseCode    string              `json:"responseCode"`
	ResponseMessage string              `json:"responseMessage"`
	Message         string              `json:"Message,omitempty"` // Error message
}

// GetPaymentMethods fetches available payment methods from Duitku
func (p *DuitkuProvider) GetPaymentMethods(ctx context.Context, amount int64) ([]PaymentMethodInfo, error) {
	// Format datetime: yyyy-MM-dd HH:mm:ss
	datetime := time.Now().Format("2006-01-02 15:04:05")

	// Generate SHA256 signature: sha256(merchantCode + amount + datetime + apiKey)
	signatureString := fmt.Sprintf("%s%d%s%s",
		p.merchantCode,
		amount,
		datetime,
		p.merchantKey,
	)
	hash := sha256.Sum256([]byte(signatureString))
	signature := hex.EncodeToString(hash[:])

	// Build request
	reqBody := map[string]interface{}{
		"merchantcode": p.merchantCode,
		"amount":       amount,
		"datetime":     datetime,
		"signature":    signature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[Duitku GetPaymentMethods] Request: %s", string(jsonBody))

	// Create HTTP request
	url := p.getBaseURL() + "/paymentmethod/getpaymentmethod"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[Duitku GetPaymentMethods] Response: %s", string(body))

	// Parse response
	var pmResp GetPaymentMethodsResponse
	if err := json.Unmarshal(body, &pmResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	// Check for error
	if pmResp.Message != "" {
		return nil, fmt.Errorf("duitku error: %s", pmResp.Message)
	}
	if pmResp.ResponseCode != "00" {
		return nil, fmt.Errorf("duitku error: %s - %s", pmResp.ResponseCode, pmResp.ResponseMessage)
	}

	return pmResp.PaymentFee, nil
}

// DuitkuCallbackData represents callback data from Duitku
type DuitkuCallbackData struct {
	MerchantCode    string `json:"merchantCode" form:"merchantCode"`
	Amount          string `json:"amount" form:"amount"`
	MerchantOrderID string `json:"merchantOrderId" form:"merchantOrderId"`
	ProductDetail   string `json:"productDetail" form:"productDetail"`
	PaymentCode     string `json:"paymentCode" form:"paymentCode"`
	ResultCode      string `json:"resultCode" form:"resultCode"`
	Reference       string `json:"reference" form:"reference"`
	Signature       string `json:"signature" form:"signature"`
}

// HandleNotification processes Duitku callback notifications
func (p *DuitkuProvider) HandleNotification(ctx context.Context, body []byte) (*NotificationResult, error) {
	// Parse callback data
	var callback DuitkuCallbackData
	if err := json.Unmarshal(body, &callback); err != nil {
		// Try parsing as Form Data (fallback)
		values, errForm := url.ParseQuery(string(body))
		if errForm == nil && values.Get("merchantCode") != "" {
			callback.MerchantCode = values.Get("merchantCode")
			callback.Amount = values.Get("amount")
			callback.MerchantOrderID = values.Get("merchantOrderId")
			callback.ProductDetail = values.Get("productDetail")
			callback.PaymentCode = values.Get("paymentCode")
			callback.ResultCode = values.Get("resultCode")
			callback.Reference = values.Get("reference")
			callback.Signature = values.Get("signature")
		} else {
			return nil, fmt.Errorf("failed to parse callback: %w", err)
		}
	}

	// Verify signature
	if !p.verifySignatureFromBytes(body) {
		return nil, fmt.Errorf("invalid callback signature")
	}

	// Determine status based on resultCode
	// 00 = Success, 01 = Pending, 02 = Failed/Canceled
	var status string
	var isSuccess bool
	var isPending bool

	switch callback.ResultCode {
	case "00":
		status = "settlement"
		isSuccess = true
	case "01":
		status = "pending"
		isPending = true
	default:
		status = "failure"
	}

	// Parse amount
	amount, _ := strconv.ParseFloat(callback.Amount, 64)

	return &NotificationResult{
		OrderID:           callback.MerchantOrderID,
		TransactionID:     callback.Reference,
		TransactionStatus: status,
		PaymentType:       callback.PaymentCode,
		GrossAmount:       amount,
		FraudStatus:       "",
		IsSuccess:         isSuccess,
		IsPending:         isPending,
		TransactionTime:   nil,
		SettlementTime:    nil,
	}, nil
}

// VerifySignature verifies the callback signature
func (p *DuitkuProvider) VerifySignature(data map[string]interface{}) bool {
	merchantOrderID, _ := data["merchantOrderId"].(string)
	amountStr, _ := data["amount"].(string)
	signature, _ := data["signature"].(string)

	if merchantOrderID == "" || amountStr == "" || signature == "" {
		return false
	}

	// Parse amount as int64
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		// Try parsing as float and convert
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return false
		}
		amount = int64(amountFloat)
	}

	// Generate expected signature
	expectedSignature := p.generateSignature(merchantOrderID, amount)

	return signature == expectedSignature
}

// verifySignatureFromBytes verifies signature from raw bytes (JSON or form-urlencoded)
func (p *DuitkuProvider) verifySignatureFromBytes(body []byte) bool {
	// First try JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		return p.VerifySignature(data)
	}
	
	// If JSON fails, try form-urlencoded
	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("[Duitku] Failed to parse callback body: %v", err)
		return false
	}
	
	merchantOrderID := values.Get("merchantOrderId")
	amountStr := values.Get("amount")
	signature := values.Get("signature")
	
	if merchantOrderID == "" || amountStr == "" || signature == "" {
		log.Printf("[Duitku] Missing required fields: orderID=%s, amount=%s, sig=%s", merchantOrderID, amountStr, signature != "")
		return false
	}
	
	// Parse amount
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("[Duitku] Failed to parse amount: %s", amountStr)
			return false
		}
		amount = int64(amountFloat)
	}
	
	// Generate expected signature and compare
	expectedSignature := p.generateSignatureForCallback(merchantOrderID, amount)
	
	if signature != expectedSignature {
		log.Printf("[Duitku] Signature mismatch: received=%s, expected=%s", signature, expectedSignature)
		log.Printf("[Duitku] Signature input: merchantCode=%s, amount=%d, orderId=%s", p.merchantCode, amount, merchantOrderID)
	}
	
	return signature == expectedSignature
}

// GetTransactionStatus retrieves transaction status from Duitku
func (p *DuitkuProvider) GetTransactionStatus(ctx context.Context, orderID string) (*TransactionStatus, error) {
	// Generate signature for check transaction
	// signature = md5(merchantCode + merchantOrderId + merchantKey)
	signatureString := fmt.Sprintf("%s%s%s", p.merchantCode, orderID, p.merchantKey)
	hash := md5.Sum([]byte(signatureString))
	signature := hex.EncodeToString(hash[:])

	// Build request
	reqBody := map[string]string{
		"merchantCode":    p.merchantCode,
		"merchantOrderId": orderID,
		"signature":       signature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := p.getBaseURL() + "/transactionStatus"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var statusResp struct {
		MerchantOrderID string `json:"merchantOrderId"`
		Reference       string `json:"reference"`
		Amount          string `json:"amount"`
		StatusCode      string `json:"statusCode"`
		StatusMessage   string `json:"statusMessage"`
	}

	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Map status
	var status string
	switch statusResp.StatusCode {
	case "00":
		status = "settlement"
	case "01":
		status = "pending"
	default:
		status = "failure"
	}

	amount, _ := strconv.ParseFloat(statusResp.Amount, 64)

	return &TransactionStatus{
		OrderID:           statusResp.MerchantOrderID,
		TransactionID:     statusResp.Reference,
		TransactionStatus: status,
		GrossAmount:       amount,
	}, nil
}

// IsProduction returns whether provider is in production mode
func (p *DuitkuProvider) IsProduction() bool {
	return p.isProduction
}

// GetMerchantCode returns the merchant code (for frontend display)
func (p *DuitkuProvider) GetMerchantCode() string {
	return p.merchantCode
}
