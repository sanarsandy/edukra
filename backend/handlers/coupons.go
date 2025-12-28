package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

var couponRepo *postgres.CouponRepository

func initCouponRepo() {
	if couponRepo == nil {
		couponRepo = postgres.NewCouponRepository(db.DB)
	}
}

// ==================== ADMIN ENDPOINTS ====================

// ListCoupons returns all coupons (Admin only)
// GET /api/admin/coupons
func ListCoupons(c echo.Context) error {
	initCouponRepo()

	filters := make(map[string]interface{})
	
	if active := c.QueryParam("active"); active != "" {
		filters["is_active"] = active == "true"
	}

	coupons, err := couponRepo.List(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch coupons"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"coupons": coupons,
		"total":   len(coupons),
	})
}

// CreateCoupon creates a new coupon (Admin only)
// POST /api/admin/coupons
func CreateCoupon(c echo.Context) error {
	initCouponRepo()

	var req domain.CreateCouponRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Code == "" || req.DiscountType == "" || req.DiscountValue <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code, discount_type, and discount_value are required"})
	}

	// Validate discount type
	if req.DiscountType != "percentage" && req.DiscountType != "fixed" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "discount_type must be 'percentage' or 'fixed'"})
	}

	// Validate percentage range
	if req.DiscountType == "percentage" && req.DiscountValue > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Percentage discount cannot exceed 100%"})
	}

	// Parse dates
	validFrom := time.Now()
	if req.ValidFrom != nil && *req.ValidFrom != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidFrom); err == nil {
			validFrom = parsed
		}
	}

	var validUntil *time.Time
	if req.ValidUntil != nil && *req.ValidUntil != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidUntil); err == nil {
			validUntil = &parsed
		}
	}

	perUserLimit := 1
	if req.PerUserLimit != nil {
		perUserLimit = *req.PerUserLimit
	}

	coupon := &domain.Coupon{
		Code:          strings.ToUpper(strings.TrimSpace(req.Code)),
		DiscountType:  domain.DiscountType(req.DiscountType),
		DiscountValue: req.DiscountValue,
		MaxDiscount:   req.MaxDiscount,
		CourseID:      req.CourseID,
		UsageLimit:    req.UsageLimit,
		PerUserLimit:  perUserLimit,
		ValidFrom:     validFrom,
		ValidUntil:    validUntil,
		IsActive:      true,
	}

	if err := couponRepo.Create(coupon); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Coupon code already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create coupon"})
	}

	return c.JSON(http.StatusCreated, coupon)
}

// GetCoupon returns a single coupon by ID (Admin only)
// GET /api/admin/coupons/:id
func GetCoupon(c echo.Context) error {
	initCouponRepo()

	id := c.Param("id")
	coupon, err := couponRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch coupon"})
	}
	if coupon == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Coupon not found"})
	}

	return c.JSON(http.StatusOK, coupon)
}

// UpdateCoupon updates an existing coupon (Admin only)
// PUT /api/admin/coupons/:id
func UpdateCoupon(c echo.Context) error {
	initCouponRepo()

	id := c.Param("id")
	coupon, err := couponRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch coupon"})
	}
	if coupon == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Coupon not found"})
	}

	var req domain.UpdateCouponRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Apply updates
	if req.Code != nil {
		coupon.Code = strings.ToUpper(strings.TrimSpace(*req.Code))
	}
	if req.DiscountType != nil {
		coupon.DiscountType = domain.DiscountType(*req.DiscountType)
	}
	if req.DiscountValue != nil {
		coupon.DiscountValue = *req.DiscountValue
	}
	if req.MaxDiscount != nil {
		coupon.MaxDiscount = req.MaxDiscount
	}
	if req.CourseID != nil {
		coupon.CourseID = req.CourseID
	}
	if req.UsageLimit != nil {
		coupon.UsageLimit = req.UsageLimit
	}
	if req.PerUserLimit != nil {
		coupon.PerUserLimit = *req.PerUserLimit
	}
	if req.ValidFrom != nil && *req.ValidFrom != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidFrom); err == nil {
			coupon.ValidFrom = parsed
		}
	}
	if req.ValidUntil != nil && *req.ValidUntil != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidUntil); err == nil {
			coupon.ValidUntil = &parsed
		}
	}
	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}

	if err := couponRepo.Update(coupon); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update coupon"})
	}

	return c.JSON(http.StatusOK, coupon)
}

// DeleteCoupon deletes a coupon (Admin only)
// DELETE /api/admin/coupons/:id
func DeleteCoupon(c echo.Context) error {
	initCouponRepo()

	id := c.Param("id")
	if err := couponRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete coupon"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Coupon deleted successfully"})
}

// ==================== INSTRUCTOR ENDPOINTS ====================

// ListInstructorCoupons returns coupons created by the instructor
// GET /api/instructor/coupons
func ListInstructorCoupons(c echo.Context) error {
	initCouponRepo()

	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	filters := map[string]interface{}{
		"instructor_id": userID,
	}

	coupons, err := couponRepo.List(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch coupons"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"coupons": coupons,
		"total":   len(coupons),
	})
}

// CreateInstructorCoupon creates a coupon by instructor (for their courses)
// POST /api/instructor/coupons
func CreateInstructorCoupon(c echo.Context) error {
	initCouponRepo()

	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req domain.CreateCouponRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Code == "" || req.DiscountType == "" || req.DiscountValue <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code, discount_type, and discount_value are required"})
	}

	// Instructor must specify a course they own (optional based on business rules)
	// For now, we allow instructor to create coupons for their own courses or all their courses

	// Parse dates
	validFrom := time.Now()
	if req.ValidFrom != nil && *req.ValidFrom != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidFrom); err == nil {
			validFrom = parsed
		}
	}

	var validUntil *time.Time
	if req.ValidUntil != nil && *req.ValidUntil != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.ValidUntil); err == nil {
			validUntil = &parsed
		}
	}

	perUserLimit := 1
	if req.PerUserLimit != nil {
		perUserLimit = *req.PerUserLimit
	}

	coupon := &domain.Coupon{
		Code:          strings.ToUpper(strings.TrimSpace(req.Code)),
		DiscountType:  domain.DiscountType(req.DiscountType),
		DiscountValue: req.DiscountValue,
		MaxDiscount:   req.MaxDiscount,
		CourseID:      req.CourseID,
		InstructorID:  &userID,
		UsageLimit:    req.UsageLimit,
		PerUserLimit:  perUserLimit,
		ValidFrom:     validFrom,
		ValidUntil:    validUntil,
		IsActive:      true,
	}

	if err := couponRepo.Create(coupon); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Coupon code already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create coupon"})
	}

	return c.JSON(http.StatusCreated, coupon)
}

// ==================== PUBLIC ENDPOINTS (Checkout) ====================

// ValidateCoupon validates a coupon code at checkout
// POST /api/coupons/validate
func ValidateCoupon(c echo.Context) error {
	initCouponRepo()
	initPaymentRepos() // for course repo

	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req domain.ValidateCouponRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Code == "" || req.CourseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code and course_id are required"})
	}

	// Get course to calculate discount
	course, err := courseRepoCheckout.GetByID(req.CourseID)
	if err != nil || course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}

	// Get coupon by code
	coupon, err := couponRepo.GetByCode(req.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate coupon"})
	}
	if coupon == nil {
		return c.JSON(http.StatusOK, domain.ValidateCouponResponse{
			Valid:   false,
			Message: "Kode kupon tidak ditemukan",
		})
	}

	// Validate coupon for this user and course
	valid, message := couponRepo.ValidateCouponForUser(coupon.ID, userID, req.CourseID)
	if !valid {
		return c.JSON(http.StatusOK, domain.ValidateCouponResponse{
			Valid:   false,
			Message: message,
		})
	}

	// Determine effective price (use amount from request if provided, otherwise calculate from course)
	effectivePrice := course.Price
	if req.Amount > 0 {
		// Frontend sent the effective price (after course discount)
		effectivePrice = req.Amount
	} else if course.DiscountPrice != nil && *course.DiscountPrice > 0 {
		// Check if course has an active discount
		discountValid := true
		if course.DiscountValidUntil != nil {
			discountValid = course.DiscountValidUntil.After(time.Now())
		}
		if discountValid {
			effectivePrice = *course.DiscountPrice
		}
	}

	// Calculate discount based on EFFECTIVE price (not original price)
	discountAmount := coupon.CalculateDiscount(effectivePrice)
	finalPrice := effectivePrice - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	return c.JSON(http.StatusOK, domain.ValidateCouponResponse{
		Valid:          true,
		Coupon:         coupon, // Include coupon data for frontend
		DiscountType:   string(coupon.DiscountType),
		DiscountValue:  coupon.DiscountValue,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		OriginalPrice:  course.Price,
		CouponID:       coupon.ID,
	})
}
