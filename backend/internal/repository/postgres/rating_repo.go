package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// RatingRepository handles course rating data access
type RatingRepository struct {
	db *sqlx.DB
}

// NewRatingRepository creates a new rating repository
func NewRatingRepository(db *sqlx.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

// Create creates a new course rating
func (r *RatingRepository) Create(rating *domain.CourseRating) error {
	query := `
		INSERT INTO course_ratings (user_id, course_id, rating, review)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		rating.UserID,
		rating.CourseID,
		rating.Rating,
		rating.Review,
	).Scan(&rating.ID, &rating.CreatedAt, &rating.UpdatedAt)
}

// Update updates an existing rating
func (r *RatingRepository) Update(rating *domain.CourseRating) error {
	query := `
		UPDATE course_ratings 
		SET rating = $1, review = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRow(query,
		rating.Rating,
		rating.Review,
		rating.ID,
	).Scan(&rating.UpdatedAt)
}

// GetByCourse retrieves ratings for a course
func (r *RatingRepository) GetByCourse(courseID string, limit int) ([]*domain.CourseRating, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	ratings := []*domain.CourseRating{}
	err := r.db.Select(&ratings, `
		SELECT 
			cr.id, cr.user_id, cr.course_id, cr.rating, cr.review, cr.created_at, cr.updated_at,
			COALESCE(u.full_name, u.email) as user_name
		FROM course_ratings cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.course_id = $1
		ORDER BY cr.created_at DESC
		LIMIT $2
	`, courseID, limit)

	if err == sql.ErrNoRows {
		return ratings, nil
	}
	return ratings, err
}

// GetByUserAndCourse retrieves a user's rating for a specific course
func (r *RatingRepository) GetByUserAndCourse(userID, courseID string) (*domain.CourseRating, error) {
	var rating domain.CourseRating
	err := r.db.Get(&rating, `
		SELECT 
			cr.id, cr.user_id, cr.course_id, cr.rating, cr.review, cr.created_at, cr.updated_at,
			COALESCE(u.full_name, u.email) as user_name
		FROM course_ratings cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.user_id = $1 AND cr.course_id = $2
	`, userID, courseID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetByID retrieves a rating by ID
func (r *RatingRepository) GetByID(id string) (*domain.CourseRating, error) {
	var rating domain.CourseRating
	err := r.db.Get(&rating, `
		SELECT 
			cr.id, cr.user_id, cr.course_id, cr.rating, cr.review, cr.created_at, cr.updated_at,
			COALESCE(u.full_name, u.email) as user_name
		FROM course_ratings cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.id = $1
	`, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// GetCourseStats retrieves rating statistics for a course
func (r *RatingRepository) GetCourseStats(courseID string) (*domain.CourseRatingStats, error) {
	var stats domain.CourseRatingStats
	err := r.db.Get(&stats, `
		SELECT 
			COALESCE(AVG(rating)::numeric(3,2), 0) as average_rating,
			COUNT(*) as total_ratings
		FROM course_ratings
		WHERE course_id = $1
	`, courseID)

	if err != nil {
		return &domain.CourseRatingStats{}, nil
	}
	return &stats, nil
}

// Delete deletes a rating
func (r *RatingRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM course_ratings WHERE id = $1`, id)
	return err
}

// Exists checks if a rating exists for user and course
func (r *RatingRepository) Exists(userID, courseID string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM course_ratings 
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID)
	return exists, err
}

// GetRecentByCourse retrieves recent ratings for a course with time filter
func (r *RatingRepository) GetRecentByCourse(courseID string, since time.Time, limit int) ([]*domain.CourseRating, error) {
	ratings := []*domain.CourseRating{}
	err := r.db.Select(&ratings, `
		SELECT 
			cr.id, cr.user_id, cr.course_id, cr.rating, cr.review, cr.created_at, cr.updated_at,
			COALESCE(u.full_name, u.email) as user_name
		FROM course_ratings cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.course_id = $1 AND cr.created_at >= $2
		ORDER BY cr.created_at DESC
		LIMIT $3
	`, courseID, since, limit)

	return ratings, err
}

// GetAll retrieves all ratings with pagination (for admin)
func (r *RatingRepository) GetAll(limit, offset int) ([]*domain.CourseRating, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	ratings := []*domain.CourseRating{}
	err := r.db.Select(&ratings, `
		SELECT 
			cr.id, cr.user_id, cr.course_id, cr.rating, cr.review, cr.created_at, cr.updated_at,
			COALESCE(u.full_name, u.email) as user_name
		FROM course_ratings cr
		JOIN users u ON u.id = cr.user_id
		ORDER BY cr.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err == sql.ErrNoRows {
		return ratings, nil
	}
	return ratings, err
}

// Count returns total number of ratings
func (r *RatingRepository) Count() (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM course_ratings`)
	return count, err
}

// GetAverageAll returns the average rating across all courses
func (r *RatingRepository) GetAverageAll() (float64, error) {
	var avg float64
	err := r.db.Get(&avg, `SELECT COALESCE(AVG(rating)::numeric(3,2), 0) FROM course_ratings`)
	return avg, err
}

// GetBulkCourseStats retrieves rating stats for multiple courses
func (r *RatingRepository) GetBulkCourseStats(courseIDs []string) (map[string]*domain.CourseRatingStats, error) {
	if len(courseIDs) == 0 {
		return make(map[string]*domain.CourseRatingStats), nil
	}

	type courseStats struct {
		CourseID      string  `db:"course_id"`
		AverageRating float64 `db:"average_rating"`
		TotalRatings  int     `db:"total_ratings"`
	}

	var stats []courseStats
	query, args, err := sqlx.In(`
		SELECT 
			course_id,
			COALESCE(AVG(rating)::numeric(3,2), 0) as average_rating,
			COUNT(*) as total_ratings
		FROM course_ratings
		WHERE course_id IN (?)
		GROUP BY course_id
	`, courseIDs)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.Select(&stats, query, args...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*domain.CourseRatingStats)
	for _, s := range stats {
		result[s.CourseID] = &domain.CourseRatingStats{
			AverageRating: s.AverageRating,
			TotalRatings:  s.TotalRatings,
		}
	}
	return result, nil
}

