package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type BlogRepository struct {
	db *sqlx.DB
}

func NewBlogRepository(db *sqlx.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// ==================== Posts ====================

func (r *BlogRepository) CreatePost(post *domain.BlogPost, categoryIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO blog_posts (slug, title, excerpt, content, thumbnail_url, author_id, status, published_at, meta_title, meta_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	var publishedAt *time.Time
	if post.Status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	err = tx.QueryRowx(query,
		post.Slug, post.Title, post.Excerpt, post.Content, post.ThumbnailURL,
		post.AuthorID, post.Status, publishedAt, post.MetaTitle, post.MetaDescription,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	// Insert category relations
	if len(categoryIDs) > 0 {
		for _, catID := range categoryIDs {
			_, err = tx.Exec(`INSERT INTO blog_post_categories (post_id, category_id) VALUES ($1, $2)`, post.ID, catID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *BlogRepository) UpdatePost(id string, updates map[string]interface{}, categoryIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build dynamic update query
	if len(updates) > 0 {
		setParts := []string{}
		args := []interface{}{}
		i := 1

		for key, value := range updates {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, i))
			args = append(args, value)
			i++
		}

		args = append(args, id)
		query := fmt.Sprintf("UPDATE blog_posts SET %s WHERE id = $%d", strings.Join(setParts, ", "), i)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	// Update categories if provided (replace all)
	if categoryIDs != nil {
		// Delete existing
		_, err = tx.Exec(`DELETE FROM blog_post_categories WHERE post_id = $1`, id)
		if err != nil {
			return err
		}

		// Insert new
		for _, catID := range categoryIDs {
			_, err = tx.Exec(`INSERT INTO blog_post_categories (post_id, category_id) VALUES ($1, $2)`, id, catID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *BlogRepository) DeletePost(id string) error {
	_, err := r.db.Exec(`DELETE FROM blog_posts WHERE id = $1`, id)
	return err
}

func (r *BlogRepository) GetPostByID(id string) (*domain.BlogPost, error) {
	post := &domain.BlogPost{}
	err := r.db.Get(post, `SELECT * FROM blog_posts WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Load author
	if post.AuthorID != nil {
		author := &domain.User{}
		err = r.db.Get(author, `SELECT id, email, full_name, avatar FROM users WHERE id = $1`, *post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	// Load categories
	categories, _ := r.GetPostCategories(post.ID)
	post.Categories = categories

	return post, nil
}

func (r *BlogRepository) GetPostBySlug(slug string) (*domain.BlogPost, error) {
	post := &domain.BlogPost{}
	err := r.db.Get(post, `SELECT * FROM blog_posts WHERE slug = $1`, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Load author
	if post.AuthorID != nil {
		author := &domain.User{}
		err = r.db.Get(author, `SELECT id, email, full_name, avatar FROM users WHERE id = $1`, *post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	// Load categories
	categories, _ := r.GetPostCategories(post.ID)
	post.Categories = categories

	return post, nil
}

func (r *BlogRepository) ListPosts(status string, page, perPage int) (*domain.BlogPostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// Build query based on status filter
	var whereClause string
	var args []interface{}
	argIdx := 1

	if status != "" && status != "all" {
		whereClause = fmt.Sprintf("WHERE status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM blog_posts %s", whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, err
	}

	// Fetch posts
	query := fmt.Sprintf(`
		SELECT * FROM blog_posts %s 
		ORDER BY created_at DESC 
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)
	args = append(args, perPage, offset)

	posts := []domain.BlogPost{}
	err = r.db.Select(&posts, query, args...)
	if err != nil {
		return nil, err
	}

	// Load authors for each post
	for i := range posts {
		if posts[i].AuthorID != nil {
			author := &domain.User{}
			err = r.db.Get(author, `SELECT id, email, full_name, avatar FROM users WHERE id = $1`, *posts[i].AuthorID)
			if err == nil {
				posts[i].Author = author
			}
		}
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &domain.BlogPostListResponse{
		Posts:      posts,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (r *BlogRepository) IncrementViewCount(id string) error {
	_, err := r.db.Exec(`UPDATE blog_posts SET view_count = view_count + 1 WHERE id = $1`, id)
	return err
}

// ==================== Categories ====================

func (r *BlogRepository) CreateCategory(category *domain.BlogCategory) error {
	query := `
		INSERT INTO blog_categories (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowx(query, category.Name, category.Slug, category.Description).
		Scan(&category.ID, &category.CreatedAt)
}

func (r *BlogRepository) UpdateCategory(id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setParts := []string{}
	args := []interface{}{}
	i := 1

	for key, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE blog_categories SET %s WHERE id = $%d", strings.Join(setParts, ", "), i)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *BlogRepository) DeleteCategory(id string) error {
	_, err := r.db.Exec(`DELETE FROM blog_categories WHERE id = $1`, id)
	return err
}

func (r *BlogRepository) GetCategoryByID(id string) (*domain.BlogCategory, error) {
	category := &domain.BlogCategory{}
	err := r.db.Get(category, `SELECT * FROM blog_categories WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (r *BlogRepository) ListCategories() ([]domain.BlogCategory, error) {
	categories := []domain.BlogCategory{}
	err := r.db.Select(&categories, `SELECT * FROM blog_categories ORDER BY name ASC`)
	return categories, err
}

func (r *BlogRepository) GetPostCategories(postID string) ([]domain.BlogCategory, error) {
	categories := []domain.BlogCategory{}
	query := `
		SELECT c.* FROM blog_categories c
		INNER JOIN blog_post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = $1
		ORDER BY c.name ASC
	`
	err := r.db.Select(&categories, query, postID)
	return categories, err
}
