package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// LessonRepository implements domain.LessonRepository using PostgreSQL
type LessonRepository struct {
	db *sqlx.DB
}

// NewLessonRepository creates a new LessonRepository
func NewLessonRepository(db *sqlx.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// scanLesson is a helper function to scan a lesson row
func (r *LessonRepository) scanLesson(row interface{ Scan(...interface{}) error }) (*domain.Lesson, error) {
	var lesson domain.Lesson
	var videoURL, content, parentID sql.NullString
	var isContainer sql.NullBool

	err := row.Scan(
		&lesson.ID, &lesson.CourseID, &parentID, &isContainer, &lesson.Title, &lesson.Description,
		&lesson.OrderIndex, &lesson.ContentType, &videoURL, &content, &lesson.VideoDuration,
		&lesson.SecurityLevel, &lesson.IsPreview, &lesson.CreatedAt, &lesson.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if videoURL.Valid {
		lesson.VideoURL = &videoURL.String
	}
	if content.Valid {
		lesson.Content = &content.String
	}
	if parentID.Valid {
		lesson.ParentID = &parentID.String
	}
	if isContainer.Valid {
		lesson.IsContainer = isContainer.Bool
	}

	return &lesson, nil
}

// GetByID retrieves a lesson by ID
func (r *LessonRepository) GetByID(id string) (*domain.Lesson, error) {
	query := `
		SELECT id, course_id, parent_id, is_container, title, description, order_index, content_type,
		       video_url, content, video_duration, security_level, is_preview, created_at, updated_at
		FROM lessons WHERE id = $1
	`

	lesson, err := r.scanLesson(r.db.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return lesson, err
}

// Create inserts a new lesson
func (r *LessonRepository) Create(lesson *domain.Lesson) error {
	query := `
		INSERT INTO lessons (course_id, parent_id, is_container, title, description, order_index, content_type,
		                     video_url, content, video_duration, security_level, is_preview,
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	now := time.Now()
	lesson.CreatedAt = now
	lesson.UpdatedAt = now

	if lesson.SecurityLevel == "" {
		lesson.SecurityLevel = domain.SecuritySignedURL
	}

	return r.db.QueryRow(query,
		lesson.CourseID, lesson.ParentID, lesson.IsContainer, lesson.Title, lesson.Description, lesson.OrderIndex,
		lesson.ContentType, lesson.VideoURL, lesson.Content, lesson.VideoDuration,
		lesson.SecurityLevel, lesson.IsPreview, lesson.CreatedAt, lesson.UpdatedAt,
	).Scan(&lesson.ID)
}

// Update updates an existing lesson
func (r *LessonRepository) Update(lesson *domain.Lesson) error {
	query := `
		UPDATE lessons SET
			parent_id = $2, is_container = $3, title = $4, description = $5, order_index = $6, content_type = $7,
			video_url = $8, content = $9, video_duration = $10, security_level = $11, 
			is_preview = $12, updated_at = $13
		WHERE id = $1
	`

	lesson.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		lesson.ID, lesson.ParentID, lesson.IsContainer, lesson.Title, lesson.Description, lesson.OrderIndex,
		lesson.ContentType, lesson.VideoURL, lesson.Content, lesson.VideoDuration,
		lesson.SecurityLevel, lesson.IsPreview, lesson.UpdatedAt,
	)
	return err
}

// Delete removes a lesson by ID (CASCADE will delete children)
func (r *LessonRepository) Delete(id string) error {
	query := `DELETE FROM lessons WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListByCourse retrieves all lessons for a course ordered by index (flat list)
func (r *LessonRepository) ListByCourse(courseID string) ([]*domain.Lesson, error) {
	query := `
		SELECT id, course_id, parent_id, is_container, title, description, order_index, content_type,
		       video_url, content, video_duration, security_level, is_preview, created_at, updated_at
		FROM lessons 
		WHERE course_id = $1
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		lesson, err := r.scanLesson(rows)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// ListByParent retrieves lessons by parent_id (nil for root level)
func (r *LessonRepository) ListByParent(courseID string, parentID *string) ([]*domain.Lesson, error) {
	var query string
	var rows *sql.Rows
	var err error

	if parentID == nil {
		query = `
			SELECT id, course_id, parent_id, is_container, title, description, order_index, content_type,
			       video_url, content, video_duration, security_level, is_preview, created_at, updated_at
			FROM lessons 
			WHERE course_id = $1 AND parent_id IS NULL
			ORDER BY order_index ASC
		`
		rows, err = r.db.Query(query, courseID)
	} else {
		query = `
			SELECT id, course_id, parent_id, is_container, title, description, order_index, content_type,
			       video_url, content, video_duration, security_level, is_preview, created_at, updated_at
			FROM lessons 
			WHERE course_id = $1 AND parent_id = $2
			ORDER BY order_index ASC
		`
		rows, err = r.db.Query(query, courseID, *parentID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		lesson, err := r.scanLesson(rows)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// GetTree returns all lessons for a course as a nested tree structure
func (r *LessonRepository) GetTree(courseID string) ([]*domain.Lesson, error) {
	// First, get all lessons flat
	allLessons, err := r.ListByCourse(courseID)
	if err != nil {
		return nil, err
	}

	// Build a map for quick lookup
	lessonMap := make(map[string]*domain.Lesson)
	for _, lesson := range allLessons {
		lesson.Children = []*domain.Lesson{} // Initialize empty children slice
		lessonMap[lesson.ID] = lesson
	}

	// Build tree structure
	var rootLessons []*domain.Lesson
	for _, lesson := range allLessons {
		if lesson.ParentID == nil {
			rootLessons = append(rootLessons, lesson)
		} else {
			parent, exists := lessonMap[*lesson.ParentID]
			if exists {
				parent.Children = append(parent.Children, lesson)
			} else {
				// Orphaned lesson, add to root
				rootLessons = append(rootLessons, lesson)
			}
		}
	}

	return rootLessons, nil
}

// ReorderLessons updates the order of lessons within the same parent
func (r *LessonRepository) ReorderLessons(courseID string, lessonIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE lessons SET order_index = $1, updated_at = $2 WHERE id = $3 AND course_id = $4`
	now := time.Now()

	for i, lessonID := range lessonIDs {
		_, err := tx.Exec(query, i, now, lessonID, courseID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// MoveLesson moves a lesson to a new parent with a new order
func (r *LessonRepository) MoveLesson(lessonID string, newParentID *string, newOrderIndex int) error {
	query := `
		UPDATE lessons SET 
			parent_id = $2, 
			order_index = $3, 
			updated_at = $4
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.db.Exec(query, lessonID, newParentID, newOrderIndex, now)
	return err
}

// CountByCourse returns total lesson count for a course
func (r *LessonRepository) CountByCourse(courseID string) (int, error) {
	query := `SELECT COUNT(*) FROM lessons WHERE course_id = $1`
	var count int
	err := r.db.QueryRow(query, courseID).Scan(&count)
	return count, err
}

// CountByParent returns lesson count for a specific parent (nil for root level)
func (r *LessonRepository) CountByParent(courseID string, parentID *string) (int, error) {
	var query string
	var count int
	var err error

	if parentID == nil {
		query = `SELECT COUNT(*) FROM lessons WHERE course_id = $1 AND parent_id IS NULL`
		err = r.db.QueryRow(query, courseID).Scan(&count)
	} else {
		query = `SELECT COUNT(*) FROM lessons WHERE course_id = $1 AND parent_id = $2`
		err = r.db.QueryRow(query, courseID, *parentID).Scan(&count)
	}

	return count, err
}

// GetPreviewLessons returns lessons marked as preview for a course
func (r *LessonRepository) GetPreviewLessons(courseID string) ([]*domain.Lesson, error) {
	query := `
		SELECT id, course_id, parent_id, is_container, title, description, order_index, content_type,
		       video_url, content, video_duration, security_level, is_preview, created_at, updated_at
		FROM lessons 
		WHERE course_id = $1 AND is_preview = true
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		lesson, err := r.scanLesson(rows)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}
