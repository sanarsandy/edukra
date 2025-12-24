package domain

import "time"

// LessonProgress represents user's progress on a specific lesson
type LessonProgress struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	LessonID    string     `json:"lesson_id" db:"lesson_id"`
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	WatchTime   int        `json:"watch_time" db:"watch_time"` // seconds
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// ActivityLog represents a user learning activity
type ActivityLog struct {
	ID            string                 `json:"id" db:"id"`
	UserID        string                 `json:"user_id" db:"user_id"`
	ActivityType  string                 `json:"activity_type" db:"activity_type"` // lesson_complete, quiz_pass, course_complete, enroll
	ReferenceID   *string                `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType *string                `json:"reference_type,omitempty" db:"reference_type"` // lesson, quiz, course
	Description   string                 `json:"description" db:"description"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// Activity Types
const (
	ActivityLessonComplete = "lesson_complete"
	ActivityQuizPass       = "quiz_pass"
	ActivityQuizFail       = "quiz_fail"
	ActivityCourseComplete = "course_complete"
	ActivityEnroll         = "enroll"
	ActivityVideoWatch     = "video_watch"
)

// CourseProgress aggregates lesson progress for a course
type CourseProgress struct {
	CourseID           string            `json:"course_id"`
	TotalLessons       int               `json:"total_lessons"`
	CompletedLessons   int               `json:"completed_lessons"`
	ProgressPercentage int               `json:"progress_percentage"`
	TotalWatchTime     int               `json:"total_watch_time"` // seconds
	LessonProgress     []*LessonProgress `json:"lesson_progress,omitempty"`
}

// LearningStreak tracks consecutive learning days
type LearningStreak struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"user_id" db:"user_id"`
	CurrentStreak    int        `json:"current_streak" db:"current_streak"`
	LongestStreak    int        `json:"longest_streak" db:"longest_streak"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// UserLearningStats aggregates all learning statistics for a user
type UserLearningStats struct {
	TotalWatchTime   int `json:"total_watch_time"`   // seconds
	QuizzesPassed    int `json:"quizzes_passed"`
	LessonsCompleted int `json:"lessons_completed"`
	CoursesCompleted int `json:"courses_completed"`
	CurrentStreak    int `json:"current_streak"`
}

// UpdateLessonProgressRequest for API request
type UpdateLessonProgressRequest struct {
	IsCompleted bool `json:"is_completed"`
	WatchTime   int  `json:"watch_time,omitempty"`
}

// UpdateWatchTimeRequest for API request
type UpdateWatchTimeRequest struct {
	Seconds int `json:"seconds"`
}

// BulkProgressRequest for syncing multiple lessons at once
type BulkProgressRequest struct {
	LessonIDs []string `json:"lesson_ids"`
}
