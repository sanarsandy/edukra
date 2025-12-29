package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/payment"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	campaignRepoCheckout *postgres.CampaignRepository
)

func initCampaignCheckoutRepos() {
	if campaignRepoCheckout == nil && db.DB != nil {
		campaignRepoCheckout = postgres.NewCampaignRepository(db.DB)
	}
	initPaymentRepos()
}

// CampaignCheckoutRequest represents the campaign checkout request body
type CampaignCheckoutRequest struct {
	CampaignID    string `json:"campaign_id" validate:"required"`
	Email         string `json:"email" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	FullName      string `json:"full_name,omitempty"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	CouponCode    string `json:"coupon_code,omitempty"`
}

// CampaignCheckoutResponse represents the campaign checkout response
type CampaignCheckoutResponse struct {
	TransactionID  string     `json:"transaction_id"`
	OrderID        string     `json:"order_id"`
	PaymentURL     string     `json:"payment_url"`
	ExpiredAt      *time.Time `json:"expired_at,omitempty"`
	IsFree         bool       `json:"is_free"`
	Message        string     `json:"message,omitempty"`
	OriginalAmount float64    `json:"original_amount,omitempty"`
	DiscountAmount float64    `json:"discount_amount,omitempty"`
	FinalAmount    float64    `json:"final_amount,omitempty"`
	IsNewUser      bool       `json:"is_new_user"`
}

// validateEmail checks if email is valid
func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePhone checks if phone is valid (Indonesian format)
func validatePhone(phone string) bool {
	// Remove common prefixes and check length
	phone = strings.TrimPrefix(phone, "+62")
	phone = strings.TrimPrefix(phone, "62")
	phone = strings.TrimPrefix(phone, "0")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	// Check if it's a valid length (9-12 digits)
	if len(phone) < 9 || len(phone) > 12 {
		return false
	}
	
	// Check if all characters are digits
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// normalizePhone normalizes phone to +62 format
func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	if strings.HasPrefix(phone, "+62") {
		return phone
	}
	if strings.HasPrefix(phone, "62") {
		return "+" + phone
	}
	if strings.HasPrefix(phone, "0") {
		return "+62" + phone[1:]
	}
	return "+62" + phone
}

// generateRandomPassword generates a random password
func generateRandomPassword() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// findOrCreateGuestUser finds existing user by email or creates a new guest user
func findOrCreateGuestUser(email, phone, fullName string) (*domain.User, bool, error) {
	userRepo := postgres.NewUserRepository(db.DB)
	
	// Try to find existing user (use empty tenant for now)
	existingUser, err := userRepo.GetByEmail("", email)
	if err == nil && existingUser != nil {
		// Update phone if not set
		if existingUser.Phone == nil || *existingUser.Phone == "" {
			normalizedPhone := normalizePhone(phone)
			existingUser.Phone = &normalizedPhone
			userRepo.Update(existingUser)
		}
		return existingUser, false, nil
	}
	
	// Create new guest user
	password := generateRandomPassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, false, fmt.Errorf("failed to hash password: %v", err)
	}
	
	normalizedPhone := normalizePhone(phone)
	if fullName == "" {
		// Use email prefix as name
		parts := strings.Split(email, "@")
		fullName = parts[0]
	}
	
	newUser := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         domain.RoleStudent,
		AuthProvider: "guest", // Mark as guest registration
		Phone:        &normalizedPhone,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	if err := userRepo.Create(newUser); err != nil {
		return nil, false, fmt.Errorf("failed to create user: %v", err)
	}
	
	log.Printf("[CampaignCheckout] Created new guest user: %s (%s)", email, newUser.ID)
	return newUser, true, nil
}

// CampaignCheckout handles checkout from campaign pages for guest users
// POST /api/campaign-checkout
func CampaignCheckout(c echo.Context) error {
	initCampaignCheckoutRepos()

	// Check if payment is enabled
	if !getSettingBool("payment_enabled", false) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Payment module is not enabled",
		})
	}

	// Parse request
	var req CampaignCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.CampaignID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign ID is required"})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}
	if req.Phone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Phone number is required"})
	}
	if req.PaymentMethod == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payment method is required"})
	}

	// Validate email format
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !validateEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email format"})
	}

	// Validate phone format
	if !validatePhone(req.Phone) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid phone number format"})
	}

	// Get campaign
	campaign, err := campaignRepoCheckout.GetByID(req.CampaignID)
	if err != nil {
		log.Printf("[CampaignCheckout] Failed to fetch campaign: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch campaign"})
	}
	if campaign == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Campaign not found"})
	}
	if !campaign.IsActive {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign is not active"})
	}

	// Check if campaign has linked course
	if campaign.CourseID == nil || *campaign.CourseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign has no linked course"})
	}

	// Get course
	course, err := courseRepoCheckout.GetByID(*campaign.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}

	// Find or create user
	user, isNewUser, err := findOrCreateGuestUser(req.Email, req.Phone, req.FullName)
	if err != nil {
		log.Printf("[CampaignCheckout] Failed to find/create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process user"})
	}

	// Check if already enrolled
	enrolled, err := enrollmentRepoCheckout.IsEnrolled(user.ID, *campaign.CourseID)
	if err == nil && enrolled {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Anda sudah terdaftar di kursus ini",
		})
	}

	// Calculate price
	originalPrice := course.Price
	finalPrice := course.Price

	// Apply course discount if available
	if course.DiscountPrice != nil && *course.DiscountPrice > 0 {
		discountValid := true
		if course.DiscountValidUntil != nil {
			discountValid = course.DiscountValidUntil.After(time.Now())
		}
		if discountValid {
			finalPrice = *course.DiscountPrice
		}
	}

	// Apply coupon if provided
	var appliedCoupon *domain.Coupon
	var couponDiscountAmount float64 = 0
	
	if req.CouponCode != "" {
		initCouponRepo()
		coupon, err := couponRepo.GetByCode(req.CouponCode)
		if err != nil {
			log.Printf("[CampaignCheckout] Failed to fetch coupon: %v", err)
		} else if coupon != nil {
			valid, message := couponRepo.ValidateCouponForUser(coupon.ID, user.ID, *campaign.CourseID)
			if !valid {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
			}
			
			appliedCoupon = coupon
			couponDiscountAmount = coupon.CalculateDiscount(finalPrice)
			finalPrice = finalPrice - couponDiscountAmount
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Kode kupon tidak ditemukan"})
		}
	}

	// If free after discounts, enroll directly
	if finalPrice <= 0 {
		// Record coupon usage
		if appliedCoupon != nil {
			usage := &domain.CouponUsage{
				CouponID:        appliedCoupon.ID,
				UserID:          user.ID,
				DiscountApplied: couponDiscountAmount,
			}
			couponRepo.RecordUsage(usage)
			couponRepo.IncrementUsage(appliedCoupon.ID)
		}

		// Enroll directly
		enrollment := &postgres.Enrollment{
			UserID:   user.ID,
			CourseID: *campaign.CourseID,
		}
		if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
		}

		return c.JSON(http.StatusOK, CampaignCheckoutResponse{
			IsFree:         true,
			Message:        "Selamat! Anda berhasil terdaftar secara gratis.",
			OriginalAmount: originalPrice,
			DiscountAmount: originalPrice,
			FinalAmount:    0,
			IsNewUser:      isNewUser,
		})
	}

	// Check for free course
	if course.Price == 0 {
		enrollment := &postgres.Enrollment{
			UserID:   user.ID,
			CourseID: *campaign.CourseID,
		}
		if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
		}

		return c.JSON(http.StatusOK, CampaignCheckoutResponse{
			IsFree:    true,
			Message:   "Selamat! Anda berhasil terdaftar di kursus gratis ini.",
			IsNewUser: isNewUser,
		})
	}

	// Cancel existing pending transactions
	if err := paymentTxRepo.CancelPendingByUserAndCourse(user.ID, *campaign.CourseID); err != nil {
		log.Printf("[CampaignCheckout] Failed to cancel existing transactions: %v", err)
	}

	// Initialize payment provider
	if paymentProvider == nil {
		InitPaymentProvider()
		if paymentProvider == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Payment provider not configured",
			})
		}
	}

	// Generate order ID
	orderID := fmt.Sprintf("CAM-%s-%d", uuid.New().String()[:8], time.Now().UnixMilli()%100000)

	// Create transaction record
	tx := &postgres.Transaction{
		UserID:         user.ID,
		PaymentGateway: paymentProvider.GetName(),
		Amount:         finalPrice,
		Currency:       course.Currency,
		Status:         "pending",
	}
	tx.CourseID = campaign.CourseID
	tx.OrderID = &orderID
	
	totalDiscount := originalPrice - finalPrice
	if totalDiscount > 0 || appliedCoupon != nil {
		tx.OriginalAmount = &originalPrice
		tx.DiscountAmount = &totalDiscount
		if appliedCoupon != nil {
			tx.CouponID = &appliedCoupon.ID
		}
	}

	if err := paymentTxRepo.Create(tx); err != nil {
		log.Printf("[CampaignCheckout] Failed to create transaction: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create transaction"})
	}

	// Create payment via provider
	ctx := c.Request().Context()
	
	// Determine callback URL
	callbackURL := fmt.Sprintf("%s://%s/api/webhooks/duitku", c.Scheme(), c.Request().Host)
	returnURL := fmt.Sprintf("%s://%s/payment/success?order_id=%s", c.Scheme(), c.Request().Host, orderID)
	
	// Override with frontend URL if available
	frontendURL := getSettingValue("frontend_url", "")
	if frontendURL != "" {
		returnURL = fmt.Sprintf("%s/payment/success?order_id=%s", frontendURL, orderID)
	}

	paymentReq := &payment.CreateTransactionRequest{
		OrderID:       orderID,
		Amount:        finalPrice,
		Currency:      course.Currency,
		CustomerName:  user.FullName,
		CustomerEmail: user.Email,
		ItemName:      course.Title,
		ItemID:        course.ID,
		ItemCategory:  "course",
		PaymentMethod: req.PaymentMethod,
		ReturnURL:     returnURL,
		CallbackURL:   callbackURL,
	}

	log.Printf("[CampaignCheckout] Creating payment: OrderID=%s, User=%s, Amount=%.2f, Method=%s",
		orderID, user.Email, finalPrice, req.PaymentMethod)

	paymentResp, err := paymentProvider.CreateTransaction(ctx, paymentReq)
	if err != nil {
		log.Printf("[CampaignCheckout] Failed to create payment: %v", err)
		paymentTxRepo.UpdateStatus(tx.ID, "failure")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal membuat pembayaran: " + err.Error(),
		})
	}

	// Update transaction with payment info
	if paymentResp.ExpiredAt != nil {
		paymentTxRepo.UpdateSnapToken(tx.ID, paymentResp.SnapToken, paymentResp.PaymentURL, *paymentResp.ExpiredAt)
	}

	log.Printf("[CampaignCheckout] Payment created successfully: OrderID=%s, PaymentURL=%s",
		orderID, paymentResp.PaymentURL)

	return c.JSON(http.StatusOK, CampaignCheckoutResponse{
		TransactionID:  tx.ID,
		OrderID:        orderID,
		PaymentURL:     paymentResp.PaymentURL,
		ExpiredAt:      paymentResp.ExpiredAt,
		IsFree:         false,
		OriginalAmount: originalPrice,
		DiscountAmount: totalDiscount,
		FinalAmount:    finalPrice,
		IsNewUser:      isNewUser,
	})
}

// GetTransactionStatus returns transaction status by order ID (public endpoint)
// GET /api/transaction-status/:order_id
func GetTransactionStatus(c echo.Context) error {
	initPaymentRepos()

	orderID := c.Param("order_id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Order ID is required"})
	}

	tx, err := paymentTxRepo.GetByOrderID(orderID)
	if err != nil || tx == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Transaction not found"})
	}

	// Get course info
	var courseName string
	var courseSlug string
	if tx.CourseID != nil {
		course, _ := courseRepoCheckout.GetByID(*tx.CourseID)
		if course != nil {
			courseName = course.Title
			courseSlug = course.Slug
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"order_id":    orderID,
		"status":      tx.Status,
		"amount":      tx.Amount,
		"course_name": courseName,
		"course_slug": courseSlug,
		"created_at":  tx.CreatedAt,
	})
}
