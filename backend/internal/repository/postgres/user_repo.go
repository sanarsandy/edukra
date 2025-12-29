package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// UserRepository implements domain.UserRepository using PostgreSQL
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, full_name, avatar_url, 
		       bio, phone, google_id, auth_provider, is_active, metadata, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	var user domain.User
	var tenantID, avatarURL, bio, phone, googleID, passwordHash sql.NullString
	
	var metadata []byte
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &tenantID, &user.Email, &passwordHash, &user.Role,
		&user.FullName, &avatarURL, &bio, &phone, &googleID, &user.AuthProvider,
		&user.IsActive, &metadata, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tenantID.Valid {
		user.TenantID = &tenantID.String
	}
	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if len(metadata) > 0 {
		json.Unmarshal(metadata, &user.Metadata)
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by tenant and email
func (r *UserRepository) GetByEmail(tenantID, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, full_name, avatar_url,
		       bio, phone, google_id, auth_provider, is_active, metadata, created_at, updated_at
		FROM users 
		WHERE email = $1 AND (tenant_id = $2 OR tenant_id IS NULL)
	`
	
	var user domain.User
	var tid, avatarURL, bio, phone, googleID sql.NullString
	
	var metadata []byte
	
	err := r.db.QueryRow(query, email, tenantID).Scan(
		&user.ID, &tid, &user.Email, &user.PasswordHash, &user.Role,
		&user.FullName, &avatarURL, &bio, &phone, &googleID, &user.AuthProvider,
		&user.IsActive, &metadata, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tid.Valid {
		user.TenantID = &tid.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if len(metadata) > 0 {
		json.Unmarshal(metadata, &user.Metadata)
	}
	
	return &user, nil
}

// GetByGoogleID retrieves a user by Google ID
func (r *UserRepository) GetByGoogleID(googleID string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, role, full_name, avatar_url,
		       google_id, auth_provider, is_active, metadata, created_at, updated_at
		FROM users WHERE google_id = $1
	`
	
	var user domain.User
	var tenantID, avatarURL, gID sql.NullString
	
	var metadata []byte
	
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID, &tenantID, &user.Email, &user.PasswordHash, &user.Role,
		&user.FullName, &avatarURL, &gID, &user.AuthProvider,
		&user.IsActive, &metadata, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if tenantID.Valid {
		user.TenantID = &tenantID.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if gID.Valid {
		user.GoogleID = &gID.String
	}
	if len(metadata) > 0 {
		json.Unmarshal(metadata, &user.Metadata)
	}
	
	return &user, nil
}

// Create inserts a new user
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (tenant_id, email, password_hash, role, full_name, avatar_url,
		                   bio, phone, google_id, auth_provider, is_active, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true
	
	metadata, _ := json.Marshal(user.Metadata)
	if user.Metadata == nil {
		metadata = []byte("{}")
	}

	return r.db.QueryRow(query,
		user.TenantID, user.Email, user.PasswordHash, user.Role, user.FullName,
		user.AvatarURL, user.Bio, user.Phone, user.GoogleID, user.AuthProvider, 
		user.IsActive, metadata, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users SET
			email = $2, role = $3, full_name = $4, avatar_url = $5,
			is_active = $6, metadata = $7, updated_at = $8, bio = $9, phone = $10
		WHERE id = $1
	`
	
	user.UpdatedAt = time.Now()
	
	metadata, _ := json.Marshal(user.Metadata)
	if user.Metadata == nil {
		metadata = []byte("{}")
	}

	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Role, user.FullName,
		user.AvatarURL, user.IsActive, metadata, user.UpdatedAt,
		user.Bio, user.Phone,
	)
	return err
}

// Delete removes a user by ID
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListByTenant retrieves users for a tenant with pagination
func (r *UserRepository) ListByTenant(tenantID string, limit, offset int) ([]*domain.User, error) {
	var query string
	var rows *sql.Rows
	var err error
	
	// If tenantID is empty or 'default', just get users with NULL tenant_id
	if tenantID == "" || tenantID == "default" {
		query = `
			SELECT id, tenant_id, email, password_hash, role, full_name, avatar_url,
			       google_id, auth_provider, is_active, metadata, created_at, updated_at
			FROM users 
			WHERE tenant_id IS NULL
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = r.db.Query(query, limit, offset)
	} else {
		query = `
			SELECT id, tenant_id, email, password_hash, role, full_name, avatar_url,
			       google_id, auth_provider, is_active, metadata, created_at, updated_at
			FROM users 
			WHERE tenant_id = $1 OR tenant_id IS NULL
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.db.Query(query, tenantID, limit, offset)
	}
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var tid, avatarURL, googleID, passwordHash sql.NullString
		
		var metadata []byte
		
		err := rows.Scan(
			&user.ID, &tid, &user.Email, &passwordHash, &user.Role,
			&user.FullName, &avatarURL, &googleID, &user.AuthProvider,
			&user.IsActive, &metadata, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if tid.Valid {
			user.TenantID = &tid.String
		}
		if passwordHash.Valid {
			user.PasswordHash = passwordHash.String
		}
		if avatarURL.Valid {
			user.AvatarURL = &avatarURL.String
		}
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if len(metadata) > 0 {
			json.Unmarshal(metadata, &user.Metadata)
		}
		
		users = append(users, &user)
	}
	
	return users, nil
}

// CountByTenant returns total user count for a tenant
func (r *UserRepository) CountByTenant(tenantID string) (int, error) {
	var query string
	var count int
	var err error
	
	if tenantID == "" || tenantID == "default" {
		query = `SELECT COUNT(*) FROM users WHERE tenant_id IS NULL`
		err = r.db.QueryRow(query).Scan(&count)
	} else {
		query = `SELECT COUNT(*) FROM users WHERE tenant_id = $1 OR tenant_id IS NULL`
		err = r.db.QueryRow(query, tenantID).Scan(&count)
	}
	return count, err
}

// UpdatePassword updates only the password hash for a user
func (r *UserRepository) UpdatePassword(userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, userID, passwordHash, time.Now())
	return err
}
// GetMonthlyGrowth returns total new users for last 6 months
func (r *UserRepository) GetMonthlyGrowth(tenantID string) ([]int, []string, error) {
	// Query to get count of new users grouped by month
	query := `
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', NOW()) - INTERVAL '5 months',
				date_trunc('month', NOW()),
				'1 month'::interval
			) as month
		)
		SELECT 
			m.month,
			COUNT(u.id) as total
		FROM months m
		LEFT JOIN users u ON 
			date_trunc('month', u.created_at) = m.month 
			AND ($1 = '' OR $1 = 'default' OR u.tenant_id::text = $1 OR u.tenant_id IS NULL)
		GROUP BY m.month
		ORDER BY m.month ASC
	`
	
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	
	var users []int
	var months []string
	
	for rows.Next() {
		var month time.Time
		var total int
		if err := rows.Scan(&month, &total); err != nil {
			return nil, nil, err
		}
		months = append(months, month.Format("Jan"))
		users = append(users, total)
	}
	
	return users, months, nil
}
