package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Enrollment represents a user's enrollment in a course
type Enrollment struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	CourseID           string     `json:"course_id"`
	TransactionID      *string    `json:"transaction_id,omitempty"`
	ProgressPercentage int        `json:"progress_percentage"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	EnrolledAt         time.Time  `json:"enrolled_at"`
}

// EnrollmentRepository handles enrollment data access
type EnrollmentRepository struct {
	db *sqlx.DB
}

// NewEnrollmentRepository creates a new EnrollmentRepository
func NewEnrollmentRepository(db *sqlx.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// GetByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetByID(id string) (*Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, transaction_id, progress_percentage,
		       completed_at, enrolled_at
		FROM enrollments WHERE id = $1
	`
	
	var enrollment Enrollment
	var transactionID sql.NullString
	var completedAt sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID,
		&transactionID, &enrollment.ProgressPercentage,
		&completedAt, &enrollment.EnrolledAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if transactionID.Valid {
		enrollment.TransactionID = &transactionID.String
	}
	if completedAt.Valid {
		enrollment.CompletedAt = &completedAt.Time
	}
	
	return &enrollment, nil
}

// GetByUserAndCourse checks if user is enrolled in a course
func (r *EnrollmentRepository) GetByUserAndCourse(userID, courseID string) (*Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, transaction_id, progress_percentage,
		       completed_at, enrolled_at
		FROM enrollments WHERE user_id = $1 AND course_id = $2
	`
	
	var enrollment Enrollment
	var transactionID sql.NullString
	var completedAt sql.NullTime
	
	err := r.db.QueryRow(query, userID, courseID).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.CourseID,
		&transactionID, &enrollment.ProgressPercentage,
		&completedAt, &enrollment.EnrolledAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if transactionID.Valid {
		enrollment.TransactionID = &transactionID.String
	}
	if completedAt.Valid {
		enrollment.CompletedAt = &completedAt.Time
	}
	
	return &enrollment, nil
}

// Create inserts a new enrollment
func (r *EnrollmentRepository) Create(enrollment *Enrollment) error {
	query := `
		INSERT INTO enrollments (user_id, course_id, transaction_id, progress_percentage, enrolled_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	enrollment.EnrolledAt = time.Now()
	enrollment.ProgressPercentage = 0
	
	return r.db.QueryRow(query,
		enrollment.UserID, enrollment.CourseID, enrollment.TransactionID,
		enrollment.ProgressPercentage, enrollment.EnrolledAt,
	).Scan(&enrollment.ID)
}

// UpdateProgress updates the progress percentage of an enrollment
func (r *EnrollmentRepository) UpdateProgress(id string, progressPercentage int) error {
	query := `UPDATE enrollments SET progress_percentage = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, progressPercentage)
	return err
}

// MarkCompleted marks an enrollment as completed
func (r *EnrollmentRepository) MarkCompleted(id string) error {
	query := `UPDATE enrollments SET progress_percentage = 100, completed_at = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, time.Now())
	return err
}

// ListByUser retrieves all enrollments for a user
func (r *EnrollmentRepository) ListByUser(userID string, limit, offset int) ([]*Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, transaction_id, progress_percentage,
		       completed_at, enrolled_at
		FROM enrollments 
		WHERE user_id = $1
		ORDER BY enrolled_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanEnrollments(rows)
}

// ListByCourse retrieves all enrollments for a course
func (r *EnrollmentRepository) ListByCourse(courseID string, limit, offset int) ([]*Enrollment, error) {
	query := `
		SELECT id, user_id, course_id, transaction_id, progress_percentage,
		       completed_at, enrolled_at
		FROM enrollments 
		WHERE course_id = $1
		ORDER BY enrolled_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, courseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanEnrollments(rows)
}

// CountByUser returns total enrollment count for a user
func (r *EnrollmentRepository) CountByUser(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM enrollments WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// CountCompletedByUser returns completed course count for a user
func (r *EnrollmentRepository) CountCompletedByUser(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM enrollments WHERE user_id = $1 AND completed_at IS NOT NULL`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// Delete removes an enrollment
func (r *EnrollmentRepository) Delete(id string) error {
	query := `DELETE FROM enrollments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateProgressByUserAndCourse updates progress for a user's enrollment in a course
func (r *EnrollmentRepository) UpdateProgressByUserAndCourse(userID, courseID string, progressPercentage int) error {
	if progressPercentage >= 100 {
		query := `UPDATE enrollments SET progress_percentage = 100, completed_at = $3 
		          WHERE user_id = $1 AND course_id = $2`
		_, err := r.db.Exec(query, userID, courseID, time.Now())
		return err
	}
	query := `UPDATE enrollments SET progress_percentage = $3 WHERE user_id = $1 AND course_id = $2`
	_, err := r.db.Exec(query, userID, courseID, progressPercentage)
	return err
}

// Helper to scan enrollment rows
func (r *EnrollmentRepository) scanEnrollments(rows *sql.Rows) ([]*Enrollment, error) {
	var enrollments []*Enrollment
	for rows.Next() {
		var enrollment Enrollment
		var transactionID sql.NullString
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&enrollment.ID, &enrollment.UserID, &enrollment.CourseID,
			&transactionID, &enrollment.ProgressPercentage,
			&completedAt, &enrollment.EnrolledAt,
		)
		if err != nil {
			return nil, err
		}
		
		if transactionID.Valid {
			enrollment.TransactionID = &transactionID.String
		}
		if completedAt.Valid {
			enrollment.CompletedAt = &completedAt.Time
		}
		
		enrollments = append(enrollments, &enrollment)
	}
	
	return enrollments, nil
}
