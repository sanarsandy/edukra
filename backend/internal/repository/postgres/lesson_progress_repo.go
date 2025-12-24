package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// LessonProgressRepository handles lesson progress data access
type LessonProgressRepository struct {
	db *sqlx.DB
}

// NewLessonProgressRepository creates a new lesson progress repository
func NewLessonProgressRepository(db *sqlx.DB) *LessonProgressRepository {
	return &LessonProgressRepository{db: db}
}

// GetByUserAndLesson retrieves progress for a specific lesson
func (r *LessonProgressRepository) GetByUserAndLesson(userID, lessonID string) (*domain.LessonProgress, error) {
	var progress domain.LessonProgress
	err := r.db.Get(&progress, `
		SELECT id, user_id, lesson_id, is_completed, watch_time, completed_at, created_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`, userID, lessonID)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// Upsert creates or updates lesson progress
func (r *LessonProgressRepository) Upsert(progress *domain.LessonProgress) error {
	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, is_completed, watch_time, completed_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			is_completed = EXCLUDED.is_completed,
			watch_time = CASE 
				WHEN EXCLUDED.watch_time > lesson_progress.watch_time THEN EXCLUDED.watch_time 
				ELSE lesson_progress.watch_time 
			END,
			completed_at = COALESCE(lesson_progress.completed_at, EXCLUDED.completed_at),
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, 
		progress.UserID, 
		progress.LessonID, 
		progress.IsCompleted, 
		progress.WatchTime,
		progress.CompletedAt,
	).Scan(&progress.ID, &progress.CreatedAt, &progress.UpdatedAt)
}

// MarkComplete marks a lesson as completed
func (r *LessonProgressRepository) MarkComplete(userID, lessonID string) error {
	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, is_completed, completed_at)
		VALUES ($1, $2, true, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			is_completed = true,
			completed_at = COALESCE(lesson_progress.completed_at, CURRENT_TIMESTAMP),
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, lessonID)
	return err
}

// UpdateWatchTime increments watch time for a lesson
func (r *LessonProgressRepository) UpdateWatchTime(userID, lessonID string, seconds int) error {
	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, watch_time)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			watch_time = lesson_progress.watch_time + EXCLUDED.watch_time,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, lessonID, seconds)
	return err
}

// ListByUserAndCourse retrieves all lesson progress for a course
func (r *LessonProgressRepository) ListByUserAndCourse(userID, courseID string) ([]*domain.LessonProgress, error) {
	var progress []*domain.LessonProgress
	err := r.db.Select(&progress, `
		SELECT lp.id, lp.user_id, lp.lesson_id, lp.is_completed, lp.watch_time, lp.completed_at, lp.created_at, lp.updated_at
		FROM lesson_progress lp
		JOIN lessons l ON l.id = lp.lesson_id
		WHERE lp.user_id = $1 AND l.course_id = $2 AND l.is_container = false
		ORDER BY l.order_index
	`, userID, courseID)
	
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// GetCourseProgress calculates course-level progress
func (r *LessonProgressRepository) GetCourseProgress(userID, courseID string) (*domain.CourseProgress, error) {
	// Get total lessons in course (exclude containers/modules)
	var totalLessons int
	err := r.db.Get(&totalLessons, `SELECT COUNT(*) FROM lessons WHERE course_id = $1 AND is_container = false`, courseID)
	if err != nil {
		return nil, err
	}

	// Get completed lessons and total watch time (exclude containers)
	var completedLessons, totalWatchTime int
	err = r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN lp.is_completed THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(lp.watch_time), 0)
		FROM lesson_progress lp
		JOIN lessons l ON l.id = lp.lesson_id
		WHERE lp.user_id = $1 AND l.course_id = $2 AND l.is_container = false
	`, userID, courseID).Scan(&completedLessons, &totalWatchTime)
	if err != nil {
		return nil, err
	}

	// Get individual lesson progress
	lessonProgress, err := r.ListByUserAndCourse(userID, courseID)
	if err != nil {
		return nil, err
	}

	progressPercentage := 0
	if totalLessons > 0 {
		progressPercentage = (completedLessons * 100) / totalLessons
	}

	return &domain.CourseProgress{
		CourseID:           courseID,
		TotalLessons:       totalLessons,
		CompletedLessons:   completedLessons,
		ProgressPercentage: progressPercentage,
		TotalWatchTime:     totalWatchTime,
		LessonProgress:     lessonProgress,
	}, nil
}

// BulkMarkComplete marks multiple lessons as completed
func (r *LessonProgressRepository) BulkMarkComplete(userID string, lessonIDs []string) error {
	if len(lessonIDs) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, is_completed, completed_at)
		VALUES ($1, $2, true, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			is_completed = true,
			completed_at = COALESCE(lesson_progress.completed_at, CURRENT_TIMESTAMP),
			updated_at = CURRENT_TIMESTAMP
	`

	for _, lessonID := range lessonIDs {
		if _, err := tx.Exec(query, userID, lessonID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetTotalWatchTime returns total watch time for a user
func (r *LessonProgressRepository) GetTotalWatchTime(userID string) (int, error) {
	var total int
	err := r.db.Get(&total, `SELECT COALESCE(SUM(watch_time), 0) FROM lesson_progress WHERE user_id = $1`, userID)
	return total, err
}

// GetCompletedLessonIDs returns list of completed lesson IDs for a course
func (r *LessonProgressRepository) GetCompletedLessonIDs(userID, courseID string) ([]string, error) {
	var ids []string
	err := r.db.Select(&ids, `
		SELECT lp.lesson_id
		FROM lesson_progress lp
		JOIN lessons l ON l.id = lp.lesson_id
		WHERE lp.user_id = $1 AND l.course_id = $2 AND lp.is_completed = true
	`, userID, courseID)
	return ids, err
}
