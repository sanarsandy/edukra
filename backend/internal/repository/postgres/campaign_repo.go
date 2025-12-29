package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

type CampaignRepository struct {
	db *sqlx.DB
}

func NewCampaignRepository(db *sqlx.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

// Create creates a new campaign
func (r *CampaignRepository) Create(c *domain.Campaign) error {
	query := `
		INSERT INTO campaigns (
			slug, is_active, course_id, title, meta_description, og_image_url,
			blocks, styles, html_content, css_content, gjs_data,
			start_date, end_date,
			view_count, click_count, conversion_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13,
			$14, $15, $16, $17, $18
		) RETURNING id
	`

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	// Ensure JSON fields have defaults
	if c.Blocks == nil {
		defaultBlocks := domain.GetDefaultBlocks()
		c.Blocks, _ = json.Marshal(defaultBlocks)
	}
	if c.Styles == nil {
		c.Styles = []byte("{}")
	}
	if c.GJSData == nil {
		c.GJSData = []byte("{}")
	}

	return r.db.QueryRow(query,
		c.Slug, c.IsActive, c.CourseID, c.Title, c.MetaDesc, c.OGImageURL,
		c.Blocks, c.Styles, c.HTMLContent, c.CSSContent, c.GJSData,
		c.StartDate, c.EndDate,
		c.ViewCount, c.ClickCount, c.ConversionCount, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
}

// GetByID retrieves a campaign by ID (admin)
func (r *CampaignRepository) GetByID(id string) (*domain.Campaign, error) {
	query := `
		SELECT 
			c.id, c.slug, c.is_active, c.course_id, c.title, c.meta_description, c.og_image_url,
			c.blocks, c.styles, c.html_content, c.css_content, c.gjs_data,
			c.start_date, c.end_date,
			c.view_count, c.click_count, c.conversion_count, c.created_at, c.updated_at,
			course.id, course.title, course.slug, course.price, course.discount_price, course.thumbnail_url
		FROM campaigns c
		LEFT JOIN courses course ON c.course_id = course.id
		WHERE c.id = $1
	`

	var camp domain.Campaign
	var courseID, courseTitle, courseSlug, courseThumbnail sql.NullString
	var coursePrice, courseDiscountPrice sql.NullFloat64
	// Handle nullable GrapeJS fields
	var htmlContent, cssContent, gjsData sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&camp.ID, &camp.Slug, &camp.IsActive, &camp.CourseID, &camp.Title, &camp.MetaDesc, &camp.OGImageURL,
		&camp.Blocks, &camp.Styles, &htmlContent, &cssContent, &gjsData,
		&camp.StartDate, &camp.EndDate,
		&camp.ViewCount, &camp.ClickCount, &camp.ConversionCount, &camp.CreatedAt, &camp.UpdatedAt,
		&courseID, &courseTitle, &courseSlug, &coursePrice, &courseDiscountPrice, &courseThumbnail,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Assign nullable GrapeJS fields
	if htmlContent.Valid {
		camp.HTMLContent = &htmlContent.String
	}
	if cssContent.Valid {
		camp.CSSContent = &cssContent.String
	}
	if gjsData.Valid {
		camp.GJSData = json.RawMessage(gjsData.String)
	}

	if courseID.Valid {
		camp.Course = &domain.Course{
			ID:           courseID.String,
			Title:        courseTitle.String,
			Slug:         courseSlug.String,
			Price:        coursePrice.Float64,
			ThumbnailURL: &courseThumbnail.String,
		}
		if courseDiscountPrice.Valid {
			camp.Course.DiscountPrice = &courseDiscountPrice.Float64
		}
	}

	return &camp, nil
}


// GetBySlug retrieves an active campaign by slug (public)
func (r *CampaignRepository) GetBySlug(slug string) (*domain.Campaign, error) {
	query := `
		SELECT 
			c.id, c.slug, c.is_active, c.course_id, c.title, c.meta_description, c.og_image_url,
			c.blocks, c.styles, c.html_content, c.css_content, c.gjs_data,
			c.start_date, c.end_date,
			c.view_count, c.click_count, c.conversion_count, c.created_at, c.updated_at,
			course.id, course.title, course.slug, course.price, course.discount_price, course.thumbnail_url, course.description,
			u.id, u.full_name, u.avatar_url, u.bio
		FROM campaigns c
		LEFT JOIN courses course ON c.course_id = course.id
		LEFT JOIN users u ON course.instructor_id = u.id
		WHERE c.slug = $1 AND c.is_active = true
	`

	var camp domain.Campaign
	var courseID, courseTitle, courseSlug, courseThumbnail, courseDesc sql.NullString
	var coursePrice, courseDiscountPrice sql.NullFloat64
	var instrID, instrName, instrAvatar, instrBio sql.NullString
	// Handle nullable GrapeJS fields
	var htmlContent, cssContent, gjsData sql.NullString

	err := r.db.QueryRow(query, slug).Scan(
		&camp.ID, &camp.Slug, &camp.IsActive, &camp.CourseID, &camp.Title, &camp.MetaDesc, &camp.OGImageURL,
		&camp.Blocks, &camp.Styles, &htmlContent, &cssContent, &gjsData,
		&camp.StartDate, &camp.EndDate,
		&camp.ViewCount, &camp.ClickCount, &camp.ConversionCount, &camp.CreatedAt, &camp.UpdatedAt,
		&courseID, &courseTitle, &courseSlug, &coursePrice, &courseDiscountPrice, &courseThumbnail, &courseDesc,
		&instrID, &instrName, &instrAvatar, &instrBio,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Assign nullable GrapeJS fields
	if htmlContent.Valid {
		camp.HTMLContent = &htmlContent.String
	}
	if cssContent.Valid {
		camp.CSSContent = &cssContent.String
	}
	if gjsData.Valid {
		camp.GJSData = json.RawMessage(gjsData.String)
	}

	if courseID.Valid {
		camp.Course = &domain.Course{
			ID:           courseID.String,
			Title:        courseTitle.String,
			Slug:         courseSlug.String,
			Price:        coursePrice.Float64,
			ThumbnailURL: &courseThumbnail.String,
			Description:  courseDesc.String,
		}
		if courseDiscountPrice.Valid {
			camp.Course.DiscountPrice = &courseDiscountPrice.Float64
		}
		if instrID.Valid {
			camp.Course.Instructor = &domain.User{
				ID:        instrID.String,
				FullName:  instrName.String,
				AvatarURL: &instrAvatar.String,
				Bio:       &instrBio.String,
			}
		}
	}

	return &camp, nil
}

// List retrieves all campaigns (admin)
func (r *CampaignRepository) List(limit, offset int) ([]*domain.Campaign, error) {
	query := `
		SELECT 
			c.id, c.slug, c.is_active, c.course_id, c.title, c.meta_description, c.og_image_url,
			c.blocks, c.styles, c.start_date, c.end_date,
			c.view_count, c.click_count, c.conversion_count, c.created_at, c.updated_at,
			course.id, course.title, course.thumbnail_url
		FROM campaigns c
		LEFT JOIN courses course ON c.course_id = course.id
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.Campaign
	for rows.Next() {
		var camp domain.Campaign
		var courseID, courseTitle, courseThumbnail sql.NullString

		err := rows.Scan(
			&camp.ID, &camp.Slug, &camp.IsActive, &camp.CourseID, &camp.Title, &camp.MetaDesc, &camp.OGImageURL,
			&camp.Blocks, &camp.Styles, &camp.StartDate, &camp.EndDate,
			&camp.ViewCount, &camp.ClickCount, &camp.ConversionCount, &camp.CreatedAt, &camp.UpdatedAt,
			&courseID, &courseTitle, &courseThumbnail,
		)
		if err != nil {
			return nil, err
		}

		if courseID.Valid {
			camp.Course = &domain.Course{
				ID:           courseID.String,
				Title:        courseTitle.String,
				ThumbnailURL: &courseThumbnail.String,
			}
		}

		campaigns = append(campaigns, &camp)
	}

	return campaigns, nil
}

// Update updates a campaign
func (r *CampaignRepository) Update(c *domain.Campaign) error {
	query := `
		UPDATE campaigns SET
			slug = $2, is_active = $3, course_id = $4, title = $5, 
			meta_description = $6, og_image_url = $7,
			blocks = $8, styles = $9, html_content = $10, css_content = $11, gjs_data = $12,
			start_date = $13, end_date = $14,
			updated_at = $15
		WHERE id = $1
	`

	c.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		c.ID, c.Slug, c.IsActive, c.CourseID, c.Title,
		c.MetaDesc, c.OGImageURL,
		c.Blocks, c.Styles, c.HTMLContent, c.CSSContent, c.GJSData,
		c.StartDate, c.EndDate,
		c.UpdatedAt,
	)

	return err
}

// Delete deletes a campaign
func (r *CampaignRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM campaigns WHERE id = $1`, id)
	return err
}

// IncrementViewCount increments the view counter
func (r *CampaignRepository) IncrementViewCount(id string) error {
	_, err := r.db.Exec(`UPDATE campaigns SET view_count = view_count + 1 WHERE id = $1`, id)
	return err
}

// IncrementClickCount increments the click counter
func (r *CampaignRepository) IncrementClickCount(id string) error {
	_, err := r.db.Exec(`UPDATE campaigns SET click_count = click_count + 1 WHERE id = $1`, id)
	return err
}

// IncrementConversionCount increments the conversion counter
func (r *CampaignRepository) IncrementConversionCount(id string) error {
	_, err := r.db.Exec(`UPDATE campaigns SET conversion_count = conversion_count + 1 WHERE id = $1`, id)
	return err
}

// TrackEvent records a detailed analytics event
func (r *CampaignRepository) TrackEvent(event *domain.CampaignAnalytics) error {
	query := `
		INSERT INTO campaign_analytics (
			campaign_id, event_type, transaction_id,
			visitor_ip, user_agent, referer,
			utm_source, utm_medium, utm_campaign, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id
	`

	event.CreatedAt = time.Now()

	return r.db.QueryRow(query,
		event.CampaignID, event.EventType, event.TransactionID,
		event.VisitorIP, event.UserAgent, event.Referer,
		event.UTMSource, event.UTMMedium, event.UTMCampaign, event.CreatedAt,
	).Scan(&event.ID)
}

// GetAnalyticsSummary returns analytics summary for a campaign
func (r *CampaignRepository) GetAnalyticsSummary(campaignID string) (map[string]int, error) {
	query := `
		SELECT event_type, COUNT(*) as count
		FROM campaign_analytics
		WHERE campaign_id = $1
		GROUP BY event_type
	`

	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{"view": 0, "click": 0, "conversion": 0}
	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, err
		}
		result[eventType] = count
	}

	return result, nil
}

// SlugExists checks if a slug is already taken
func (r *CampaignRepository) SlugExists(slug string, excludeID string) (bool, error) {
	query := `SELECT COUNT(*) FROM campaigns WHERE slug = $1 AND id != $2`
	var count int
	err := r.db.QueryRow(query, slug, excludeID).Scan(&count)
	return count > 0, err
}
