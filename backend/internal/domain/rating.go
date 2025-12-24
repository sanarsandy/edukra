package domain

import "time"

// CourseRating represents a user's rating and review for a course
type CourseRating struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Rating    int       `json:"rating" db:"rating"` // 1-5
	Review    *string   `json:"review,omitempty" db:"review"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// Joined fields
	UserName string `json:"user_name,omitempty" db:"user_name"`
}

// CourseRatingStats holds aggregated rating statistics
type CourseRatingStats struct {
	AverageRating float64 `json:"average_rating" db:"average_rating"`
	TotalRatings  int     `json:"total_ratings" db:"total_ratings"`
}

// CreateRatingRequest is the request body for creating a rating
type CreateRatingRequest struct {
	Rating int     `json:"rating"` // 1-5
	Review *string `json:"review,omitempty"`
}

// UpdateRatingRequest is the request body for updating a rating
type UpdateRatingRequest struct {
	Rating int     `json:"rating"` // 1-5
	Review *string `json:"review,omitempty"`
}
