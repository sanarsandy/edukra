package handlers

import (
	"fmt"
	"io"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/payment"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/service"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

var (
	paymentProvider        payment.PaymentProvider
	paymentTxRepo          *postgres.TransactionRepository
	courseRepoCheckout     *postgres.CourseRepository
	enrollmentRepoCheckout *postgres.EnrollmentRepository
)

func initPaymentRepos() {
	if paymentTxRepo == nil && db.DB != nil {
		paymentTxRepo = postgres.NewTransactionRepository(db.DB)
	}
	if courseRepoCheckout == nil && db.DB != nil {
		courseRepoCheckout = postgres.NewCourseRepository(db.DB)
	}
	if enrollmentRepoCheckout == nil && db.DB != nil {
		enrollmentRepoCheckout = postgres.NewEnrollmentRepository(db.DB)
	}
}

// InitPaymentProvider initializes the payment provider based on settings
func InitPaymentProvider() {
	// Get selected provider
	provider := getSettingValue("payment_provider", "midtrans")

	switch provider {
	case "duitku":
		initDuitkuProvider()
	default:
		initMidtransProvider()
	}
}

// initMidtransProvider initializes Midtrans payment provider
func initMidtransProvider() {
	serverKey := getSettingValue("payment_midtrans_server_key", "")
	clientKey := getSettingValue("payment_midtrans_client_key", "")
	isProduction := getSettingValue("payment_midtrans_is_production", "false") == "true"

	if serverKey == "" {
		log.Println("Midtrans not configured: server key is empty")
		return
	}

	paymentProvider = payment.NewMidtransProvider(payment.MidtransConfig{
		ServerKey:    serverKey,
		ClientKey:    clientKey,
		IsProduction: isProduction,
	})

	log.Printf("Midtrans payment provider initialized (production: %v)", isProduction)
}

// initDuitkuProvider initializes Duitku payment provider
func initDuitkuProvider() {
	merchantCode := getSettingValue("payment_duitku_merchant_code", "")
	merchantKey := getSettingValue("payment_duitku_merchant_key", "")
	isProduction := getSettingValue("payment_duitku_is_production", "false") == "true"

	if merchantCode == "" || merchantKey == "" {
		log.Println("Duitku not configured: merchant code or key is empty")
		return
	}

	paymentProvider = payment.NewDuitkuProvider(payment.DuitkuConfig{
		MerchantCode: merchantCode,
		MerchantKey:  merchantKey,
		IsProduction: isProduction,
	})

	log.Printf("Duitku payment provider initialized (production: %v)", isProduction)
}

// CheckoutRequest represents the checkout request body
type CheckoutRequest struct {
	CourseID      string `json:"course_id" validate:"required"`
	CouponCode    string `json:"coupon_code,omitempty"` // Optional coupon code
	ReturnURL     string `json:"return_url,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"` // Duitku payment method code (e.g., "BC", "M2", "I1")
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	TransactionID  string     `json:"transaction_id"`
	OrderID        string     `json:"order_id"`
	SnapToken      string     `json:"snap_token,omitempty"`
	PaymentURL     string     `json:"payment_url,omitempty"`
	ClientKey      string     `json:"client_key,omitempty"`
	ExpiredAt      *time.Time `json:"expired_at,omitempty"`
	IsFree         bool       `json:"is_free"`
	Message        string     `json:"message,omitempty"`
	OriginalAmount float64    `json:"original_amount,omitempty"`
	DiscountAmount float64    `json:"discount_amount,omitempty"`
	FinalAmount    float64    `json:"final_amount,omitempty"`
}

// CreateCheckout creates a new payment checkout
// POST /api/checkout
func CreateCheckout(c echo.Context) error {
	initPaymentRepos()

	// Check if payment is enabled
	if !getSettingBool("payment_enabled", false) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Payment module is not enabled",
		})
	}

	// Get user from JWT
	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Parse request
	var req CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.CourseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Course ID is required"})
	}

	// Get course
	course, err := courseRepoCheckout.GetByID(req.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}

	// Check if user is already enrolled
	enrolled, err := enrollmentRepoCheckout.IsEnrolled(userID, req.CourseID)
	if err == nil && enrolled {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "You are already enrolled in this course",
		})
	}

	// Get user info for customer details
	userRepo := postgres.NewUserRepository(db.DB)
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	// Check if course is free
	if course.Price == 0 {
		// Free course - enroll directly
		enrollment := &postgres.Enrollment{
			UserID:   userID,
			CourseID: req.CourseID,
		}
		if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
		}

		return c.JSON(http.StatusOK, CheckoutResponse{
			IsFree:  true,
			Message: "Successfully enrolled in free course",
		})
	}

	// === APPLY COURSE DISCOUNT FIRST ===
	var appliedCoupon *domain.Coupon
	var couponDiscountAmount float64 = 0
	originalPrice := course.Price
	
	// Check if course has an active discount
	priceAfterCourseDiscount := course.Price
	if course.DiscountPrice != nil && *course.DiscountPrice > 0 {
		// Check if discount is still valid
		discountValid := true
		if course.DiscountValidUntil != nil {
			discountValid = course.DiscountValidUntil.After(time.Now())
		}
		if discountValid {
			priceAfterCourseDiscount = *course.DiscountPrice
			log.Printf("[Checkout] Course discount applied: Original=%.2f, Discounted=%.2f", 
				originalPrice, priceAfterCourseDiscount)
		}
	}
	
	finalPrice := priceAfterCourseDiscount

	// === COUPON VALIDATION ===
	if req.CouponCode != "" {
		initCouponRepo()
		coupon, err := couponRepo.GetByCode(req.CouponCode)
		if err != nil {
			log.Printf("[Checkout] Failed to fetch coupon: %v", err)
		} else if coupon != nil {
			// Validate coupon for user and course
			valid, message := couponRepo.ValidateCouponForUser(coupon.ID, userID, req.CourseID)
			if !valid {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
			}
			
			// Calculate discount based on price AFTER course discount
			appliedCoupon = coupon
			couponDiscountAmount = coupon.CalculateDiscount(priceAfterCourseDiscount)
			finalPrice = priceAfterCourseDiscount - couponDiscountAmount
			
			log.Printf("[Checkout] Coupon %s applied: PriceAfterCourseDiscount=%.2f, CouponDiscount=%.2f, Final=%.2f", 
				coupon.Code, priceAfterCourseDiscount, couponDiscountAmount, finalPrice)
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Kode kupon tidak ditemukan"})
		}
	}

	// If discount makes it free, enroll directly
	if finalPrice <= 0 {
		// Record coupon usage first
		if appliedCoupon != nil {
			usage := &domain.CouponUsage{
				CouponID:        appliedCoupon.ID,
				UserID:          userID,
				DiscountApplied: couponDiscountAmount,
			}
			if err := couponRepo.RecordUsage(usage); err != nil {
				log.Printf("[Checkout] Failed to record coupon usage: %v", err)
			}
			couponRepo.IncrementUsage(appliedCoupon.ID)
		}
		
		// Enroll directly
		enrollment := &postgres.Enrollment{
			UserID:   userID,
			CourseID: req.CourseID,
		}
		if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
		}

		totalDiscount := originalPrice - finalPrice
		return c.JSON(http.StatusOK, CheckoutResponse{
			IsFree:         true,
			Message:        "Diskon berhasil diterapkan. Anda terdaftar gratis!",
			OriginalAmount: originalPrice,
			DiscountAmount: totalDiscount,
			FinalAmount:    0,
		})
	}

	// Cancel existing pending transaction for this user and course to ensure fresh checkout
	if err := paymentTxRepo.CancelPendingByUserAndCourse(userID, req.CourseID); err != nil {
		log.Printf("[Checkout] Failed to cancel existing transactions: %v", err)
		// Continue anyway to try creating new transaction
	}

	// Check if payment provider is initialized
	if paymentProvider == nil {
		InitPaymentProvider()
		if paymentProvider == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Payment provider not configured",
			})
		}
	}

	// Generate unique order ID
	orderID := fmt.Sprintf("LMS-%s-%d", uuid.New().String()[:8], time.Now().UnixMilli()%100000)

	// Create transaction record with dynamic provider name
	tx := &postgres.Transaction{
		UserID:         userID,
		PaymentGateway: paymentProvider.GetName(),
		Amount:         finalPrice, // Use discounted price
		Currency:       course.Currency,
		Status:         "pending",
	}
	tx.CourseID = &req.CourseID
	tx.OrderID = &orderID
	// Store discount info in metadata (includes both course discount and coupon discount)
	totalDiscount := originalPrice - finalPrice
	if totalDiscount > 0 || appliedCoupon != nil {
		tx.OriginalAmount = &originalPrice
		tx.DiscountAmount = &totalDiscount
		if appliedCoupon != nil {
			tx.CouponID = &appliedCoupon.ID
		}
	}

	if err := paymentTxRepo.Create(tx); err != nil {
		log.Printf("[Checkout] Failed to create transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create transaction"})
	}

	// Create payment via provider - USE FINAL PRICE (after all discounts)
	ctx := c.Request().Context()
	paymentReq := &payment.CreateTransactionRequest{
		OrderID:       orderID,
		Amount:        finalPrice, // FIXED: Use final price after all discounts
		Currency:      course.Currency,
		CustomerName:  user.FullName,
		CustomerEmail: user.Email,
		ItemName:      course.Title,
		ItemID:        course.ID,
		ItemCategory:  "course",
		PaymentMethod: req.PaymentMethod, // For Duitku payment method selection
		ReturnURL:     req.ReturnURL,
		CallbackURL:   fmt.Sprintf("%s://%s/api/webhooks/duitku", c.Scheme(), c.Request().Host),
	}
	
	log.Printf("[Checkout] Creating payment: OrderID=%s, OriginalPrice=%.2f, FinalPrice=%.2f", 
		orderID, originalPrice, finalPrice)

	paymentResp, err := paymentProvider.CreateTransaction(ctx, paymentReq)
	if err != nil {
		log.Printf("[Checkout] Failed to create Midtrans transaction: %v", err)
		paymentTxRepo.UpdateStatus(tx.ID, "failure")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create payment: " + err.Error(),
		})
	}

	// Update transaction with snap token
	if paymentResp.ExpiredAt != nil {
		paymentTxRepo.UpdateSnapToken(tx.ID, paymentResp.SnapToken, paymentResp.PaymentURL, *paymentResp.ExpiredAt)
	}

	// Get client key for frontend
	midtransProvider, ok := paymentProvider.(*payment.MidtransProvider)
	var clientKey string
	if ok {
		clientKey = midtransProvider.GetClientKey()
	}

	return c.JSON(http.StatusOK, CheckoutResponse{
		TransactionID: tx.ID,
		OrderID:       orderID,
		SnapToken:     paymentResp.SnapToken,
		PaymentURL:    paymentResp.PaymentURL,
		ClientKey:     clientKey,
		ExpiredAt:     paymentResp.ExpiredAt,
		IsFree:        false,
	})
}

// MidtransWebhook handles Midtrans payment notifications
// POST /api/webhooks/midtrans
func MidtransWebhook(c echo.Context) error {
	initPaymentRepos()

	// Read request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("[Midtrans Webhook] Failed to read body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
	}

	log.Printf("[Midtrans Webhook] Received: %s", string(body))

	// Initialize provider if needed
	if paymentProvider == nil {
		InitPaymentProvider()
	}

	if paymentProvider == nil {
		log.Printf("[Midtrans Webhook] Payment provider not configured")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Provider not configured"})
	}

	// Process notification
	ctx := c.Request().Context()
	result, err := paymentProvider.HandleNotification(ctx, body)
	if err != nil {
		log.Printf("[Midtrans Webhook] Failed to handle notification: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("[Midtrans Webhook] Processed: order_id=%s, status=%s, is_success=%v",
		result.OrderID, result.TransactionStatus, result.IsSuccess)

	// Get transaction by order ID
	tx, err := paymentTxRepo.GetByOrderID(result.OrderID)
	if err != nil || tx == nil {
		log.Printf("[Midtrans Webhook] Transaction not found: %s", result.OrderID)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}

	// Update transaction
	paymentType := result.PaymentType
	fraudStatus := result.FraudStatus
	transactionID := result.TransactionID
	
	err = paymentTxRepo.UpdateFromCallback(
		result.OrderID,
		result.TransactionStatus,
		&paymentType,
		&fraudStatus,
		result.TransactionTime,
		result.SettlementTime,
		&transactionID,
	)
	if err != nil {
		log.Printf("[Midtrans Webhook] Failed to update transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update transaction"})
	}

	// If payment successful, create enrollment and record coupon usage
	if result.IsSuccess && tx.CourseID != nil {
		// Check if already enrolled (idempotency)
		enrolled, _ := enrollmentRepoCheckout.IsEnrolled(tx.UserID, *tx.CourseID)
		if !enrolled {
			enrollment := &postgres.Enrollment{
				UserID:        tx.UserID,
				CourseID:      *tx.CourseID,
				TransactionID: &tx.ID,
			}
			if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
				log.Printf("[Midtrans Webhook] Failed to create enrollment: %v", err)
			} else {
				log.Printf("[Midtrans Webhook] Enrollment created for user %s, course %s", tx.UserID, *tx.CourseID)
			}
		}
		
		// Record coupon usage if a coupon was applied
		if tx.CouponID != nil && tx.DiscountAmount != nil {
			initCouponRepo()
			usage := &domain.CouponUsage{
				CouponID:        *tx.CouponID,
				UserID:          tx.UserID,
				TransactionID:   &tx.ID,
				DiscountApplied: *tx.DiscountAmount,
			}
			if err := couponRepo.RecordUsage(usage); err != nil {
				log.Printf("[Midtrans Webhook] Failed to record coupon usage: %v", err)
			} else {
				couponRepo.IncrementUsage(*tx.CouponID)
				log.Printf("[Midtrans Webhook] Coupon usage recorded for coupon %s", *tx.CouponID)
			}
		}
		
		// Send WhatsApp notification (async)
		go handlePaymentSuccessNotification(tx.UserID, *tx.CourseID)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// DuitkuWebhook handles Duitku payment callback notifications
// POST /api/webhooks/duitku
func DuitkuWebhook(c echo.Context) error {
	initPaymentRepos()

	// Read request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("[Duitku Webhook] Failed to read body: %v", err)
		return c.String(http.StatusBadRequest, "Failed to read body")
	}

	log.Printf("[Duitku Webhook] Received: %s", string(body))

	// Initialize provider if needed
	if paymentProvider == nil {
		InitPaymentProvider()
	}

	if paymentProvider == nil || paymentProvider.GetName() != "duitku" {
		log.Printf("[Duitku Webhook] Payment provider not configured for Duitku")
		return c.String(http.StatusServiceUnavailable, "Provider not configured")
	}

	// Process notification
	ctx := c.Request().Context()
	result, err := paymentProvider.HandleNotification(ctx, body)
	if err != nil {
		log.Printf("[Duitku Webhook] Failed to handle notification: %v", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("[Duitku Webhook] Processed: order_id=%s, status=%s, is_success=%v",
		result.OrderID, result.TransactionStatus, result.IsSuccess)

	// Get transaction by order ID
	tx, err := paymentTxRepo.GetByOrderID(result.OrderID)
	if err != nil || tx == nil {
		log.Printf("[Duitku Webhook] Transaction not found: %s", result.OrderID)
		return c.String(http.StatusNotFound, "Transaction not found")
	}

	// Update transaction
	paymentType := result.PaymentType
	transactionID := result.TransactionID

	err = paymentTxRepo.UpdateFromCallback(
		result.OrderID,
		result.TransactionStatus,
		&paymentType,
		nil, // no fraud status for Duitku
		result.TransactionTime,
		result.SettlementTime,
		&transactionID,
	)
	if err != nil {
		log.Printf("[Duitku Webhook] Failed to update transaction: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to update transaction")
	}

	// If payment successful
	if result.IsSuccess {
		// Check for webinar_only via metadata
		var metadata map[string]string
		if len(tx.Metadata) > 0 {
			json.Unmarshal(tx.Metadata, &metadata)
		}
		
		webinarID := metadata["webinar_id"]

		if tx.CourseID != nil {
			// Normal course enrollment logic
			enrolled, _ := enrollmentRepoCheckout.IsEnrolled(tx.UserID, *tx.CourseID)
			if !enrolled {
				enrollment := &postgres.Enrollment{
					UserID:        tx.UserID,
					CourseID:      *tx.CourseID,
					TransactionID: &tx.ID,
				}
				if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
					log.Printf("[Duitku Webhook] Failed to create enrollment: %v", err)
				} else {
					log.Printf("[Duitku Webhook] Enrollment created for user %s, course %s", tx.UserID, *tx.CourseID)
					
					// Send payment success notification and handle webinar registration
					go handlePaymentSuccessNotification(tx.UserID, *tx.CourseID)
				}
			}
		} else if webinarID != "" {
			// Webinar Only logic (no course enrollment)
			log.Printf("[Duitku Webhook] Webinar Only payment for webinar %s", webinarID)
			go handleWebinarPaymentSuccess(tx.UserID, webinarID)
		}
		
		// Record coupon usage if a coupon was applied
		if tx.CouponID != nil && tx.DiscountAmount != nil {
			initCouponRepo()
			usage := &domain.CouponUsage{
				CouponID:        *tx.CouponID,
				UserID:          tx.UserID,
				TransactionID:   &tx.ID,
				DiscountApplied: *tx.DiscountAmount,
			}
			if err := couponRepo.RecordUsage(usage); err != nil {
				log.Printf("[Duitku Webhook] Failed to record coupon usage: %v", err)
			} else {
				couponRepo.IncrementUsage(*tx.CouponID)
				log.Printf("[Duitku Webhook] Coupon usage recorded for coupon %s", *tx.CouponID)
			}
		}
	}

	// Duitku expects "SUCCESS" response
	return c.String(http.StatusOK, "SUCCESS")
}

// handlePaymentSuccessNotification handles post-payment notifications (Webinar or General)
func handlePaymentSuccessNotification(userID, courseID string) {
	webinarRepo := postgres.NewWebinarRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)
	courseRepo := postgres.NewCourseRepository(db.DB)

	// Get user info
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		log.Printf("[Notification] Failed to get user %s: %v", userID, err)
		return
	}
	
	log.Printf("[Notification] Processing for user %s (%s), Phone: %v", user.ID, user.FullName, user.Phone)
	if user.Phone != nil {
		log.Printf("[Notification] Phone value: '%s'", *user.Phone)
	} else {
		log.Printf("[Notification] Phone is NIL")
	}

	// Get LMS URL
	lmsURL := getSettingValue("frontend_url", "")
	if lmsURL == "" {
		lmsURL = os.Getenv("FRONTEND_URL")
	}
	if lmsURL == "" {
		lmsURL = "https://lms.edukra.com"
	}

	// 1. Check for upcoming webinars
	webinars, err := webinarRepo.GetUpcomingByCourse(courseID)
	if err != nil {
		log.Printf("[Notification] Failed to check webinars for course %s: %v", courseID, err)
	}
	log.Printf("[Notification] Found %d upcoming webinars for course %s", len(webinars), courseID)

	// CASE A: Webinar Course (Has upcoming webinars)
	if len(webinars) > 0 {
		// Register user to each webinar and create reminders
		for _, webinar := range webinars {
			if err := webinarRepo.RegisterUser(webinar.ID, userID, "purchase"); err != nil {
				log.Printf("[Webinar] Failed to register user %s to webinar %s: %v", userID, webinar.ID, err)
				continue
			}
			log.Printf("[Webinar] User %s registered to webinar %s", userID, webinar.ID)
			webinarRepo.CreateReminders(webinar.ID, userID, webinar.ScheduledAt)
		}

		// Send Webinar Confirmation WA (includes meeting details)
		if user.Phone != nil && *user.Phone != "" {
			webinar := webinars[0]
			data := service.WebinarConfirmationData{
				UserName:        user.FullName,
				UserEmail:       user.Email,
				WebinarTitle:    webinar.Title,
				WebinarDate:     webinar.ScheduledAt.Format("02 January 2006"),
				WebinarTime:     webinar.ScheduledAt.Format("15:04"),
				DurationMinutes: webinar.DurationMinutes,
				LMSUrl:          lmsURL,
			}
			
			if webinar.MeetingURL != nil {
				data.MeetingURL = *webinar.MeetingURL
			}
			if webinar.MeetingPassword != nil {
				data.MeetingPassword = *webinar.MeetingPassword
			}
			
			service.SendWebinarConfirmationAsync(*user.Phone, data)
			log.Printf("[Notification] Webinar confirmation sent to %s", *user.Phone)
		}
		return
	}

	// CASE B: Regular Course (No webinars) -> Send General Payment Success
	if user.Phone != nil && *user.Phone != "" {
		course, err := courseRepo.GetByID(courseID)
		if err != nil || course == nil {
			log.Printf("[Notification] Failed to get course %s: %v", courseID, err)
			return
		}

		go func() {
			if err := service.GetWhatsAppService().SendPaymentSuccess(*user.Phone, user.FullName, course.Title, lmsURL); err != nil {
				log.Printf("[Notification] Failed to send payment success WA to %s: %v", *user.Phone, err)
			} else {
				log.Printf("[Notification] General payment success WA sent to %s", *user.Phone)
			}
		}()
	}
}

// SimulatePaymentSuccess simulates a successful payment callback (FOR TESTING ONLY)
// POST /api/test/simulate-payment
func SimulatePaymentSuccess(c echo.Context) error {
	initPaymentRepos()

	var req struct {
		OrderID string `json:"order_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.OrderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "order_id is required"})
	}

	log.Printf("[SimulatePayment] Simulating payment for order: %s", req.OrderID)

	// Get transaction by order ID
	tx, err := paymentTxRepo.GetByOrderID(req.OrderID)
	if err != nil || tx == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}

	if tx.Status == "success" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Transaction already successful"})
	}

	// Update transaction to success
	now := time.Now()
	err = paymentTxRepo.UpdateFromCallback(
		req.OrderID,
		"settlement", // Duitku success status
		nil,          // payment_type
		nil,          // fraud_status
		&now,         // transaction_time
		&now,         // settlement_time
		nil,          // gateway_transaction_id
	)
	if err != nil {
		log.Printf("[SimulatePayment] Failed to update transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update transaction"})
	}

	// Create enrollment if not exists
	if tx.CourseID != nil {
		enrolled, _ := enrollmentRepoCheckout.IsEnrolled(tx.UserID, *tx.CourseID)
		if !enrolled {
			enrollment := &postgres.Enrollment{
				UserID:        tx.UserID,
				CourseID:      *tx.CourseID,
				TransactionID: &tx.ID,
			}
			if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
				log.Printf("[SimulatePayment] Failed to create enrollment: %v", err)
			} else {
				log.Printf("[SimulatePayment] Enrollment created for user %s, course %s", tx.UserID, *tx.CourseID)
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Payment simulated successfully",
		"order_id": req.OrderID,
		"status":   "success",
	})
}

// GetMyTransactions returns current user's transactions
// GET /api/my/transactions
func GetMyTransactions(c echo.Context) error {
	initPaymentRepos()

	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Parse pagination
	limit := 20
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.QueryParam("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	transactions, err := paymentTxRepo.ListByUser(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch transactions"})
	}

	// Enrich with course info
	type TransactionWithCourse struct {
		*postgres.Transaction
		Course *domain.Course `json:"course,omitempty"`
	}

	var enriched []TransactionWithCourse
	for _, t := range transactions {
		tc := TransactionWithCourse{Transaction: t}
		if t.CourseID != nil {
			course, _ := courseRepoCheckout.GetByID(*t.CourseID)
			tc.Course = course
		}
		enriched = append(enriched, tc)
	}

	total, _ := paymentTxRepo.CountByUser(userID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": enriched,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetPaymentSettings returns payment configuration for admin
// GET /api/admin/payment/settings
func GetPaymentSettings(c echo.Context) error {
	settings := map[string]interface{}{
		"enabled":                getSettingBool("payment_enabled", false),
		"provider":               getSettingValue("payment_provider", "midtrans"),
		// Midtrans settings
		"midtrans_client_key":    getSettingValue("payment_midtrans_client_key", ""),
		"midtrans_server_key":    maskString(getSettingValue("payment_midtrans_server_key", "")),
		"midtrans_is_production": getSettingValue("payment_midtrans_is_production", "false") == "true",
		// Duitku settings
		"duitku_merchant_code":   getSettingValue("payment_duitku_merchant_code", ""),
		"duitku_merchant_key":    maskString(getSettingValue("payment_duitku_merchant_key", "")),
		"duitku_is_production":   getSettingValue("payment_duitku_is_production", "false") == "true",
	}

	return c.JSON(http.StatusOK, settings)
}

// UpdatePaymentSettings updates payment configuration
// PUT /api/admin/payment/settings
func UpdatePaymentSettings(c echo.Context) error {
	var req struct {
		Enabled               *bool  `json:"enabled"`
		Provider              string `json:"provider"`
		// Midtrans
		MidtransServerKey     string `json:"midtrans_server_key,omitempty"`
		MidtransClientKey     string `json:"midtrans_client_key,omitempty"`
		MidtransIsProduction  *bool  `json:"midtrans_is_production,omitempty"`
		// Duitku
		DuitkuMerchantCode    string `json:"duitku_merchant_code,omitempty"`
		DuitkuMerchantKey     string `json:"duitku_merchant_key,omitempty"`
		DuitkuIsProduction    *bool  `json:"duitku_is_production,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update settings
	if req.Enabled != nil {
		setSettingValue("payment_enabled", fmt.Sprintf("%v", *req.Enabled))
	}
	if req.Provider != "" {
		setSettingValue("payment_provider", req.Provider)
	}

	// Midtrans settings
	if req.MidtransServerKey != "" && !isMasked(req.MidtransServerKey) {
		setSettingValue("payment_midtrans_server_key", req.MidtransServerKey)
	}
	if req.MidtransClientKey != "" {
		setSettingValue("payment_midtrans_client_key", req.MidtransClientKey)
	}
	if req.MidtransIsProduction != nil {
		setSettingValue("payment_midtrans_is_production", fmt.Sprintf("%v", *req.MidtransIsProduction))
	}

	// Duitku settings
	if req.DuitkuMerchantCode != "" {
		setSettingValue("payment_duitku_merchant_code", req.DuitkuMerchantCode)
	}
	if req.DuitkuMerchantKey != "" && !isMasked(req.DuitkuMerchantKey) {
		setSettingValue("payment_duitku_merchant_key", req.DuitkuMerchantKey)
	}
	if req.DuitkuIsProduction != nil {
		setSettingValue("payment_duitku_is_production", fmt.Sprintf("%v", *req.DuitkuIsProduction))
	}

	// Reinitialize payment provider
	InitPaymentProvider()

	return c.JSON(http.StatusOK, map[string]string{"message": "Payment settings updated"})
}

// isMasked checks if string is a masked value
func isMasked(s string) bool {
	return strings.Contains(s, "****")
}

// GetCheckoutConfig returns payment config for frontend (non-sensitive)
// GET /api/checkout/config
func GetCheckoutConfig(c echo.Context) error {
	enabled := getSettingBool("payment_enabled", false)
	providerName := getSettingValue("payment_provider", "midtrans")

	var clientKey string
	var isProduction bool

	if paymentProvider != nil {
		switch p := paymentProvider.(type) {
		case *payment.MidtransProvider:
			clientKey = p.GetClientKey()
			isProduction = p.IsProduction()
		case *payment.DuitkuProvider:
			clientKey = p.GetMerchantCode() // For Duitku, merchant code is used
			isProduction = p.IsProduction()
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"enabled":       enabled,
		"provider":      providerName,
		"client_key":    clientKey,
		"is_production": isProduction,
	})
}

// Helper to mask sensitive strings
func maskString(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// GetPaymentMethods returns available payment methods from payment provider
// GET /api/checkout/payment-methods?amount=10000
func GetPaymentMethods(c echo.Context) error {
	// Check if payment is enabled
	if !getSettingBool("payment_enabled", false) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Payment module is not enabled",
		})
	}

	// Get amount from query parameter
	amountStr := c.QueryParam("amount")
	if amountStr == "" {
		amountStr = "10000" // Default amount
	}
	
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid amount parameter",
		})
	}

	// Initialize provider if not already done
	if paymentProvider == nil {
		InitPaymentProvider()
	}

	// Only Duitku supports GetPaymentMethods
	providerName := getSettingValue("payment_provider", "midtrans")
	if providerName != "duitku" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"provider":         providerName,
			"payment_methods":  []interface{}{},
			"message":          "Payment methods list only available for Duitku provider",
		})
	}

	duitkuProvider, ok := paymentProvider.(*payment.DuitkuProvider)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get Duitku provider",
		})
	}

	// Fetch payment methods from Duitku
	methods, err := duitkuProvider.GetPaymentMethods(c.Request().Context(), amount)
	if err != nil {
		log.Printf("Failed to get payment methods: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch payment methods: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"provider":        providerName,
		"payment_methods": methods,
	})
}

func handleWebinarPaymentSuccess(userID, webinarID string) {
	// Initialize repos
	userRepo := postgres.NewUserRepository(db.DB)
	webinarRepo := postgres.NewWebinarRepository(db.DB)

	// 1. Get User
	user, err := userRepo.GetByID(userID)
	if err != nil {
		log.Printf("[WebinarOnlyPayment] Failed to get user: %v", err)
		return
	}

	// 2. Get Webinar
	webinar, err := webinarRepo.GetByID(webinarID)
	if err != nil {
		log.Printf("[WebinarOnlyPayment] Failed to get webinar: %v", err)
		return
	}

	// 3. Register user to webinar
	// Check if already registered
	registered, _ := webinarRepo.IsUserRegistered(webinarID, userID)
	if !registered {
		err = webinarRepo.RegisterUser(webinarID, userID, "payment_gateway")
		if err != nil {
			log.Printf("[WebinarOnlyPayment] Failed to register user to webinar: %v", err)
			// Continue to notification? Maybe better to retry or log error. 
			// Attempt notification anyway as manual fallback.
		} else {
			log.Printf("[WebinarOnlyPayment] User %s registered to webinar %s", userID, webinarID)
		}
	}

	// 4. Send WhatsApp Notification
	if user.Phone != nil && *user.Phone != "" {
		data := service.WebinarConfirmationData{
			UserName:        user.FullName,
			UserEmail:       user.Email,
			WebinarTitle:    webinar.Title,
			WebinarDate:     webinar.ScheduledAt.Format("02 January 2006"),
			WebinarTime:     webinar.ScheduledAt.Format("15:04"),
			DurationMinutes: webinar.DurationMinutes,
		}

		if webinar.MeetingURL != nil {
			data.MeetingURL = *webinar.MeetingURL
		}
		if webinar.MeetingPassword != nil {
			data.MeetingPassword = *webinar.MeetingPassword
		}

		service.SendWebinarOnlyConfirmationAsync(*user.Phone, data)
		log.Printf("[WebinarOnlyPayment] Confirmation sent to %s", *user.Phone)
	}
}
