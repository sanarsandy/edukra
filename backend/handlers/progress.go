package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var progressRepo *postgres.LessonProgressRepository
var activityRepo *postgres.ActivityLogRepository

func initProgressRepos() {
	if progressRepo == nil && db.DB != nil {
		progressRepo = postgres.NewLessonProgressRepository(db.DB)
	}
	if activityRepo == nil && db.DB != nil {
		activityRepo = postgres.NewActivityLogRepository(db.DB)
	}
	initLessonRepos()
	initCourseRepos()
	initEnrollmentRepos()
}

// GetLessonProgress returns progress for a specific lesson
func GetLessonProgress(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	lessonID := c.Param("lessonId")
	progress, err := progressRepo.GetByUserAndLesson(userID, lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch progress"})
	}

	if progress == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"lesson_id":    lessonID,
			"is_completed": false,
			"watch_time":   0,
		})
	}

	return c.JSON(http.StatusOK, progress)
}

// UpdateLessonProgress updates progress for a lesson
func UpdateLessonProgressHandler(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	lessonID := c.Param("lessonId")

	var req domain.UpdateLessonProgressRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	progress := &domain.LessonProgress{
		UserID:      userID,
		LessonID:    lessonID,
		IsCompleted: req.IsCompleted,
		WatchTime:   req.WatchTime,
	}

	if req.IsCompleted {
		// Mark as complete
		if err := progressRepo.MarkComplete(userID, lessonID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update progress"})
		}

		// Log activity
		lesson, _ := lessonRepo.GetByID(lessonID)
		if lesson != nil {
			course, _ := courseRepo.GetByID(lesson.CourseID)
			courseName := ""
			if course != nil {
				courseName = course.Title
			}
			activityRepo.LogLessonComplete(userID, lessonID, lesson.Title, courseName)
			activityRepo.UpdateStreak(userID)

			// Calculate and sync course progress to enrollment
			courseProgress, _ := progressRepo.GetCourseProgress(userID, lesson.CourseID)
			if courseProgress != nil {
				// Sync progress to enrollment table
				enrollmentRepo.UpdateProgressByUserAndCourse(userID, lesson.CourseID, courseProgress.ProgressPercentage)
				
				// Check if course is now complete
				if courseProgress.ProgressPercentage >= 100 {
					if course != nil {
						activityRepo.LogCourseComplete(userID, lesson.CourseID, course.Title)
						// Auto-generate certificate
						GenerateCertificateForUser(userID, lesson.CourseID)
					}
				}
			}
		}
	} else {
		if err := progressRepo.Upsert(progress); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update progress"})
		}
	}

	// Return updated progress
	updated, _ := progressRepo.GetByUserAndLesson(userID, lessonID)
	if updated == nil {
		updated = progress
	}

	return c.JSON(http.StatusOK, updated)
}

// GetCourseProgress returns aggregated progress for a course
func GetCourseProgressHandler(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")
	progress, err := progressRepo.GetCourseProgress(userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course progress"})
	}

	return c.JSON(http.StatusOK, progress)
}

// UpdateWatchTime increments watch time for a lesson
func UpdateWatchTime(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	lessonID := c.Param("lessonId")

	var req domain.UpdateWatchTimeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Seconds <= 0 || req.Seconds > 3600 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid watch time"})
	}

	if err := progressRepo.UpdateWatchTime(userID, lessonID, req.Seconds); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update watch time"})
	}

	// Update streak on activity
	activityRepo.UpdateStreak(userID)

	return c.JSON(http.StatusOK, map[string]string{"message": "Watch time updated"})
}

// BulkUpdateProgress syncs multiple lessons at once
func BulkUpdateProgress(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	courseID := c.Param("courseId")

	var req domain.BulkProgressRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if len(req.LessonIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No lesson IDs provided"})
	}

	// Mark all lessons as complete
	if err := progressRepo.BulkMarkComplete(userID, req.LessonIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update progress"})
	}

	// Update streak
	activityRepo.UpdateStreak(userID)

	// Return updated course progress
	progress, _ := progressRepo.GetCourseProgress(userID, courseID)
	return c.JSON(http.StatusOK, progress)
}

// GetRecentActivities returns recent activities for the current user
func GetRecentActivities(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	activities, err := activityRepo.ListByUser(userID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch activities"})
	}

	if activities == nil {
		activities = []*domain.ActivityLog{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"activities": activities,
		"count":      len(activities),
	})
}

// GetLearningStats returns learning statistics for the current user
func GetLearningStats(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	// Get activity-based stats
	stats, err := activityRepo.GetUserStats(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stats"})
	}

	// Get total watch time
	totalWatchTime, _ := progressRepo.GetTotalWatchTime(userID)
	stats.TotalWatchTime = totalWatchTime

	// Get current streak
	streak, _ := activityRepo.GetStreak(userID)
	stats.CurrentStreak = streak

	return c.JSON(http.StatusOK, stats)
}

// GetWeeklyActivities returns activity counts per day for the last 7 days
func GetWeeklyActivities(c echo.Context) error {
	initProgressRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	// Get activities from the last 7 days
	activities, err := activityRepo.ListByUser(userID, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch activities"})
	}

	// Aggregate by day (last 7 days)
	dayNames := []string{"Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"}
	result := make([]map[string]interface{}, 7)
	
	// Initialize all days with 0 count
	for i := 0; i < 7; i++ {
		result[i] = map[string]interface{}{
			"name":  dayNames[i],
			"count": 0,
		}
	}

	// Count activities per day based on current streak info
	// Simple implementation: show activity on days when there's any activity
	for _, activity := range activities {
		if activity.CreatedAt.IsZero() {
			continue
		}
		// Get day of week (0 = Sunday, so adjust)
		weekday := int(activity.CreatedAt.Weekday())
		dayIdx := weekday - 1
		if dayIdx < 0 {
			dayIdx = 6 // Sunday
		}
		if dayIdx >= 0 && dayIdx < 7 {
			result[dayIdx]["count"] = result[dayIdx]["count"].(int) + 1
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"weekly": result,
	})
}
