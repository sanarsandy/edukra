package postgres

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// CertificateRepository handles certificate data access
type CertificateRepository struct {
	db *sqlx.DB
}

// NewCertificateRepository creates a new certificate repository
func NewCertificateRepository(db *sqlx.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

// GetByUser retrieves all certificates for a user
func (r *CertificateRepository) GetByUser(userID string) ([]*domain.Certificate, error) {
	certificates := []*domain.Certificate{}
	err := r.db.Select(&certificates, `
		SELECT 
			c.id, c.user_id, c.course_id, c.certificate_number, c.issued_at, c.pdf_url,
			co.title as course_name,
			u.full_name as user_name
		FROM certificates c
		JOIN courses co ON co.id = c.course_id
		JOIN users u ON u.id = c.user_id
		WHERE c.user_id = $1
		ORDER BY c.issued_at DESC
	`, userID)
	
	if err == sql.ErrNoRows {
		return certificates, nil
	}
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

// GetByID retrieves a certificate by ID
func (r *CertificateRepository) GetByID(id string) (*domain.Certificate, error) {
	var cert domain.Certificate
	err := r.db.Get(&cert, `
		SELECT 
			c.id, c.user_id, c.course_id, c.certificate_number, c.issued_at, c.pdf_url,
			co.title as course_name,
			u.full_name as user_name
		FROM certificates c
		JOIN courses co ON co.id = c.course_id
		JOIN users u ON u.id = c.user_id
		WHERE c.id = $1
	`, id)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// Create creates a new certificate
func (r *CertificateRepository) Create(cert *domain.Certificate) error {
	query := `
		INSERT INTO certificates (user_id, course_id, certificate_number)
		VALUES ($1, $2, $3)
		RETURNING id, issued_at
	`
	return r.db.QueryRow(query,
		cert.UserID,
		cert.CourseID,
		cert.CertificateNumber,
	).Scan(&cert.ID, &cert.IssuedAt)
}

// ExistsForCourse checks if a certificate already exists for user and course
func (r *CertificateRepository) ExistsForCourse(userID, courseID string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM certificates 
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID)
	return exists, err
}

// GetByCourseAndUser retrieves a certificate by user and course
func (r *CertificateRepository) GetByCourseAndUser(userID, courseID string) (*domain.Certificate, error) {
	var cert domain.Certificate
	err := r.db.Get(&cert, `
		SELECT 
			c.id, c.user_id, c.course_id, c.certificate_number, c.issued_at, c.pdf_url,
			co.title as course_name,
			u.full_name as user_name
		FROM certificates c
		JOIN courses co ON co.id = c.course_id
		JOIN users u ON u.id = c.user_id
		WHERE c.user_id = $1 AND c.course_id = $2
	`, userID, courseID)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// Count returns total certificates for a user
func (r *CertificateRepository) Count(userID string) (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM certificates WHERE user_id = $1`, userID)
	return count, err
}

// GenerateCertificateNumber generates a unique certificate number
func (r *CertificateRepository) GenerateCertificateNumber() string {
	year := time.Now().Year()
	bytes := make([]byte, 4)
	rand.Read(bytes)
	randomPart := hex.EncodeToString(bytes)
	return fmt.Sprintf("CERT-%d-%s", year, randomPart)
}
