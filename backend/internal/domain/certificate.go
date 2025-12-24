package domain

import "time"

// Certificate represents a course completion certificate
type Certificate struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	CourseID          string    `json:"course_id" db:"course_id"`
	CertificateNumber string    `json:"certificate_number" db:"certificate_number"`
	IssuedAt          time.Time `json:"issued_at" db:"issued_at"`
	PdfURL            *string   `json:"pdf_url,omitempty" db:"pdf_url"`
	// Joined fields (not in DB)
	CourseName string `json:"course_name,omitempty" db:"course_name"`
	UserName   string `json:"user_name,omitempty" db:"user_name"`
}
