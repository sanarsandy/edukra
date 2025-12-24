package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// ActivityLogRepository handles activity log data access
type ActivityLogRepository struct {
	db *sqlx.DB
}

// NewActivityLogRepository creates a new activity log repository
func NewActivityLogRepository(db *sqlx.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

// Create adds a new activity log entry
func (r *ActivityLogRepository) Create(log *domain.ActivityLog) error {
	metadata, err := json.Marshal(log.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO activity_logs (user_id, activity_type, reference_id, reference_type, description, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query,
		log.UserID,
		log.ActivityType,
		log.ReferenceID,
		log.ReferenceType,
		log.Description,
		metadata,
	).Scan(&log.ID, &log.CreatedAt)
}

// ListByUser retrieves recent activities for a user
func (r *ActivityLogRepository) ListByUser(userID string, limit int) ([]*domain.ActivityLog, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	rows, err := r.db.Query(`
		SELECT id, user_id, activity_type, reference_id, reference_type, description, metadata, created_at
		FROM activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.ActivityLog
	for rows.Next() {
		var log domain.ActivityLog
		var metadata []byte
		var refID, refType sql.NullString

		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.ActivityType,
			&refID,
			&refType,
			&log.Description,
			&metadata,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}

		if refID.Valid {
			log.ReferenceID = &refID.String
		}
		if refType.Valid {
			log.ReferenceType = &refType.String
		}
		if len(metadata) > 0 {
			json.Unmarshal(metadata, &log.Metadata)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// LogLessonComplete logs a lesson completion activity
func (r *ActivityLogRepository) LogLessonComplete(userID, lessonID, lessonTitle, courseName string) error {
	log := &domain.ActivityLog{
		UserID:        userID,
		ActivityType:  domain.ActivityLessonComplete,
		ReferenceID:   &lessonID,
		ReferenceType: strPtr("lesson"),
		Description:   "Menyelesaikan materi \"" + lessonTitle + "\"",
		Metadata: map[string]interface{}{
			"lesson_title": lessonTitle,
			"course_name":  courseName,
		},
	}
	return r.Create(log)
}

// LogQuizResult logs a quiz attempt result
func (r *ActivityLogRepository) LogQuizResult(userID, quizID, quizTitle string, score int, passed bool) error {
	activityType := domain.ActivityQuizFail
	description := "Tidak lulus kuis dengan nilai " + string(rune(score)) + "%"
	if passed {
		activityType = domain.ActivityQuizPass
		description = "Lulus kuis \"" + quizTitle + "\" dengan nilai " + itoa(score) + "%"
	}

	log := &domain.ActivityLog{
		UserID:        userID,
		ActivityType:  activityType,
		ReferenceID:   &quizID,
		ReferenceType: strPtr("quiz"),
		Description:   description,
		Metadata: map[string]interface{}{
			"quiz_title": quizTitle,
			"score":      score,
			"passed":     passed,
		},
	}
	return r.Create(log)
}

// LogCourseComplete logs a course completion
func (r *ActivityLogRepository) LogCourseComplete(userID, courseID, courseTitle string) error {
	log := &domain.ActivityLog{
		UserID:        userID,
		ActivityType:  domain.ActivityCourseComplete,
		ReferenceID:   &courseID,
		ReferenceType: strPtr("course"),
		Description:   "Menyelesaikan kursus \"" + courseTitle + "\"",
		Metadata: map[string]interface{}{
			"course_title": courseTitle,
		},
	}
	return r.Create(log)
}

// LogEnroll logs an enrollment activity
func (r *ActivityLogRepository) LogEnroll(userID, courseID, courseTitle string) error {
	log := &domain.ActivityLog{
		UserID:        userID,
		ActivityType:  domain.ActivityEnroll,
		ReferenceID:   &courseID,
		ReferenceType: strPtr("course"),
		Description:   "Mendaftar kursus \"" + courseTitle + "\"",
		Metadata: map[string]interface{}{
			"course_title": courseTitle,
		},
	}
	return r.Create(log)
}

// GetUserStats returns aggregated stats for a user
func (r *ActivityLogRepository) GetUserStats(userID string) (*domain.UserLearningStats, error) {
	stats := &domain.UserLearningStats{}

	// Count quizzes passed
	err := r.db.Get(&stats.QuizzesPassed, `
		SELECT COUNT(*) FROM activity_logs 
		WHERE user_id = $1 AND activity_type = $2
	`, userID, domain.ActivityQuizPass)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Count lessons completed
	err = r.db.Get(&stats.LessonsCompleted, `
		SELECT COUNT(*) FROM activity_logs 
		WHERE user_id = $1 AND activity_type = $2
	`, userID, domain.ActivityLessonComplete)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Count courses completed
	err = r.db.Get(&stats.CoursesCompleted, `
		SELECT COUNT(*) FROM activity_logs 
		WHERE user_id = $1 AND activity_type = $2
	`, userID, domain.ActivityCourseComplete)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return stats, nil
}

// UpdateStreak updates the learning streak for a user
func (r *ActivityLogRepository) UpdateStreak(userID string) error {
	today := time.Now().Format("2006-01-02")
	
	// Get current streak info
	var currentStreak, longestStreak int
	var lastDate sql.NullString
	
	err := r.db.QueryRow(`
		SELECT current_streak, longest_streak, last_activity_date 
		FROM learning_streaks 
		WHERE user_id = $1
	`, userID).Scan(&currentStreak, &longestStreak, &lastDate)
	
	if err == sql.ErrNoRows {
		// First activity, create record
		_, err = r.db.Exec(`
			INSERT INTO learning_streaks (user_id, current_streak, longest_streak, last_activity_date)
			VALUES ($1, 1, 1, $2)
		`, userID, today)
		return err
	}
	if err != nil {
		return err
	}

	// Already logged today
	if lastDate.Valid && lastDate.String == today {
		return nil
	}

	// Check if streak continues
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if lastDate.Valid && lastDate.String == yesterday {
		currentStreak++
		if currentStreak > longestStreak {
			longestStreak = currentStreak
		}
	} else {
		// Streak broken
		currentStreak = 1
	}

	_, err = r.db.Exec(`
		UPDATE learning_streaks 
		SET current_streak = $1, longest_streak = $2, last_activity_date = $3, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $4
	`, currentStreak, longestStreak, today, userID)
	
	return err
}

// GetStreak returns the current streak for a user
func (r *ActivityLogRepository) GetStreak(userID string) (int, error) {
	var streak int
	err := r.db.Get(&streak, `SELECT COALESCE(current_streak, 0) FROM learning_streaks WHERE user_id = $1`, userID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return streak, err
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}
