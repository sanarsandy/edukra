package domain

import (
	"time"
)

// BlogPost represents a blog article
type BlogPost struct {
	ID              string     `json:"id" db:"id"`
	Slug            string     `json:"slug" db:"slug"`
	Title           string     `json:"title" db:"title"`
	Excerpt         *string    `json:"excerpt,omitempty" db:"excerpt"`
	Content         string     `json:"content" db:"content"`
	ThumbnailURL    *string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	AuthorID        *string    `json:"author_id,omitempty" db:"author_id"`
	Status          string     `json:"status" db:"status"`
	PublishedAt     *time.Time `json:"published_at,omitempty" db:"published_at"`
	MetaTitle       *string    `json:"meta_title,omitempty" db:"meta_title"`
	MetaDescription *string    `json:"meta_description,omitempty" db:"meta_description"`
	ViewCount       int        `json:"view_count" db:"view_count"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Relations (not stored in blog_posts table)
	Author     *User          `json:"author,omitempty" db:"-"`
	Categories []BlogCategory `json:"categories,omitempty" db:"-"`
}

// BlogCategory represents a category for blog posts
type BlogCategory struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Request/Response structs

type CreateBlogPostRequest struct {
	Title           string   `json:"title" validate:"required,min=3,max=500"`
	Slug            string   `json:"slug" validate:"required,min=3,max=255"`
	Excerpt         *string  `json:"excerpt,omitempty"`
	Content         string   `json:"content" validate:"required"`
	ThumbnailURL    *string  `json:"thumbnail_url,omitempty"`
	Status          string   `json:"status" validate:"required,oneof=draft published"`
	MetaTitle       *string  `json:"meta_title,omitempty"`
	MetaDescription *string  `json:"meta_description,omitempty"`
	CategoryIDs     []string `json:"category_ids,omitempty"`
}

type UpdateBlogPostRequest struct {
	Title           *string  `json:"title,omitempty"`
	Slug            *string  `json:"slug,omitempty"`
	Excerpt         *string  `json:"excerpt,omitempty"`
	Content         *string  `json:"content,omitempty"`
	ThumbnailURL    *string  `json:"thumbnail_url,omitempty"`
	Status          *string  `json:"status,omitempty"`
	MetaTitle       *string  `json:"meta_title,omitempty"`
	MetaDescription *string  `json:"meta_description,omitempty"`
	CategoryIDs     []string `json:"category_ids,omitempty"`
}

type BlogPostListResponse struct {
	Posts      []BlogPost `json:"posts"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}

type CreateBlogCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Slug        string  `json:"slug" validate:"required,min=2,max=100"`
	Description *string `json:"description,omitempty"`
}

// BlogRepository interface for database operations
type BlogRepository interface {
	// Posts
	CreatePost(post *BlogPost, categoryIDs []string) error
	UpdatePost(id string, updates map[string]interface{}, categoryIDs []string) error
	DeletePost(id string) error
	GetPostByID(id string) (*BlogPost, error)
	GetPostBySlug(slug string) (*BlogPost, error)
	ListPosts(status string, page, perPage int) (*BlogPostListResponse, error)
	IncrementViewCount(id string) error

	// Categories
	CreateCategory(category *BlogCategory) error
	UpdateCategory(id string, updates map[string]interface{}) error
	DeleteCategory(id string) error
	GetCategoryByID(id string) (*BlogCategory, error)
	ListCategories() ([]BlogCategory, error)
	GetPostCategories(postID string) ([]BlogCategory, error)
}
