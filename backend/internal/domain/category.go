package domain

import "time"

// Category represents a course category
type Category struct {
	ID          string    `json:"id" db:"id"`
	TenantID    *string   `json:"tenant_id,omitempty" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	CourseCount int       `json:"course_count" db:"course_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetByID(id string) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id string) error
	List(tenantID string, limit, offset int) ([]*Category, error)
	Count(tenantID string) (int, error)
}
