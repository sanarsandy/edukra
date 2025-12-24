package postgres

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// CategoryRepository handles category database operations
type CategoryRepository struct {
	db *sqlx.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// List retrieves all categories
func (r *CategoryRepository) List(tenantID string, limit, offset int) ([]*domain.Category, error) {
	var query string
	var args []interface{}
	
	if tenantID == "" || tenantID == "default" {
		query = `
			SELECT id, tenant_id, name, slug, description, icon, course_count, created_at, updated_at
			FROM categories 
			WHERE tenant_id IS NULL
			ORDER BY name ASC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	} else {
		query = `
			SELECT id, tenant_id, name, slug, description, icon, course_count, created_at, updated_at
			FROM categories 
			WHERE tenant_id = $1 OR tenant_id IS NULL
			ORDER BY name ASC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{tenantID, limit, offset}
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []*domain.Category
	for rows.Next() {
		var cat domain.Category
		var tid, desc, icon sql.NullString
		
		err := rows.Scan(
			&cat.ID, &tid, &cat.Name, &cat.Slug, &desc, &icon,
			&cat.CourseCount, &cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if tid.Valid {
			cat.TenantID = &tid.String
		}
		if desc.Valid {
			cat.Description = desc.String
		}
		if icon.Valid {
			cat.Icon = icon.String
		}
		
		categories = append(categories, &cat)
	}
	
	return categories, nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id string) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, name, slug, description, icon, course_count, created_at, updated_at
		FROM categories WHERE id = $1
	`
	
	var cat domain.Category
	var tid, desc, icon sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&cat.ID, &tid, &cat.Name, &cat.Slug, &desc, &icon,
		&cat.CourseCount, &cat.CreatedAt, &cat.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tid.Valid {
		cat.TenantID = &tid.String
	}
	if desc.Valid {
		cat.Description = desc.String
	}
	if icon.Valid {
		cat.Icon = icon.String
	}
	
	return &cat, nil
}

// Create inserts a new category
func (r *CategoryRepository) Create(cat *domain.Category) error {
	query := `
		INSERT INTO categories (tenant_id, name, slug, description, icon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	
	now := time.Now()
	cat.CreatedAt = now
	cat.UpdatedAt = now
	
	// Generate slug from name if not provided
	if cat.Slug == "" {
		cat.Slug = strings.ToLower(strings.ReplaceAll(cat.Name, " ", "-"))
	}
	
	return r.db.QueryRow(query,
		cat.TenantID, cat.Name, cat.Slug, cat.Description, cat.Icon,
		cat.CreatedAt, cat.UpdatedAt,
	).Scan(&cat.ID)
}

// Update updates an existing category
func (r *CategoryRepository) Update(cat *domain.Category) error {
	query := `
		UPDATE categories 
		SET name = $1, slug = $2, description = $3, icon = $4, updated_at = $5
		WHERE id = $6
	`
	
	cat.UpdatedAt = time.Now()
	
	if cat.Slug == "" {
		cat.Slug = strings.ToLower(strings.ReplaceAll(cat.Name, " ", "-"))
	}
	
	_, err := r.db.Exec(query, cat.Name, cat.Slug, cat.Description, cat.Icon, cat.UpdatedAt, cat.ID)
	return err
}

// Delete removes a category
func (r *CategoryRepository) Delete(id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Count returns total category count
func (r *CategoryRepository) Count(tenantID string) (int, error) {
	var query string
	var count int
	var err error
	
	if tenantID == "" || tenantID == "default" {
		query = `SELECT COUNT(*) FROM categories WHERE tenant_id IS NULL`
		err = r.db.QueryRow(query).Scan(&count)
	} else {
		query = `SELECT COUNT(*) FROM categories WHERE tenant_id = $1 OR tenant_id IS NULL`
		err = r.db.QueryRow(query, tenantID).Scan(&count)
	}
	return count, err
}
