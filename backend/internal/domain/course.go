package domain

import (
	"time"
)

// Course represents a course in the LMS
type Course struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	InstructorID *string   `json:"instructor_id,omitempty"`
	CategoryID   *string   `json:"category_id,omitempty"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	LessonsCount int       `json:"lessons_count"`
	Duration     string    `json:"duration"`
	IsPublished  bool      `json:"is_published"`
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Related data (populated on demand)
	Instructor *User     `json:"instructor,omitempty"`
	Category   *Category `json:"category,omitempty"`
	Lessons    []*Lesson `json:"lessons,omitempty"`
}

// CreateCourseRequest represents a request to create a course
type CreateCourseRequest struct {
	TenantID     string  `json:"tenant_id"`
	InstructorID string  `json:"instructor_id"`
	CategoryID   string  `json:"category_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	LessonsCount int     `json:"lessons_count"`
	Duration     string  `json:"duration"`
	IsPublished  bool    `json:"is_published"`
}

// UpdateCourseRequest represents a request to update a course
type UpdateCourseRequest struct {
	Title        *string  `json:"title,omitempty"`
	Description  *string  `json:"description,omitempty"`
	CategoryID   *string  `json:"category_id,omitempty"`
	ThumbnailURL *string  `json:"thumbnail_url,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	LessonsCount *int     `json:"lessons_count,omitempty"`
	Duration     *string  `json:"duration,omitempty"`
	IsPublished  *bool    `json:"is_published,omitempty"`
	IsFeatured   *bool    `json:"is_featured,omitempty"`
}

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	GetByID(id string) (*Course, error)
	GetBySlug(tenantID, slug string) (*Course, error)
	Create(course *Course) error
	Update(course *Course) error
	Delete(id string) error
	ListByTenant(tenantID string, limit, offset int) ([]*Course, error)
	ListByInstructor(instructorID string, limit, offset int) ([]*Course, error)
	ListPublished(tenantID string, limit, offset int) ([]*Course, error)
}

// CourseService defines the interface for course business logic
type CourseService interface {
	GetCourse(id string) (*Course, error)
	GetCourseBySlug(tenantID, slug string) (*Course, error)
	CreateCourse(tenantID string, req *CreateCourseRequest) (*Course, error)
	UpdateCourse(id string, req *UpdateCourseRequest) (*Course, error)
	DeleteCourse(id string) error
	ListCourses(tenantID string, limit, offset int) ([]*Course, error)
	PublishCourse(id string) error
}
