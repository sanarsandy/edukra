package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var categoryRepo *postgres.CategoryRepository

func initCategoryRepos() {
	if categoryRepo == nil && db.DB != nil {
		categoryRepo = postgres.NewCategoryRepository(db.DB)
	}
}

// ListCategories returns all categories
func ListCategories(c echo.Context) error {
	initCategoryRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	categories, err := categoryRepo.List(tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch categories"})
	}
	
	total, _ := categoryRepo.Count(tenantID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
		"total":      total,
	})
}

// GetCategory returns a single category by ID
func GetCategory(c echo.Context) error {
	initCategoryRepos()
	
	id := c.Param("id")
	
	category, err := categoryRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch category"})
	}
	
	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	
	return c.JSON(http.StatusOK, category)
}

// CreateCategory creates a new category (admin only)
func CreateCategory(c echo.Context) error {
	initCategoryRepos()
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}
	
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}
	
	err := categoryRepo.Create(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create category: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates a category (admin only)
func UpdateCategory(c echo.Context) error {
	initCategoryRepos()
	
	id := c.Param("id")
	
	category, err := categoryRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch category"})
	}
	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	
	err = categoryRepo.Update(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update category"})
	}
	
	return c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category (admin only)
func DeleteCategory(c echo.Context) error {
	initCategoryRepos()
	
	id := c.Param("id")
	
	category, err := categoryRepo.GetByID(id)
	if err != nil || category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	
	err = categoryRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete category"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
