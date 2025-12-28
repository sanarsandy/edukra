package domain

import (
	"strings"
	"time"
)

// DiscountType represents the type of discount
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID            string        `json:"id"`
	TenantID      *string       `json:"tenant_id,omitempty"`
	Code          string        `json:"code"`
	DiscountType  DiscountType  `json:"discount_type"`
	DiscountValue float64       `json:"discount_value"`
	MaxDiscount   *float64      `json:"max_discount,omitempty"`
	CourseID      *string       `json:"course_id,omitempty"`
	InstructorID  *string       `json:"instructor_id,omitempty"`
	UsageLimit    *int          `json:"usage_limit,omitempty"`
	UsageCount    int           `json:"usage_count"`
	PerUserLimit  int           `json:"per_user_limit"`
	ValidFrom     time.Time     `json:"valid_from"`
	ValidUntil    *time.Time    `json:"valid_until,omitempty"`
	IsActive      bool          `json:"is_active"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// CouponUsage tracks when a user uses a coupon
type CouponUsage struct {
	ID              string    `json:"id"`
	CouponID        string    `json:"coupon_id"`
	UserID          string    `json:"user_id"`
	TransactionID   *string   `json:"transaction_id,omitempty"`
	DiscountApplied float64   `json:"discount_applied"`
	UsedAt          time.Time `json:"used_at"`
}

// CreateCouponRequest represents the request to create a coupon
type CreateCouponRequest struct {
	Code          string   `json:"code" validate:"required,min=3,max=50"`
	DiscountType  string   `json:"discount_type" validate:"required,oneof=percentage fixed"`
	DiscountValue float64  `json:"discount_value" validate:"required,gt=0"`
	MaxDiscount   *float64 `json:"max_discount,omitempty"`
	CourseID      *string  `json:"course_id,omitempty"`
	UsageLimit    *int     `json:"usage_limit,omitempty"`
	PerUserLimit  *int     `json:"per_user_limit,omitempty"`
	ValidFrom     *string  `json:"valid_from,omitempty"`
	ValidUntil    *string  `json:"valid_until,omitempty"`
}

// UpdateCouponRequest represents the request to update a coupon
type UpdateCouponRequest struct {
	Code          *string  `json:"code,omitempty"`
	DiscountType  *string  `json:"discount_type,omitempty"`
	DiscountValue *float64 `json:"discount_value,omitempty"`
	MaxDiscount   *float64 `json:"max_discount,omitempty"`
	CourseID      *string  `json:"course_id,omitempty"`
	UsageLimit    *int     `json:"usage_limit,omitempty"`
	PerUserLimit  *int     `json:"per_user_limit,omitempty"`
	ValidFrom     *string  `json:"valid_from,omitempty"`
	ValidUntil    *string  `json:"valid_until,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// ValidateCouponRequest is for validating a coupon at checkout
type ValidateCouponRequest struct {
	Code     string  `json:"code" validate:"required"`
	CourseID string  `json:"course_id" validate:"required"`
	Amount   float64 `json:"amount,omitempty"` // Effective price after course discount
}

// ValidateCouponResponse returns discount info
type ValidateCouponResponse struct {
	Valid           bool    `json:"valid"`
	Message         string  `json:"message,omitempty"`
	Coupon          *Coupon `json:"coupon,omitempty"` // Include full coupon data
	DiscountType    string  `json:"discount_type,omitempty"`
	DiscountValue   float64 `json:"discount_value,omitempty"`
	DiscountAmount  float64 `json:"discount_amount,omitempty"`
	FinalPrice      float64 `json:"final_price,omitempty"`
	OriginalPrice   float64 `json:"original_price,omitempty"`
	CouponID        string  `json:"coupon_id,omitempty"`
}

// NormalizeCode uppercases and trims the coupon code
func (c *Coupon) NormalizeCode() {
	c.Code = strings.ToUpper(strings.TrimSpace(c.Code))
}

// IsValid checks if the coupon is currently valid
func (c *Coupon) IsValid() bool {
	if !c.IsActive {
		return false
	}
	
	now := time.Now()
	
	// Check valid_from
	if now.Before(c.ValidFrom) {
		return false
	}
	
	// Check valid_until
	if c.ValidUntil != nil && now.After(*c.ValidUntil) {
		return false
	}
	
	// Check usage limit
	if c.UsageLimit != nil && c.UsageCount >= *c.UsageLimit {
		return false
	}
	
	return true
}

// CalculateDiscount calculates the discount amount for a given price
func (c *Coupon) CalculateDiscount(originalPrice float64) float64 {
	var discount float64
	
	switch c.DiscountType {
	case DiscountTypePercentage:
		discount = originalPrice * (c.DiscountValue / 100)
		// Apply max discount cap if set
		if c.MaxDiscount != nil && discount > *c.MaxDiscount {
			discount = *c.MaxDiscount
		}
	case DiscountTypeFixed:
		discount = c.DiscountValue
	}
	
	// Discount cannot exceed original price
	if discount > originalPrice {
		discount = originalPrice
	}
	
	return discount
}

// CalculateFinalPrice calculates the final price after discount
func (c *Coupon) CalculateFinalPrice(originalPrice float64) float64 {
	discount := c.CalculateDiscount(originalPrice)
	return originalPrice - discount
}
