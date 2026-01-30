package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var blogRepo *postgres.BlogRepository

func initBlogRepo() {
	if blogRepo == nil {
		blogRepo = postgres.NewBlogRepository(db.DB)
	}
}

// ==================== Admin Endpoints ====================

// CreateBlogPost creates a new blog post (admin only)
func CreateBlogPost(c echo.Context) error {
	initBlogRepo()

	var req domain.CreateBlogPostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Title == "" || req.Slug == "" || req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title, slug, and content are required"})
	}

	// Get author ID from auth context
	claims := c.Get("user")
	var authorID *string
	if claims != nil {
		if userClaims, ok := claims.(map[string]interface{}); ok {
			if id, ok := userClaims["user_id"].(string); ok {
				authorID = &id
			}
		}
	}

	post := &domain.BlogPost{
		Slug:            req.Slug,
		Title:           req.Title,
		Excerpt:         req.Excerpt,
		Content:         req.Content,
		ThumbnailURL:    req.ThumbnailURL,
		AuthorID:        authorID,
		Status:          req.Status,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	if err := blogRepo.CreatePost(post, req.CategoryIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, post)
}

// UpdateBlogPost updates an existing blog post (admin only)
func UpdateBlogPost(c echo.Context) error {
	initBlogRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Post ID is required"})
	}

	var req domain.UpdateBlogPostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Excerpt != nil {
		updates["excerpt"] = *req.Excerpt
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.ThumbnailURL != nil {
		updates["thumbnail_url"] = *req.ThumbnailURL
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		// Set published_at when publishing
		if *req.Status == "published" {
			updates["published_at"] = time.Now()
		}
	}
	if req.MetaTitle != nil {
		updates["meta_title"] = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		updates["meta_description"] = *req.MetaDescription
	}

	if err := blogRepo.UpdatePost(id, updates, req.CategoryIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update post: " + err.Error()})
	}

	// Return updated post
	post, err := blogRepo.GetPostByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get updated post"})
	}

	return c.JSON(http.StatusOK, post)
}

// DeleteBlogPost deletes a blog post (admin only)
func DeleteBlogPost(c echo.Context) error {
	initBlogRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Post ID is required"})
	}

	if err := blogRepo.DeletePost(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete post"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

// ListBlogPostsAdmin lists all blog posts for admin (including drafts)
func ListBlogPostsAdmin(c echo.Context) error {
	initBlogRepo()

	status := c.QueryParam("status") // "all", "draft", "published"
	page := 1
	perPage := 20

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pp := c.QueryParam("per_page"); pp != "" {
		if parsed, err := parseInt(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	result, err := blogRepo.ListPosts(status, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list posts"})
	}

	return c.JSON(http.StatusOK, result)
}

// GetBlogPostAdmin gets a single blog post by ID (admin only)
func GetBlogPostAdmin(c echo.Context) error {
	initBlogRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Post ID is required"})
	}

	post, err := blogRepo.GetPostByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get post"})
	}
	if post == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	return c.JSON(http.StatusOK, post)
}

// ==================== Public Endpoints ====================

// ListBlogPostsPublic lists published blog posts (public)
func ListBlogPostsPublic(c echo.Context) error {
	initBlogRepo()

	page := 1
	perPage := 10

	if p := c.QueryParam("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pp := c.QueryParam("per_page"); pp != "" {
		if parsed, err := parseInt(pp); err == nil && parsed > 0 && parsed <= 50 {
			perPage = parsed
		}
	}

	result, err := blogRepo.ListPosts("published", page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list posts"})
	}

	return c.JSON(http.StatusOK, result)
}

// GetBlogPostBySlug gets a published blog post by slug (public)
func GetBlogPostBySlug(c echo.Context) error {
	initBlogRepo()

	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Slug is required"})
	}

	post, err := blogRepo.GetPostBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get post"})
	}
	if post == nil || post.Status != "published" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	// Increment view count
	_ = blogRepo.IncrementViewCount(post.ID)

	return c.JSON(http.StatusOK, post)
}

// ==================== Categories ====================

// ListBlogCategories lists all blog categories (public)
func ListBlogCategories(c echo.Context) error {
	initBlogRepo()

	categories, err := blogRepo.ListCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list categories"})
	}

	return c.JSON(http.StatusOK, categories)
}

// CreateBlogCategory creates a new category (admin only)
func CreateBlogCategory(c echo.Context) error {
	initBlogRepo()

	var req domain.CreateBlogCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Name == "" || req.Slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name and slug are required"})
	}

	category := &domain.BlogCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := blogRepo.CreateCategory(category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create category"})
	}

	return c.JSON(http.StatusCreated, category)
}

// DeleteBlogCategory deletes a category (admin only)
func DeleteBlogCategory(c echo.Context) error {
	initBlogRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category ID is required"})
	}

	if err := blogRepo.DeleteCategory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete category"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}

// Helper function
func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
