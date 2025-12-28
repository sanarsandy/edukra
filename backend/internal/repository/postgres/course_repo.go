package postgres

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// CourseRepository implements domain.CourseRepository using PostgreSQL
type CourseRepository struct {
	db *sqlx.DB
}

// NewCourseRepository creates a new CourseRepository
func NewCourseRepository(db *sqlx.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(id string) (*domain.Course, error) {
	query := `
		SELECT id, tenant_id, instructor_id, category_id, title, slug, description, thumbnail_url,
		       price, discount_price, discount_valid_until, currency, COALESCE(lessons_count, 0), COALESCE(duration, ''), is_published, is_featured, created_at, updated_at
		FROM courses WHERE id = $1
	`
	
	var course domain.Course
	var tenantID, instructorID, categoryID, thumbnailURL sql.NullString
	var discountPrice sql.NullFloat64
	var discountValidUntil sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&course.ID, &tenantID, &instructorID, &categoryID, &course.Title, &course.Slug,
		&course.Description, &thumbnailURL, &course.Price, &discountPrice, &discountValidUntil,
		&course.Currency, &course.LessonsCount, &course.Duration, &course.IsPublished, &course.IsFeatured, &course.CreatedAt, &course.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tenantID.Valid {
		course.TenantID = tenantID.String
	}
	if instructorID.Valid {
		course.InstructorID = &instructorID.String
	}
	if categoryID.Valid {
		course.CategoryID = &categoryID.String
	}
	if thumbnailURL.Valid {
		course.ThumbnailURL = &thumbnailURL.String
	}
	if discountPrice.Valid {
		course.DiscountPrice = &discountPrice.Float64
	}
	if discountValidUntil.Valid {
		course.DiscountValidUntil = &discountValidUntil.Time
	}
	
	return &course, nil
}

// GetBySlug retrieves a course by tenant and slug
func (r *CourseRepository) GetBySlug(tenantID, slug string) (*domain.Course, error) {
	query := `
		SELECT id, tenant_id, instructor_id, category_id, title, slug, description, thumbnail_url,
		       price, currency, COALESCE(lessons_count, 0), COALESCE(duration, ''), is_published, is_featured, created_at, updated_at
		FROM courses WHERE tenant_id = $1 AND slug = $2
	`
	
	var course domain.Course
	var tenantIDVal, instructorID, categoryID, thumbnailURL sql.NullString
	
	err := r.db.QueryRow(query, tenantID, slug).Scan(
		&course.ID, &tenantIDVal, &instructorID, &categoryID, &course.Title, &course.Slug,
		&course.Description, &thumbnailURL, &course.Price, &course.Currency,
		&course.LessonsCount, &course.Duration, &course.IsPublished, &course.IsFeatured, &course.CreatedAt, &course.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tenantIDVal.Valid {
		course.TenantID = tenantIDVal.String
	}
	if instructorID.Valid {
		course.InstructorID = &instructorID.String
	}
	if categoryID.Valid {
		course.CategoryID = &categoryID.String
	}
	if thumbnailURL.Valid {
		course.ThumbnailURL = &thumbnailURL.String
	}
	
	return &course, nil
}

// Create inserts a new course
func (r *CourseRepository) Create(course *domain.Course) error {
	// Generate slug from title if not provided
	if course.Slug == "" {
		course.Slug = generateSlug(course.Title)
	}
	
	query := `
		INSERT INTO courses (tenant_id, instructor_id, category_id, title, slug, description, 
		                     thumbnail_url, price, currency, lessons_count, duration, is_published, is_featured,
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`
	
	now := time.Now()
	course.CreatedAt = now
	course.UpdatedAt = now
	
	if course.Currency == "" {
		course.Currency = "IDR"
	}
	
	var tenantID interface{} = course.TenantID
	if course.TenantID == "" || course.TenantID == "default" {
		tenantID = nil
	}

	return r.db.QueryRow(query,
		tenantID, course.InstructorID, course.CategoryID, course.Title, course.Slug,
		course.Description, course.ThumbnailURL, course.Price, course.Currency,
		course.LessonsCount, course.Duration, course.IsPublished, course.IsFeatured, course.CreatedAt, course.UpdatedAt,
	).Scan(&course.ID)
}

// Update updates an existing course
func (r *CourseRepository) Update(course *domain.Course) error {
	query := `
		UPDATE courses SET
			title = $2, slug = $3, description = $4, thumbnail_url = $5,
			price = $6, discount_price = $7, discount_valid_until = $8,
			lessons_count = $9, duration = $10, is_published = $11, is_featured = $12, category_id = $13, updated_at = $14
		WHERE id = $1
	`
	
	course.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		course.ID, course.Title, course.Slug, course.Description,
		course.ThumbnailURL, course.Price, course.DiscountPrice, course.DiscountValidUntil,
		course.LessonsCount, course.Duration, course.IsPublished,
		course.IsFeatured, course.CategoryID, course.UpdatedAt,
	)
	return err
}

// Delete removes a course by ID
func (r *CourseRepository) Delete(id string) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListByTenant retrieves courses for a tenant with pagination
func (r *CourseRepository) ListByTenant(tenantID string, limit, offset int) ([]*domain.Course, error) {
	var query string
	var args []interface{}

	if tenantID == "" || tenantID == "default" {
		query = `
			SELECT c.id, c.tenant_id, c.instructor_id, c.category_id, c.title, c.slug, c.description, c.thumbnail_url,
			       c.price, c.discount_price, c.discount_valid_until, c.currency, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''), c.is_published, c.is_featured, c.created_at, c.updated_at,
				   u.full_name, cat.name
			FROM courses c
			LEFT JOIN users u ON c.instructor_id = u.id
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.tenant_id IS NULL
			ORDER BY c.created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	} else {
		query = `
			SELECT c.id, c.tenant_id, c.instructor_id, c.category_id, c.title, c.slug, c.description, c.thumbnail_url,
			       c.price, c.discount_price, c.discount_valid_until, c.currency, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''), c.is_published, c.is_featured, c.created_at, c.updated_at,
				   u.full_name, cat.name
			FROM courses c
			LEFT JOIN users u ON c.instructor_id = u.id
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.tenant_id = $1
			ORDER BY c.created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{tenantID, limit, offset}
	}
	
	return r.scanCoursesWithDetails(query, args...)
}

// ListByInstructor retrieves courses by instructor
func (r *CourseRepository) ListByInstructor(instructorID string, limit, offset int) ([]*domain.Course, error) {
	query := `
		SELECT c.id, c.tenant_id, c.instructor_id, c.category_id, c.title, c.slug, c.description, c.thumbnail_url,
		       c.price, c.discount_price, c.discount_valid_until, c.currency, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''), c.is_published, c.is_featured, c.created_at, c.updated_at,
			   u.full_name, cat.name
		FROM courses c
		LEFT JOIN users u ON c.instructor_id = u.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.instructor_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	return r.scanCoursesWithDetails(query, instructorID, limit, offset)
}

// ListPublished retrieves published courses for a tenant
func (r *CourseRepository) ListPublished(tenantID string, limit, offset int) ([]*domain.Course, error) {
	var query string
	var args []interface{}
	
	// If no tenant specified or "default", get all published courses
	if tenantID == "" || tenantID == "default" {
		query = `
			SELECT c.id, c.tenant_id, c.instructor_id, c.category_id, c.title, c.slug, c.description, c.thumbnail_url,
			       c.price, c.discount_price, c.discount_valid_until, c.currency, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''), c.is_published, c.is_featured, c.created_at, c.updated_at,
				   u.full_name, cat.name
			FROM courses c
			LEFT JOIN users u ON c.instructor_id = u.id
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.is_published = true
			ORDER BY c.is_featured DESC, c.created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	} else {
		query = `
			SELECT c.id, c.tenant_id, c.instructor_id, c.category_id, c.title, c.slug, c.description, c.thumbnail_url,
			       c.price, c.currency, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''), c.is_published, c.is_featured, c.created_at, c.updated_at,
				   u.full_name, cat.name
			FROM courses c
			LEFT JOIN users u ON c.instructor_id = u.id
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.tenant_id = $1 AND c.is_published = true
			ORDER BY c.is_featured DESC, c.created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{tenantID, limit, offset}
	}
	
	return r.scanCoursesWithDetails(query, args...)
}

// CountByTenant returns total course count for a tenant
func (r *CourseRepository) CountByTenant(tenantID string) (int, error) {
	var query string
	var args []interface{}
	
	if tenantID == "" || tenantID == "default" {
		query = `SELECT COUNT(*) FROM courses WHERE tenant_id IS NULL`
		args = []interface{}{}
	} else {
		query = `SELECT COUNT(*) FROM courses WHERE tenant_id = $1`
		args = []interface{}{tenantID}
	}
	
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// CountByInstructor returns total course count for an instructor
func (r *CourseRepository) CountByInstructor(instructorID string) (int, error) {
	query := `SELECT COUNT(*) FROM courses WHERE instructor_id = $1`
	var count int
	err := r.db.QueryRow(query, instructorID).Scan(&count)
	return count, err
}

// GetStudentCountByInstructor returns total enrolled students across all courses of an instructor
func (r *CourseRepository) GetStudentCountByInstructor(instructorID string) (int, error) {
	query := `
		SELECT COALESCE(COUNT(DISTINCT e.user_id), 0)
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE c.instructor_id = $1
	`
	var count int
	err := r.db.QueryRow(query, instructorID).Scan(&count)
	return count, err
}

// GetAverageRatingByInstructor returns average rating across all courses of an instructor
func (r *CourseRepository) GetAverageRatingByInstructor(instructorID string) (float64, error) {
	query := `
		SELECT COALESCE(AVG(cr.rating), 0)
		FROM course_ratings cr
		JOIN courses c ON cr.course_id = c.id
		WHERE c.instructor_id = $1
	`
	var rating float64
	err := r.db.QueryRow(query, instructorID).Scan(&rating)
	return rating, err
}

// Helper function to scan course rows with details
func (r *CourseRepository) scanCoursesWithDetails(query string, args ...interface{}) ([]*domain.Course, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var courses []*domain.Course
	for rows.Next() {
		var course domain.Course
		var tenantID, instructorID, categoryID, thumbnailURL, instructorName, categoryName sql.NullString
		var discountPrice sql.NullFloat64
		var discountValidUntil sql.NullTime
		
		err := rows.Scan(
			&course.ID, &tenantID, &instructorID, &categoryID, &course.Title, &course.Slug,
			&course.Description, &thumbnailURL, &course.Price, &discountPrice, &discountValidUntil,
			&course.Currency, &course.LessonsCount, &course.Duration, &course.IsPublished, &course.IsFeatured, &course.CreatedAt, &course.UpdatedAt,
			&instructorName, &categoryName,
		)
		if err != nil {
			return nil, err
		}
		
		if tenantID.Valid {
			course.TenantID = tenantID.String
		}
		if instructorID.Valid {
			course.InstructorID = &instructorID.String
		}
		if categoryID.Valid {
			course.CategoryID = &categoryID.String
		}
		if thumbnailURL.Valid {
			course.ThumbnailURL = &thumbnailURL.String
		}
		if discountPrice.Valid {
			course.DiscountPrice = &discountPrice.Float64
		}
		if discountValidUntil.Valid {
			course.DiscountValidUntil = &discountValidUntil.Time
		}
		
		// Populate related objects
		if instructorName.Valid {
			course.Instructor = &domain.User{
				ID:       *course.InstructorID,
				FullName: instructorName.String,
			}
		}
		if categoryName.Valid {
			course.Category = &domain.Category{
				ID:   *course.CategoryID,
				Name: categoryName.String,
			}
		}
		
		courses = append(courses, &course)
	}
	
	return courses, nil
}

// generateSlug creates URL-friendly slug from title
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
