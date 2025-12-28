package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// CouponRepository handles coupon database operations
type CouponRepository struct {
	db *sqlx.DB
}

// NewCouponRepository creates a new coupon repository
func NewCouponRepository(db *sqlx.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

// Create creates a new coupon
func (r *CouponRepository) Create(coupon *domain.Coupon) error {
	// Normalize code
	coupon.NormalizeCode()
	
	query := `
		INSERT INTO coupons (
			tenant_id, code, discount_type, discount_value, max_discount,
			course_id, instructor_id, usage_limit, per_user_limit,
			valid_from, valid_until, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	
	return r.db.QueryRow(
		query,
		coupon.TenantID,
		coupon.Code,
		coupon.DiscountType,
		coupon.DiscountValue,
		coupon.MaxDiscount,
		coupon.CourseID,
		coupon.InstructorID,
		coupon.UsageLimit,
		coupon.PerUserLimit,
		coupon.ValidFrom,
		coupon.ValidUntil,
		coupon.IsActive,
	).Scan(&coupon.ID, &coupon.CreatedAt, &coupon.UpdatedAt)
}

// GetByID retrieves a coupon by ID
func (r *CouponRepository) GetByID(id string) (*domain.Coupon, error) {
	query := `
		SELECT id, tenant_id, code, discount_type, discount_value, max_discount,
		       course_id, instructor_id, usage_limit, usage_count, per_user_limit,
		       valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons WHERE id = $1
	`
	
	coupon := &domain.Coupon{}
	err := r.db.QueryRow(query, id).Scan(
		&coupon.ID, &coupon.TenantID, &coupon.Code, &coupon.DiscountType, &coupon.DiscountValue,
		&coupon.MaxDiscount, &coupon.CourseID, &coupon.InstructorID, &coupon.UsageLimit,
		&coupon.UsageCount, &coupon.PerUserLimit, &coupon.ValidFrom, &coupon.ValidUntil,
		&coupon.IsActive, &coupon.CreatedAt, &coupon.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return coupon, nil
}

// GetByCode retrieves a coupon by code (case-insensitive)
func (r *CouponRepository) GetByCode(code string) (*domain.Coupon, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	
	query := `
		SELECT id, tenant_id, code, discount_type, discount_value, max_discount,
		       course_id, instructor_id, usage_limit, usage_count, per_user_limit,
		       valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons WHERE UPPER(code) = $1
	`
	
	coupon := &domain.Coupon{}
	err := r.db.QueryRow(query, code).Scan(
		&coupon.ID, &coupon.TenantID, &coupon.Code, &coupon.DiscountType, &coupon.DiscountValue,
		&coupon.MaxDiscount, &coupon.CourseID, &coupon.InstructorID, &coupon.UsageLimit,
		&coupon.UsageCount, &coupon.PerUserLimit, &coupon.ValidFrom, &coupon.ValidUntil,
		&coupon.IsActive, &coupon.CreatedAt, &coupon.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return coupon, nil
}

// List retrieves all coupons with optional filters
func (r *CouponRepository) List(filters map[string]interface{}) ([]*domain.Coupon, error) {
	query := `
		SELECT id, tenant_id, code, discount_type, discount_value, max_discount,
		       course_id, instructor_id, usage_limit, usage_count, per_user_limit,
		       valid_from, valid_until, is_active, created_at, updated_at
		FROM coupons WHERE 1=1
	`
	
	args := []interface{}{}
	argCount := 1
	
	// Apply filters
	if v, ok := filters["is_active"]; ok {
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, v)
		argCount++
	}
	
	if v, ok := filters["instructor_id"]; ok {
		query += fmt.Sprintf(" AND instructor_id = $%d", argCount)
		args = append(args, v)
		argCount++
	}
	
	if v, ok := filters["course_id"]; ok {
		query += fmt.Sprintf(" AND (course_id = $%d OR course_id IS NULL)", argCount)
		args = append(args, v)
		argCount++
	}
	
	query += " ORDER BY created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var coupons []*domain.Coupon
	for rows.Next() {
		coupon := &domain.Coupon{}
		err := rows.Scan(
			&coupon.ID, &coupon.TenantID, &coupon.Code, &coupon.DiscountType, &coupon.DiscountValue,
			&coupon.MaxDiscount, &coupon.CourseID, &coupon.InstructorID, &coupon.UsageLimit,
			&coupon.UsageCount, &coupon.PerUserLimit, &coupon.ValidFrom, &coupon.ValidUntil,
			&coupon.IsActive, &coupon.CreatedAt, &coupon.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	
	return coupons, nil
}

// Update updates a coupon
func (r *CouponRepository) Update(coupon *domain.Coupon) error {
	query := `
		UPDATE coupons SET
			code = $1, discount_type = $2, discount_value = $3, max_discount = $4,
			course_id = $5, usage_limit = $6, per_user_limit = $7,
			valid_from = $8, valid_until = $9, is_active = $10, updated_at = $11
		WHERE id = $12
	`
	
	_, err := r.db.Exec(
		query,
		coupon.Code, coupon.DiscountType, coupon.DiscountValue, coupon.MaxDiscount,
		coupon.CourseID, coupon.UsageLimit, coupon.PerUserLimit,
		coupon.ValidFrom, coupon.ValidUntil, coupon.IsActive, time.Now(), coupon.ID,
	)
	
	return err
}

// Delete deletes a coupon
func (r *CouponRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM coupons WHERE id = $1", id)
	return err
}

// IncrementUsage increments the usage count for a coupon
func (r *CouponRepository) IncrementUsage(id string) error {
	_, err := r.db.Exec("UPDATE coupons SET usage_count = usage_count + 1, updated_at = $1 WHERE id = $2", time.Now(), id)
	return err
}

// GetUserUsageCount returns how many times a user has used a specific coupon
func (r *CouponRepository) GetUserUsageCount(couponID, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM coupon_usages WHERE coupon_id = $1 AND user_id = $2",
		couponID, userID,
	).Scan(&count)
	return count, err
}

// RecordUsage records a coupon usage
func (r *CouponRepository) RecordUsage(usage *domain.CouponUsage) error {
	query := `
		INSERT INTO coupon_usages (coupon_id, user_id, transaction_id, discount_applied)
		VALUES ($1, $2, $3, $4)
		RETURNING id, used_at
	`
	
	return r.db.QueryRow(
		query,
		usage.CouponID, usage.UserID, usage.TransactionID, usage.DiscountApplied,
	).Scan(&usage.ID, &usage.UsedAt)
}

// ValidateCouponForUser checks if a coupon is valid for a specific user and course
func (r *CouponRepository) ValidateCouponForUser(couponID, userID, courseID string) (bool, string) {
	coupon, err := r.GetByID(couponID)
	if err != nil || coupon == nil {
		return false, "Kupon tidak ditemukan"
	}
	
	// Check if coupon is active and within date range
	if !coupon.IsValid() {
		if !coupon.IsActive {
			return false, "Kupon tidak aktif"
		}
		if coupon.ValidUntil != nil && time.Now().After(*coupon.ValidUntil) {
			return false, "Kupon sudah kadaluarsa"
		}
		if coupon.UsageLimit != nil && coupon.UsageCount >= *coupon.UsageLimit {
			return false, "Kupon sudah mencapai batas penggunaan"
		}
		return false, "Kupon tidak valid"
	}
	
	// Check course scope
	if coupon.CourseID != nil && *coupon.CourseID != courseID {
		return false, "Kupon tidak berlaku untuk kursus ini"
	}
	
	// Check per-user limit
	usageCount, err := r.GetUserUsageCount(couponID, userID)
	if err != nil {
		return false, "Gagal memeriksa penggunaan kupon"
	}
	if usageCount >= coupon.PerUserLimit {
		return false, "Anda sudah mencapai batas penggunaan kupon ini"
	}
	
	return true, ""
}
