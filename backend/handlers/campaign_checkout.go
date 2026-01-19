package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/payment"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var (
	campaignRepoCheckout *postgres.CampaignRepository
)

func initCampaignCheckoutRepos() {
	if campaignRepoCheckout == nil && db.DB != nil {
	campaignRepoCheckout = postgres.NewCampaignRepository(db.DB)
	}
	if webinarRepo == nil && db.DB != nil {
		webinarRepo = postgres.NewWebinarRepository(db.DB)
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
	
	// Try to find existing user by email (case-insensitive)
	// Use empty tenant for now (NULL tenant_id)
	existingUser, err := userRepo.GetByEmail("", email)
	if err != nil {
		// DB error occurred - don't create new user, return error
		log.Printf("[CampaignCheckout] DB error looking up user by email %s: %v", email, err)
		return nil, false, fmt.Errorf("failed to lookup user: %v", err)
	}
	
	if existingUser != nil {
		log.Printf("[CampaignCheckout] Found existing user for email %s: %s", email, existingUser.ID)
		
		needsUpdate := false
		
		// Update FullName if provided from checkout form
		if fullName != "" && existingUser.FullName != fullName {
			log.Printf("[CampaignCheckout] Updating fullName for user %s: old=%s, new=%s", existingUser.ID, existingUser.FullName, fullName)
			existingUser.FullName = fullName
			needsUpdate = true
		}
		
		// Update phone if not set
		// SECURITY: Do NOT overwrite existing phone number from guest checkout
		if phone != "" {
			normalizedPhone := normalizePhone(phone)
			
			if existingUser.Phone == nil || *existingUser.Phone == "" {
				log.Printf("[CampaignCheckout] Updating empty phone for user %s: new=%s", existingUser.ID, normalizedPhone)
				phoneVal := normalizedPhone
				existingUser.Phone = &phoneVal
				needsUpdate = true
			} else if *existingUser.Phone != normalizedPhone {
				log.Printf("[CampaignCheckout] SECURITY: Skipping phone update for user %s. Existing: %s, Input: %s", existingUser.ID, *existingUser.Phone, normalizedPhone)
			}
		}
		
		// Persist updates if any
		if needsUpdate {
			if err := userRepo.Update(existingUser); err != nil {
				log.Printf("[CampaignCheckout] Failed to update user: %v", err)
				// Continue anyway, don't block checkout
			} else {
				log.Printf("[CampaignCheckout] User updated successfully")
			}
		}
		
		return existingUser, false, nil
	}
	
	// No existing user found - create new guest user
	log.Printf("[CampaignCheckout] No existing user found for email %s, creating new guest user", email)
	
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
		// Check if it's a unique constraint violation (email already exists)
		// This can happen in race conditions
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			log.Printf("[CampaignCheckout] Race condition: user %s was created by another request, fetching...", email)
			// Try to fetch the existing user again
			existingUser, fetchErr := userRepo.GetByEmail("", email)
			if fetchErr == nil && existingUser != nil {
				return existingUser, false, nil
			}
		}
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

	// Validate required fields (payment_method validated later based on isFree)
	if req.CampaignID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign ID is required"})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}
	if req.Phone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Phone number is required"})
	}
	// Note: PaymentMethod validation moved after isFree determination

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
	if (campaign.CourseID == nil || *campaign.CourseID == "") && campaign.CampaignType != "webinar_only" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign has no linked course"})
	}

	// Get course
	var course *domain.Course
	if campaign.CourseID != nil && *campaign.CourseID != "" {
		fetchedCourse, err := courseRepoCheckout.GetByID(*campaign.CourseID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
		}
		course = fetchedCourse
	}
	
	if course == nil && campaign.CampaignType != "webinar_only" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}

	// Find or create user
	user, isNewUser, err := findOrCreateGuestUser(req.Email, req.Phone, req.FullName)
	if err != nil {
		log.Printf("[CampaignCheckout] Failed to find/create user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process user"})
	}

	// Check if already enrolled
	if campaign.CourseID != nil {
		enrolled, err := enrollmentRepoCheckout.IsEnrolled(user.ID, *campaign.CourseID)
		if err == nil && enrolled {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Anda sudah terdaftar di kursus ini",
			})
		}
	}

	// Determine if this is a free registration
	// Priority: 1) campaign.IsFreeWebinar override, 2) course.Price == 0
	isFreeWebinar := false
	if campaign.IsFreeWebinar != nil {
		isFreeWebinar = *campaign.IsFreeWebinar
	} else if course != nil {
		isFreeWebinar = course.Price == 0
	} else if campaign.CampaignType == "webinar_only" {
		// Webinar only without price override is assumed free for now (until we have webinar pricing)
		isFreeWebinar = true
	}

	// Validate payment method only for paid courses
	if !isFreeWebinar && req.PaymentMethod == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payment method is required for paid courses"})
	}

	// Handle free webinar registration (no payment needed)
	if isFreeWebinar {
		campaignType := campaign.CampaignType
		if campaignType == "" {
			campaignType = domain.CampaignTypeEcourseOnly // default for backward compatibility
		}

		cidSafe := "nil"
		if campaign.CourseID != nil {
			cidSafe = *campaign.CourseID
		}
		log.Printf("[CampaignCheckout] Processing free registration: type=%s, user=%s, course=%s", 
			campaignType, user.Email, cidSafe)

		switch campaignType {
		case domain.CampaignTypeWebinarOnly:
			// Webinar only - register to webinar without course enrollment
			// Use direct webinar_id if available, otherwise fall back to course-based webinar lookup
			go handleWebinarOnlyNotificationDirect(user.ID, campaign.WebinarID, campaign.CourseID)
			
			return c.JSON(http.StatusOK, CampaignCheckoutResponse{
				IsFree:    true,
				Message:   "Selamat! Anda berhasil terdaftar webinar. Informasi akan dikirim via WhatsApp.",
				IsNewUser: isNewUser,
			})

		case domain.CampaignTypeEcourseOnly:
			// Course only - enroll without webinar registration
			if campaign.CourseID == nil {
				log.Printf("[CampaignCheckout] Error: ecourse_only campaign %s has no linked course", campaign.ID)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Campaign configuration error: No linked course"})
			}
			enrollment := &postgres.Enrollment{
				UserID:   user.ID,
				CourseID: *campaign.CourseID,
			}
			if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
				log.Printf("[CampaignCheckout] Failed to enroll (ecourse_only): %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
			}

			go handleEcourseOnlyNotification(user.ID, *campaign.CourseID)

			return c.JSON(http.StatusOK, CampaignCheckoutResponse{
				IsFree:    true,
				Message:   "Selamat! Anda berhasil terdaftar. Akses kursus Anda sudah tersedia.",
				IsNewUser: isNewUser,
			})

		case domain.CampaignTypeWebinarEcourse:
			// Both webinar and course - enroll AND register to webinar
			if campaign.CourseID == nil {
				log.Printf("[CampaignCheckout] Error: webinar_ecourse campaign %s has no linked course", campaign.ID)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Campaign configuration error: No linked course"})
			}
			enrollment := &postgres.Enrollment{
				UserID:   user.ID,
				CourseID: *campaign.CourseID,
			}
			if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
				log.Printf("[CampaignCheckout] Failed to enroll (webinar_ecourse): %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
			}

			// This will handle both webinar registration and combined notification
			go handlePaymentSuccessNotification(user.ID, *campaign.CourseID)

			return c.JSON(http.StatusOK, CampaignCheckoutResponse{
				IsFree:    true,
				Message:   "Selamat! Anda berhasil terdaftar. Informasi webinar dan akses kursus akan dikirim via WhatsApp.",
				IsNewUser: isNewUser,
			})

		default:
			// Fallback to existing behavior (enroll + webinar notification)
			if campaign.CourseID == nil {
				log.Printf("[CampaignCheckout] Error: default campaign %s has no linked course", campaign.ID)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Campaign configuration error: No linked course"})
			}
			enrollment := &postgres.Enrollment{
				UserID:   user.ID,
				CourseID: *campaign.CourseID,
			}
			if err := enrollmentRepoCheckout.Create(enrollment); err != nil {
				log.Printf("[CampaignCheckout] Failed to enroll (default): %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to enroll"})
			}

			go handlePaymentSuccessNotification(user.ID, *campaign.CourseID)

			return c.JSON(http.StatusOK, CampaignCheckoutResponse{
				IsFree:    true,
				Message:   "Selamat! Anda berhasil terdaftar. Informasi akan dikirim via WhatsApp.",
				IsNewUser: isNewUser,
			})
		}
	}

	// Calculate price
	// Ensure course is not nil for paid flows
	if course == nil {
		log.Printf("[CampaignCheckout] Error: Paid campaign type %s has no linked course for pricing", campaign.CampaignType)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Konfigurasi campaign tidak valid: Webinar berbayar harus terhubung dengan kursus (untuk harga).",
		})
	}
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

		// Register to webinar if exists (for free courses)
		go handlePaymentSuccessNotification(user.ID, *campaign.CourseID)

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

		// Send payment success notification and handle webinar registration
		if user != nil {
			go handlePaymentSuccessNotification(user.ID, *campaign.CourseID)
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
	
	// Store metadata (webinar_id)
	metadata := make(map[string]string)
	if campaign.WebinarID != nil {
		metadata["webinar_id"] = *campaign.WebinarID
	}
	if len(metadata) > 0 {
		metaJSON, _ := json.Marshal(metadata)
		tx.Metadata = metaJSON
	}
	
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
	// Determine callback URL
	callbackURL := fmt.Sprintf("%s://%s/api/webhooks/duitku", c.Scheme(), c.Request().Host)
	
	// Default return URL (backend) - not ideal
	returnURL := fmt.Sprintf("%s://%s/payment/success?order_id=%s", c.Scheme(), c.Request().Host, orderID)
	
	// Try getting from Env first (Best for Docker/Prod)
	if envFrontend := os.Getenv("FRONTEND_URL"); envFrontend != "" {
		returnURL = fmt.Sprintf("%s/payment/success?order_id=%s", envFrontend, orderID)
	} else {
		// Try getting from DB Settings
		frontendURL := getSettingValue("frontend_url", "")
		if frontendURL != "" {
			returnURL = fmt.Sprintf("%s/payment/success?order_id=%s", frontendURL, orderID)
		} else if strings.Contains(c.Request().Host, "localhost") {
			// Fallback for local development if settings not set
			returnURL = fmt.Sprintf("http://localhost:3000/payment/success?order_id=%s", orderID)
		}
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

// handleWebinarOnlyNotification handles webinar-only registration (no course enrollment)
// Registers user to webinar and sends webinar-only WhatsApp notification
func handleWebinarOnlyNotification(userID, courseID string) {
	webinarRepo := postgres.NewWebinarRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)

	// Get user info
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		log.Printf("[WebinarOnly] Failed to get user %s: %v", userID, err)
		return
	}

	// Get LMS URL
	lmsURL := getSettingValue("frontend_url", "")
	if lmsURL == "" {
		lmsURL = os.Getenv("FRONTEND_URL")
	}
	if lmsURL == "" {
		lmsURL = "https://lms.edukra.com"
	}

	// Check for upcoming webinars
	webinars, err := webinarRepo.GetUpcomingByCourse(courseID)
	if err != nil {
		log.Printf("[WebinarOnly] Failed to check webinars for course %s: %v", courseID, err)
		return
	}

	if len(webinars) == 0 {
		log.Printf("[WebinarOnly] No upcoming webinars found for course %s", courseID)
		return
	}

	// Register user to each webinar
	for _, webinar := range webinars {
		if err := webinarRepo.RegisterUser(webinar.ID, userID, "campaign"); err != nil {
			log.Printf("[WebinarOnly] Failed to register user %s to webinar %s: %v", userID, webinar.ID, err)
			continue
		}
		log.Printf("[WebinarOnly] User %s registered to webinar %s", userID, webinar.ID)
		webinarRepo.CreateReminders(webinar.ID, userID, webinar.ScheduledAt)
	}

	// Send Webinar Confirmation WA
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
		log.Printf("[WebinarOnly] Webinar confirmation sent to %s", *user.Phone)
	}
}

// handleEcourseOnlyNotification handles course-only registration (no webinar)
// Sends course access notification without webinar info
func handleEcourseOnlyNotification(userID, courseID string) {
	userRepo := postgres.NewUserRepository(db.DB)
	courseRepo := postgres.NewCourseRepository(db.DB)

	// Get user info
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		log.Printf("[EcourseOnly] Failed to get user %s: %v", userID, err)
		return
	}

	// Get LMS URL
	lmsURL := getSettingValue("frontend_url", "")
	if lmsURL == "" {
		lmsURL = os.Getenv("FRONTEND_URL")
	}
	if lmsURL == "" {
		lmsURL = "https://lms.edukra.com"
	}

	// Get course info
	course, err := courseRepo.GetByID(courseID)
	if err != nil || course == nil {
		log.Printf("[EcourseOnly] Failed to get course %s: %v", courseID, err)
		return
	}

	// Send Course Access WA
	if user.Phone != nil && *user.Phone != "" {
		if err := service.GetWhatsAppService().SendPaymentSuccess(*user.Phone, user.FullName, course.Title, lmsURL); err != nil {
			log.Printf("[EcourseOnly] Failed to send course access WA to %s: %v", *user.Phone, err)
		} else {
			log.Printf("[EcourseOnly] Course access WA sent to %s", *user.Phone)
		}
	}
}

// handleWebinarOnlyNotificationDirect handles webinar-only registration with direct webinar_id
// Supports both direct webinar link and course-based webinar lookup (fallback)
func handleWebinarOnlyNotificationDirect(userID string, webinarID *string, courseID *string) {
	webinarRepo := postgres.NewWebinarRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)

	// Get user info
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		log.Printf("[WebinarOnlyDirect] Failed to get user %s: %v", userID, err)
		return
	}

	// Get LMS URL
	lmsURL := getSettingValue("frontend_url", "")
	if lmsURL == "" {
		lmsURL = os.Getenv("FRONTEND_URL")
	}
	if lmsURL == "" {
		lmsURL = "https://lms.edukra.com"
	}

	var webinar *domain.Webinar

	// If direct webinar_id is provided, use it
	if webinarID != nil && *webinarID != "" {
		webinar, err = webinarRepo.GetByID(*webinarID)
		if err != nil || webinar == nil {
			log.Printf("[WebinarOnlyDirect] Failed to get webinar %s: %v", *webinarID, err)
			return
		}

		// Register user to this specific webinar
		if err := webinarRepo.RegisterUser(webinar.ID, userID, "campaign"); err != nil {
			log.Printf("[WebinarOnlyDirect] Failed to register user %s to webinar %s: %v", userID, webinar.ID, err)
		} else {
			log.Printf("[WebinarOnlyDirect] User %s registered to webinar %s", userID, webinar.ID)
			webinarRepo.CreateReminders(webinar.ID, userID, webinar.ScheduledAt)
		}
	} else if courseID != nil && *courseID != "" {
		// Fallback: lookup webinars from course
		webinars, err := webinarRepo.GetUpcomingByCourse(*courseID)
		if err != nil || len(webinars) == 0 {
			log.Printf("[WebinarOnlyDirect] No webinars found for course %s", *courseID)
			return
		}

		// Register to all upcoming webinars for the course
		for _, w := range webinars {
			if err := webinarRepo.RegisterUser(w.ID, userID, "campaign"); err != nil {
				log.Printf("[WebinarOnlyDirect] Failed to register user %s to webinar %s: %v", userID, w.ID, err)
				continue
			}
			log.Printf("[WebinarOnlyDirect] User %s registered to webinar %s", userID, w.ID)
			webinarRepo.CreateReminders(w.ID, userID, w.ScheduledAt)
		}
		webinar = webinars[0]
	} else {
		log.Printf("[WebinarOnlyDirect] No webinar_id or course_id provided")
		return
	}

	// Send Webinar Confirmation WA
	if user.Phone != nil && *user.Phone != "" && webinar != nil {
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

		service.SendWebinarOnlyConfirmationAsync(*user.Phone, data)
		log.Printf("[WebinarOnlyDirect] Webinar confirmation sent to %s", *user.Phone)
	}
}
