package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var ratingRepo *postgres.RatingRepository
var enrollmentRepoRating *postgres.EnrollmentRepository

func initRatingRepos() {
	if ratingRepo == nil && db.DB != nil {
		ratingRepo = postgres.NewRatingRepository(db.DB)
	}
	if enrollmentRepoRating == nil && db.DB != nil {
		enrollmentRepoRating = postgres.NewEnrollmentRepository(db.DB)
	}
}

// GetCourseRatings returns ratings for a course
func GetCourseRatings(c echo.Context) error {
	initRatingRepos()

	courseID := c.Param("courseId")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	ratings, err := ratingRepo.GetByCourse(courseID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch ratings"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ratings": ratings,
		"count":   len(ratings),
	})
}

// GetCourseRatingStats returns rating statistics for a course
func GetCourseRatingStats(c echo.Context) error {
	initRatingRepos()

	courseID := c.Param("courseId")

	stats, err := ratingRepo.GetCourseStats(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch rating stats"})
	}

	return c.JSON(http.StatusOK, stats)
}

// GetMyRating returns the current user's rating for a course
func GetMyRating(c echo.Context) error {
	initRatingRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")

	rating, err := ratingRepo.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch rating"})
	}

	if rating == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"has_rated": false,
			"rating":    nil,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"has_rated": true,
		"rating":    rating,
	})
}

// CreateCourseRating creates a new rating (must be enrolled)
func CreateCourseRating(c echo.Context) error {
	initRatingRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")

	// Check if user is enrolled
	enrollment, _ := enrollmentRepoRating.GetByUserAndCourse(userID, courseID)
	if enrollment == nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You must be enrolled to rate this course"})
	}

	// Check if already rated
	exists, _ := ratingRepo.Exists(userID, courseID)
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "You have already rated this course"})
	}

	var req domain.CreateRatingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 1 and 5"})
	}

	rating := &domain.CourseRating{
		UserID:   userID,
		CourseID: courseID,
		Rating:   req.Rating,
		Review:   req.Review,
	}

	if err := ratingRepo.Create(rating); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create rating"})
	}

	return c.JSON(http.StatusCreated, rating)
}

// UpdateCourseRating updates an existing rating
func UpdateCourseRating(c echo.Context) error {
	initRatingRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")

	// Get existing rating
	rating, err := ratingRepo.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch rating"})
	}
	if rating == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Rating not found"})
	}

	var req domain.UpdateRatingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 1 and 5"})
	}

	rating.Rating = req.Rating
	rating.Review = req.Review

	if err := ratingRepo.Update(rating); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update rating"})
	}

	return c.JSON(http.StatusOK, rating)
}

// DeleteCourseRating deletes a rating (user can only delete their own)
func DeleteCourseRating(c echo.Context) error {
	initRatingRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")

	// Get existing rating
	rating, err := ratingRepo.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch rating"})
	}
	if rating == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Rating not found"})
	}

	if err := ratingRepo.Delete(rating.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete rating"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rating deleted"})
}

// ============== ADMIN HANDLERS ==============

// AdminGetAllRatings returns all ratings with pagination (admin only)
func AdminGetAllRatings(c echo.Context) error {
	initRatingRepos()

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	courseID := c.QueryParam("course_id")

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var ratings []*domain.CourseRating
	var err error

	if courseID != "" {
		ratings, err = ratingRepo.GetByCourse(courseID, limit)
	} else {
		ratings, err = ratingRepo.GetAll(limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch ratings"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ratings": ratings,
		"count":   len(ratings),
		"limit":   limit,
		"offset":  offset,
	})
}

// AdminGetCourseRatings returns all ratings for a specific course (admin only)
func AdminGetCourseRatings(c echo.Context) error {
	initRatingRepos()

	courseID := c.Param("courseId")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	ratings, err := ratingRepo.GetByCourse(courseID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch ratings"})
	}

	stats, _ := ratingRepo.GetCourseStats(courseID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ratings": ratings,
		"count":   len(ratings),
		"stats":   stats,
	})
}

// AdminDeleteRating deletes any rating by ID (admin only)
func AdminDeleteRating(c echo.Context) error {
	initRatingRepos()

	ratingID := c.Param("id")

	rating, err := ratingRepo.GetByID(ratingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch rating"})
	}
	if rating == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Rating not found"})
	}

	if err := ratingRepo.Delete(ratingID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete rating"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rating deleted"})
}

// AdminGetRatingStats returns overall rating statistics (admin only)
func AdminGetRatingStats(c echo.Context) error {
	initRatingRepos()

	totalCount, _ := ratingRepo.Count()
	avgRating, _ := ratingRepo.GetAverageAll()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_ratings":    totalCount,
		"average_rating":   avgRating,
	})
}

