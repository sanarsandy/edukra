package domain

import (
	"time"
)

// ContentType represents the type of lesson content
type ContentType string

const (
	ContentVideo ContentType = "video"
	ContentPDF   ContentType = "pdf"
	ContentQuiz  ContentType = "quiz"
	ContentText  ContentType = "text"
)

// SecurityLevel represents the content protection level
type SecurityLevel string

const (
	SecurityPublic    SecurityLevel = "public"
	SecuritySignedURL SecurityLevel = "signed_url"
	SecurityAES128    SecurityLevel = "aes_128"
	SecurityFullDRM   SecurityLevel = "full_drm"
)

// Lesson represents a lesson/module within a course
type Lesson struct {
	ID            string        `json:"id"`
	CourseID      string        `json:"course_id"`
	ParentID      *string       `json:"parent_id,omitempty"`
	IsContainer   bool          `json:"is_container"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	OrderIndex    int           `json:"order_index"`
	ContentType   ContentType   `json:"content_type"`
	VideoURL      *string       `json:"video_url,omitempty"`
	Content       *string       `json:"content,omitempty"` // Rich text content for text type
	VideoDuration int           `json:"video_duration"`    // In seconds
	SecurityLevel SecurityLevel `json:"security_level"`
	IsPreview     bool          `json:"is_preview"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Children      []*Lesson     `json:"children,omitempty"` // Nested children for tree structure
}


// ContentKey represents encryption keys for protected content
type ContentKey struct {
	ID            string     `json:"id"`
	LessonID      string     `json:"lesson_id"`
	EncryptionKey string     `json:"-"` // Never expose in JSON
	IV            string     `json:"-"` // Never expose in JSON
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// CreateLessonRequest represents a request to create a lesson
type CreateLessonRequest struct {
	CourseID      string        `json:"course_id"`
	ParentID      *string       `json:"parent_id,omitempty"`
	IsContainer   bool          `json:"is_container"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	OrderIndex    int           `json:"order_index"`
	ContentType   ContentType   `json:"content_type"`
	VideoURL      *string       `json:"video_url,omitempty"`
	Content       *string       `json:"content,omitempty"`
	SecurityLevel SecurityLevel `json:"security_level"`
	IsPreview     bool          `json:"is_preview"`
}

// UpdateLessonRequest represents a request to update a lesson
type UpdateLessonRequest struct {
	ParentID      *string        `json:"parent_id,omitempty"`
	IsContainer   *bool          `json:"is_container,omitempty"`
	Title         *string        `json:"title,omitempty"`
	Description   *string        `json:"description,omitempty"`
	OrderIndex    *int           `json:"order_index,omitempty"`
	ContentType   *ContentType   `json:"content_type,omitempty"`
	VideoURL      *string        `json:"video_url,omitempty"`
	Content       *string        `json:"content,omitempty"`
	SecurityLevel *SecurityLevel `json:"security_level,omitempty"`
	IsPreview     *bool          `json:"is_preview,omitempty"`
}

// MoveLessonRequest represents a request to move a lesson to a new parent
type MoveLessonRequest struct {
	ParentID   *string `json:"parent_id"`
	OrderIndex int     `json:"order_index"`
}

// LessonRepository defines the interface for lesson data access
type LessonRepository interface {
	GetByID(id string) (*Lesson, error)
	Create(lesson *Lesson) error
	Update(lesson *Lesson) error
	Delete(id string) error
	ListByCourse(courseID string) ([]*Lesson, error)
	ListByParent(courseID string, parentID *string) ([]*Lesson, error)
	GetTree(courseID string) ([]*Lesson, error)
	ReorderLessons(courseID string, lessonIDs []string) error
	MoveLesson(lessonID string, newParentID *string, newOrderIndex int) error
	CountByParent(courseID string, parentID *string) (int, error)
}

// ContentKeyRepository defines the interface for content key data access
type ContentKeyRepository interface {
	GetByLessonID(lessonID string) (*ContentKey, error)
	Create(key *ContentKey) error
	Delete(lessonID string) error
}

// LessonService defines the interface for lesson business logic
type LessonService interface {
	GetLesson(id string, userID string) (*Lesson, error)
	CreateLesson(req *CreateLessonRequest) (*Lesson, error)
	UpdateLesson(id string, req *UpdateLessonRequest) (*Lesson, error)
	DeleteLesson(id string) error
	ListLessons(courseID string) ([]*Lesson, error)
	GetSecureVideoURL(lessonID, userID string) (string, error)
}

// ContentProtectionService defines the interface for DRM/content protection
type ContentProtectionService interface {
	EncryptContent(lessonID string, videoPath string) (*ContentKey, error)
	GetPlaybackKey(lessonID, userID string) (string, error)
	ValidateAccess(lessonID, userID string) (bool, error)
}
