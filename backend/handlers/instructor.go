package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// Course status constants
const (
	CourseStatusDraft         = "draft"
	CourseStatusPendingReview = "pending_review"
	CourseStatusApproved      = "approved"
	CourseStatusRejected      = "rejected"
	CourseStatusPublished     = "published"
)

// ========================================
// INSTRUCTOR DASHBOARD
// ========================================

// InstructorDashboard returns instructor's dashboard stats
func InstructorDashboard(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	stats := struct {
		TotalCourses     int     `json:"total_courses"`
		PublishedCourses int     `json:"published_courses"`
		DraftCourses     int     `json:"draft_courses"`
		PendingReview    int     `json:"pending_review"`
		TotalStudents    int     `json:"total_students"`
		AverageRating    float64 `json:"average_rating"`
		TotalRevenue     float64 `json:"total_revenue"`
	}{}

	// Total courses by this instructor
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1`, userID).Scan(&stats.TotalCourses)

	// Courses by status
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1 AND status = 'published'`, userID).Scan(&stats.PublishedCourses)
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1 AND status = 'draft'`, userID).Scan(&stats.DraftCourses)
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1 AND status = 'pending_review'`, userID).Scan(&stats.PendingReview)

	// Total enrolled students across all courses
	db.DB.QueryRow(`
		SELECT COUNT(DISTINCT e.user_id) 
		FROM enrollments e 
		JOIN courses c ON e.course_id = c.id 
		WHERE c.instructor_id = $1
	`, userID).Scan(&stats.TotalStudents)

	// Average rating across all courses
	db.DB.QueryRow(`
		SELECT COALESCE(AVG(r.rating), 0) 
		FROM course_ratings r 
		JOIN courses c ON r.course_id = c.id 
		WHERE c.instructor_id = $1
	`, userID).Scan(&stats.AverageRating)

	// Get recent courses
	rows, err := db.DB.Query(`
		SELECT id, title, status, COALESCE(lessons_count, 0), created_at
		FROM courses 
		WHERE instructor_id = $1 
		ORDER BY created_at DESC 
		LIMIT 5
	`, userID)
	
	var recentCourses []map[string]interface{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, title, status string
			var lessonsCount int
			var createdAt time.Time
			rows.Scan(&id, &title, &status, &lessonsCount, &createdAt)
			recentCourses = append(recentCourses, map[string]interface{}{
				"id":            id,
				"title":         title,
				"status":        status,
				"lessons_count": lessonsCount,
				"created_at":    createdAt,
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stats":          stats,
		"recent_courses": recentCourses,
	})
}

// ========================================
// INSTRUCTOR COURSES
// ========================================

// InstructorListCourses returns courses owned by the instructor
func InstructorListCourses(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	status := c.QueryParam("status") // filter by status
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var query string
	var args []interface{}

	if status != "" && status != "all" {
		query = `
			SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.price, 
			       c.status, COALESCE(c.lessons_count, 0), c.is_published, c.created_at, c.updated_at,
				   cat.name as category_name
			FROM courses c
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.instructor_id = $1 AND c.status = $2
			ORDER BY c.updated_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{userID, status, limit, offset}
	} else {
		query = `
			SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.price, 
			       c.status, COALESCE(c.lessons_count, 0), c.is_published, c.created_at, c.updated_at,
				   cat.name as category_name
			FROM courses c
			LEFT JOIN categories cat ON c.category_id = cat.id
			WHERE c.instructor_id = $1
			ORDER BY c.updated_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{userID, limit, offset}
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data kursus"})
	}
	defer rows.Close()

	var courses []map[string]interface{}
	for rows.Next() {
		var id, title, slug, description, status string
		var thumbnailURL, categoryName sql.NullString
		var price float64
		var lessonsCount int
		var isPublished bool
		var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &title, &slug, &description, &thumbnailURL, &price,
			&status, &lessonsCount, &isPublished, &createdAt, &updatedAt, &categoryName)
		if err != nil {
			continue
		}

		course := map[string]interface{}{
			"id":            id,
			"title":         title,
			"slug":          slug,
			"description":   description,
			"price":         price,
			"status":        status,
			"lessons_count": lessonsCount,
			"is_published":  isPublished,
			"created_at":    createdAt,
			"updated_at":    updatedAt,
		}
		if thumbnailURL.Valid {
			course["thumbnail_url"] = thumbnailURL.String
		}
		if categoryName.Valid {
			course["category_name"] = categoryName.String
		}

		// Get enrolled count for this course
		var enrolledCount int
		db.DB.QueryRow(`SELECT COUNT(*) FROM enrollments WHERE course_id = $1`, id).Scan(&enrolledCount)
		course["enrolled_count"] = enrolledCount

		courses = append(courses, course)
	}

	// Get total count
	var total int
	if status != "" && status != "all" {
		db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1 AND status = $2`, userID, status).Scan(&total)
	} else {
		db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE instructor_id = $1`, userID).Scan(&total)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"courses": courses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// InstructorCreateCourse creates a new course (always draft status)
func InstructorCreateCourse(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req domain.CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Judul kursus wajib diisi"})
	}

	// Generate slug
	slug := generateSlugFromTitle(req.Title)

	// Insert course with draft status (instructor cannot publish directly)
	query := `
		INSERT INTO courses (instructor_id, category_id, title, slug, description, thumbnail_url, 
		                     price, currency, status, is_published, lessons_count, duration, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'draft', false, 0, '', NOW(), NOW())
		RETURNING id, created_at
	`

	var categoryID interface{} = nil
	if req.CategoryID != "" {
		categoryID = req.CategoryID
	}

	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	var courseID string
	var createdAt time.Time
	err = db.DB.QueryRow(query,
		userID, categoryID, req.Title, slug, req.Description, req.ThumbnailURL,
		req.Price, currency,
	).Scan(&courseID, &createdAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat kursus: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         courseID,
		"title":      req.Title,
		"slug":       slug,
		"status":     CourseStatusDraft,
		"created_at": createdAt,
		"message":    "Kursus berhasil dibuat dalam status draft",
	})
}

// InstructorGetCourse gets a single course (must be owned by instructor)
func InstructorGetCourse(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	var course struct {
		ID           string    `json:"id"`
		Title        string    `json:"title"`
		Slug         string    `json:"slug"`
		Description  string    `json:"description"`
		ThumbnailURL *string   `json:"thumbnail_url"`
		Price        float64   `json:"price"`
		Currency     string    `json:"currency"`
		Status       string    `json:"status"`
		IsPublished  bool      `json:"is_published"`
		LessonsCount int       `json:"lessons_count"`
		Duration     string    `json:"duration"`
		CategoryID   *string   `json:"category_id"`
		CategoryName *string   `json:"category_name"`
		ReviewNotes  *string   `json:"review_notes"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	var thumbnailURL, categoryID, categoryName, reviewNotes sql.NullString

	err = db.DB.QueryRow(`
		SELECT c.id, c.title, c.slug, COALESCE(c.description, ''), c.thumbnail_url, c.price, c.currency,
		       COALESCE(c.status, 'draft'), c.is_published, COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''),
			   c.category_id, cat.name, c.review_notes, c.created_at, c.updated_at
		FROM courses c
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.id = $1 AND c.instructor_id = $2
	`, courseID, userID).Scan(
		&course.ID, &course.Title, &course.Slug, &course.Description, &thumbnailURL,
		&course.Price, &course.Currency, &course.Status, &course.IsPublished,
		&course.LessonsCount, &course.Duration, &categoryID, &categoryName, &reviewNotes,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data kursus"})
	}

	if thumbnailURL.Valid {
		course.ThumbnailURL = &thumbnailURL.String
	}
	if categoryID.Valid {
		course.CategoryID = &categoryID.String
	}
	if categoryName.Valid {
		course.CategoryName = &categoryName.String
	}
	if reviewNotes.Valid {
		course.ReviewNotes = &reviewNotes.String
	}

	return c.JSON(http.StatusOK, course)
}

// InstructorUpdateCourse updates a course (only if owned and not published)
func InstructorUpdateCourse(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Check ownership and status
	var currentStatus string
	err = db.DB.QueryRow(`
		SELECT status FROM courses WHERE id = $1 AND instructor_id = $2
	`, courseID, userID).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	// Cannot edit published courses
	if currentStatus == CourseStatusPublished {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Kursus yang sudah dipublish tidak dapat diedit. Hubungi admin untuk unpublish terlebih dahulu."})
	}

	var req struct {
		Title        *string  `json:"title"`
		Description  *string  `json:"description"`
		ThumbnailURL *string  `json:"thumbnail_url"`
		Price        *float64 `json:"price"`
		CategoryID   *string  `json:"category_id"`
		Duration     *string  `json:"duration"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Title != nil {
		updates = append(updates, "title = $"+strconv.Itoa(argCount))
		args = append(args, *req.Title)
		argCount++
		// Also update slug
		updates = append(updates, "slug = $"+strconv.Itoa(argCount))
		args = append(args, generateSlugFromTitle(*req.Title))
		argCount++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(argCount))
		args = append(args, *req.Description)
		argCount++
	}
	if req.ThumbnailURL != nil {
		updates = append(updates, "thumbnail_url = $"+strconv.Itoa(argCount))
		args = append(args, *req.ThumbnailURL)
		argCount++
	}
	if req.Price != nil {
		updates = append(updates, "price = $"+strconv.Itoa(argCount))
		args = append(args, *req.Price)
		argCount++
	}
	if req.CategoryID != nil {
		updates = append(updates, "category_id = $"+strconv.Itoa(argCount))
		if *req.CategoryID == "" {
			args = append(args, nil)
		} else {
			args = append(args, *req.CategoryID)
		}
		argCount++
	}
	if req.Duration != nil {
		updates = append(updates, "duration = $"+strconv.Itoa(argCount))
		args = append(args, *req.Duration)
		argCount++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tidak ada data yang diupdate"})
	}

	// If course was rejected, reset to draft on edit
	if currentStatus == CourseStatusRejected {
		updates = append(updates, "status = $"+strconv.Itoa(argCount))
		args = append(args, CourseStatusDraft)
		argCount++
	}

	updates = append(updates, "updated_at = $"+strconv.Itoa(argCount))
	args = append(args, time.Now())
	argCount++

	// Add course ID and instructor ID to args
	args = append(args, courseID, userID)

	query := "UPDATE courses SET " + joinStrings(updates, ", ") + 
		" WHERE id = $" + strconv.Itoa(argCount) + " AND instructor_id = $" + strconv.Itoa(argCount+1)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengupdate kursus"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Kursus berhasil diupdate"})
}

// InstructorDeleteCourse deletes a course (only draft status)
func InstructorDeleteCourse(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Check ownership and status
	var currentStatus string
	err = db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft') FROM courses WHERE id = $1 AND instructor_id = $2
	`, courseID, userID).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	// Can only delete draft courses
	if currentStatus != CourseStatusDraft && currentStatus != CourseStatusRejected {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Hanya kursus dengan status draft atau rejected yang dapat dihapus",
		})
	}

	_, err = db.DB.Exec(`DELETE FROM courses WHERE id = $1 AND instructor_id = $2`, courseID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus kursus"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Kursus berhasil dihapus"})
}

// InstructorSubmitCourse submits a course for admin review
func InstructorSubmitCourse(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Check ownership and current status
	var currentStatus string
	var lessonsCount int
	err = db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft'), COALESCE(lessons_count, 0) 
		FROM courses WHERE id = $1 AND instructor_id = $2
	`, courseID, userID).Scan(&currentStatus, &lessonsCount)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	// Can only submit draft or rejected courses
	if currentStatus != CourseStatusDraft && currentStatus != CourseStatusRejected {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Hanya kursus dengan status draft atau rejected yang dapat disubmit untuk review",
		})
	}

	// Check if course has at least 1 lesson
	var actualLessons int
	db.DB.QueryRow(`SELECT COUNT(*) FROM lessons WHERE course_id = $1`, courseID).Scan(&actualLessons)
	if actualLessons == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Kursus harus memiliki minimal 1 materi sebelum disubmit",
		})
	}

	// Update status to pending_review
	_, err = db.DB.Exec(`
		UPDATE courses 
		SET status = $1, submitted_at = NOW(), review_notes = NULL, updated_at = NOW()
		WHERE id = $2 AND instructor_id = $3
	`, CourseStatusPendingReview, courseID, userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal submit kursus"})
	}

	// Create notification for admins
	db.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, reference_type)
		SELECT id, 'course_submitted', 'Kursus Baru Menunggu Review', 
		       'Ada kursus baru yang menunggu review dari instruktur', $1, 'course'
		FROM users WHERE role = 'admin'
	`, courseID)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kursus berhasil disubmit untuk review. Admin akan mereview kursus Anda.",
		"status":  CourseStatusPendingReview,
	})
}

// ========================================
// INSTRUCTOR LESSONS (reuse existing with ownership check)
// ========================================

// InstructorGetLessonTree gets lesson tree for instructor's course
func InstructorGetLessonTree(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("courseId")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1 AND instructor_id = $2)`, courseID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}

	// Reuse existing GetLessonTree handler logic
	initCourseRepos()
	lessons, err := lessonRepo.GetTree(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil materi"})
	}

	total, _ := lessonRepo.CountByCourse(courseID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lessons": lessons,
		"total":   total,
	})
}

// InstructorCreateLesson creates a lesson for instructor's course
func InstructorCreateLesson(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("courseId")

	// Verify ownership and status (can't add lessons to published course)
	var status string
	err = db.DB.QueryRow(`SELECT COALESCE(status, 'draft') FROM courses WHERE id = $1 AND instructor_id = $2`, courseID, userID).Scan(&status)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	if status == CourseStatusPublished {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Tidak dapat menambah materi ke kursus yang sudah dipublish"})
	}

	// Use shared CreateLesson logic
	return createLessonForCourse(c, courseID)
}

// InstructorUpdateLesson updates a lesson (with ownership check)
func InstructorUpdateLesson(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	lessonID := c.Param("id")

	// Verify ownership through course
	var status string
	err = db.DB.QueryRow(`
		SELECT COALESCE(c.status, 'draft') FROM lessons l
		JOIN courses c ON l.course_id = c.id
		WHERE l.id = $1 AND c.instructor_id = $2
	`, lessonID, userID).Scan(&status)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Materi tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa materi"})
	}

	if status == CourseStatusPublished {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Tidak dapat mengubah materi pada kursus yang sudah dipublish"})
	}

	// Use shared UpdateLesson logic
	return updateLessonByID(c, lessonID)
}

// InstructorDeleteLesson deletes a lesson (with ownership check)
func InstructorDeleteLesson(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	lessonID := c.Param("id")

	// Verify ownership through course
	var status string
	err = db.DB.QueryRow(`
		SELECT COALESCE(c.status, 'draft') FROM lessons l
		JOIN courses c ON l.course_id = c.id
		WHERE l.id = $1 AND c.instructor_id = $2
	`, lessonID, userID).Scan(&status)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Materi tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa materi"})
	}

	if status == CourseStatusPublished {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Tidak dapat menghapus materi pada kursus yang sudah dipublish"})
	}

	// Delete lesson
	_, err = db.DB.Exec(`DELETE FROM lessons WHERE id = $1`, lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus materi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Materi berhasil dihapus"})
}

// ========================================
// INSTRUCTOR ANALYTICS
// ========================================

// InstructorCourseStudents returns enrolled students for a course
func InstructorCourseStudents(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1 AND instructor_id = $2)`, courseID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}

	rows, err := db.DB.Query(`
		SELECT u.id, u.full_name, u.email, e.progress_percentage, e.enrolled_at, e.completed_at
		FROM enrollments e
		JOIN users u ON e.user_id = u.id
		WHERE e.course_id = $1
		ORDER BY e.enrolled_at DESC
	`, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data siswa"})
	}
	defer rows.Close()

	var students []map[string]interface{}
	for rows.Next() {
		var id, fullName, email string
		var progress int
		var enrolledAt time.Time
		var completedAt sql.NullTime

		rows.Scan(&id, &fullName, &email, &progress, &enrolledAt, &completedAt)

		student := map[string]interface{}{
			"id":          id,
			"full_name":   fullName,
			"email":       email,
			"progress":    progress,
			"enrolled_at": enrolledAt,
		}
		if completedAt.Valid {
			student["completed_at"] = completedAt.Time
		}
		students = append(students, student)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"students": students,
		"total":    len(students),
	})
}

// InstructorCourseRatings returns ratings for a course
func InstructorCourseRatings(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1 AND instructor_id = $2)`, courseID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}

	// Get stats
	var avgRating float64
	var totalRatings int
	db.DB.QueryRow(`SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM course_ratings WHERE course_id = $1`, courseID).Scan(&avgRating, &totalRatings)

	// Get individual ratings
	rows, err := db.DB.Query(`
		SELECT u.full_name, r.rating, r.review, r.created_at
		FROM course_ratings r
		JOIN users u ON r.user_id = u.id
		WHERE r.course_id = $1
		ORDER BY r.created_at DESC
		LIMIT 20
	`, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil rating"})
	}
	defer rows.Close()

	var ratings []map[string]interface{}
	for rows.Next() {
		var fullName string
		var rating int
		var review sql.NullString
		var createdAt time.Time

		rows.Scan(&fullName, &rating, &review, &createdAt)

		r := map[string]interface{}{
			"user_name":  fullName,
			"rating":     rating,
			"created_at": createdAt,
		}
		if review.Valid {
			r["review"] = review.String
		}
		ratings = append(ratings, r)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"average_rating": avgRating,
		"total_ratings":  totalRatings,
		"ratings":        ratings,
	})
}

// ========================================
// INSTRUCTOR NOTIFICATIONS
// ========================================

// InstructorNotifications returns notifications for instructor
func InstructorNotifications(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	rows, err := db.DB.Query(`
		SELECT id, type, title, message, reference_id, reference_type, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil notifikasi"})
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id, notifType, title string
		var message, refID, refType sql.NullString
		var isRead bool
		var createdAt time.Time

		rows.Scan(&id, &notifType, &title, &message, &refID, &refType, &isRead, &createdAt)

		n := map[string]interface{}{
			"id":         id,
			"type":       notifType,
			"title":      title,
			"is_read":    isRead,
			"created_at": createdAt,
		}
		if message.Valid {
			n["message"] = message.String
		}
		if refID.Valid {
			n["reference_id"] = refID.String
		}
		if refType.Valid {
			n["reference_type"] = refType.String
		}
		notifications = append(notifications, n)
	}

	// Get unread count
	var unreadCount int
	db.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`, userID).Scan(&unreadCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// MarkNotificationRead marks a notification as read
func MarkNotificationRead(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	notifID := c.Param("id")

	_, err = db.DB.Exec(`UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`, notifID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update notifikasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notifikasi ditandai sudah dibaca"})
}

// ========================================
// INSTRUCTOR QUIZ MANAGEMENT
// ========================================

// InstructorCreateQuiz creates a new quiz for a lesson
func InstructorCreateQuiz(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	lessonID := c.Param("lessonId")

	// Verify ownership through lesson -> course -> instructor
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM lessons l 
			JOIN courses c ON l.course_id = c.id 
			WHERE l.id = $1 AND c.instructor_id = $2
		)
	`, lessonID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	var req domain.Quiz
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
	}

	// Create quiz
	var quizID string
	err = db.DB.QueryRow(`
		INSERT INTO quizzes (lesson_id, title, description, time_limit, passing_score, max_attempts, show_correct_answers)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, lessonID, req.Title, req.Description, req.TimeLimit, req.PassingScore, req.MaxAttempts, req.ShowCorrectAnswers).Scan(&quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat kuis"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      quizID,
		"message": "Kuis berhasil dibuat",
	})
}

// InstructorGetQuiz returns quiz for a lesson
func InstructorGetQuiz(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	lessonID := c.Param("lessonId")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM lessons l 
			JOIN courses c ON l.course_id = c.id 
			WHERE l.id = $1 AND c.instructor_id = $2
		)
	`, lessonID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	// Get quiz
	var quiz struct {
		ID                 string `json:"id"`
		Title              string `json:"title"`
		Description        string `json:"description"`
		TimeLimit          int    `json:"time_limit"`
		PassingScore       int    `json:"passing_score"`
		MaxAttempts        int    `json:"max_attempts"`
		ShowCorrectAnswers bool   `json:"show_correct_answers"`
	}

	err = db.DB.QueryRow(`
		SELECT id, title, COALESCE(description, ''), time_limit, passing_score, max_attempts, show_correct_answers
		FROM quizzes 
		WHERE lesson_id = $1
	`, lessonID).Scan(&quiz.ID, &quiz.Title, &quiz.Description, &quiz.TimeLimit, &quiz.PassingScore, &quiz.MaxAttempts, &quiz.ShowCorrectAnswers)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Kuis tidak ditemukan"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil kuis"})
	}

	// Get questions
	rows, err := db.DB.Query(`
		SELECT id, question_type, question_text, COALESCE(explanation, ''), points, required, order_index
		FROM quiz_questions
		WHERE quiz_id = $1
		ORDER BY order_index
	`, quiz.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil soal"})
	}
	defer rows.Close()

	var questions []map[string]interface{}
	for rows.Next() {
		var qID, qType, qText, explanation string
		var points, orderIndex int
		var required bool

		rows.Scan(&qID, &qType, &qText, &explanation, &points, &required, &orderIndex)

		// Get options for this question
		optRows, _ := db.DB.Query(`
			SELECT id, option_text, is_correct
			FROM quiz_options
			WHERE question_id = $1
			ORDER BY id
		`, qID)

		var options []map[string]interface{}
		if optRows != nil {
			for optRows.Next() {
				var optID, optText string
				var isCorrect bool
				optRows.Scan(&optID, &optText, &isCorrect)
				options = append(options, map[string]interface{}{
					"id":          optID,
					"option_text": optText,
					"is_correct":  isCorrect,
				})
			}
			optRows.Close()
		}

		questions = append(questions, map[string]interface{}{
			"id":            qID,
			"question_type": qType,
			"question_text": qText,
			"explanation":   explanation,
			"points":        points,
			"required":      required,
			"order_index":   orderIndex,
			"options":       options,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"quiz":      quiz,
		"questions": questions,
	})
}

// InstructorUpdateQuiz updates a quiz
func InstructorUpdateQuiz(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	quizID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM quizzes q 
			JOIN lessons l ON q.lesson_id = l.id 
			JOIN courses c ON l.course_id = c.id 
			WHERE q.id = $1 AND c.instructor_id = $2
		)
	`, quizID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	var req domain.Quiz
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
	}

	_, err = db.DB.Exec(`
		UPDATE quizzes 
		SET title = $1, description = $2, time_limit = $3, passing_score = $4, max_attempts = $5, show_correct_answers = $6, updated_at = NOW()
		WHERE id = $7
	`, req.Title, req.Description, req.TimeLimit, req.PassingScore, req.MaxAttempts, req.ShowCorrectAnswers, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update kuis"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Kuis berhasil diupdate"})
}

// InstructorDeleteQuiz deletes a quiz
func InstructorDeleteQuiz(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	quizID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM quizzes q 
			JOIN lessons l ON q.lesson_id = l.id 
			JOIN courses c ON l.course_id = c.id 
			WHERE q.id = $1 AND c.instructor_id = $2
		)
	`, quizID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	_, err = db.DB.Exec(`DELETE FROM quizzes WHERE id = $1`, quizID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal hapus kuis"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Kuis berhasil dihapus"})
}

// InstructorAddQuestion adds a question to a quiz
func InstructorAddQuestion(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	quizID := c.Param("quizId")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM quizzes q 
			JOIN lessons l ON q.lesson_id = l.id 
			JOIN courses c ON l.course_id = c.id 
			WHERE q.id = $1 AND c.instructor_id = $2
		)
	`, quizID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	var req struct {
		QuestionType string `json:"question_type"`
		QuestionText string `json:"question_text"`
		Explanation  string `json:"explanation"`
		Points       int    `json:"points"`
		Required     bool   `json:"required"`
		Options      []struct {
			OptionText string `json:"option_text"`
			IsCorrect  bool   `json:"is_correct"`
		} `json:"options"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
	}

	// Get next order index
	var maxOrder int
	db.DB.QueryRow(`SELECT COALESCE(MAX(order_index), 0) FROM quiz_questions WHERE quiz_id = $1`, quizID).Scan(&maxOrder)

	// Insert question
	var questionID string
	err = db.DB.QueryRow(`
		INSERT INTO quiz_questions (quiz_id, question_type, question_text, explanation, points, required, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, quizID, req.QuestionType, req.QuestionText, req.Explanation, req.Points, req.Required, maxOrder+1).Scan(&questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menambah soal"})
	}

	// Insert options
	for _, opt := range req.Options {
		_, err = db.DB.Exec(`
			INSERT INTO quiz_options (question_id, option_text, is_correct)
			VALUES ($1, $2, $3)
		`, questionID, opt.OptionText, opt.IsCorrect)
		if err != nil {
			// Log but don't fail
			continue
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      questionID,
		"message": "Soal berhasil ditambahkan",
	})
}

// InstructorUpdateQuestion updates a question
func InstructorUpdateQuestion(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	questionID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM quiz_questions qs
			JOIN quizzes q ON qs.quiz_id = q.id 
			JOIN lessons l ON q.lesson_id = l.id 
			JOIN courses c ON l.course_id = c.id 
			WHERE qs.id = $1 AND c.instructor_id = $2
		)
	`, questionID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	var req struct {
		QuestionType string `json:"question_type"`
		QuestionText string `json:"question_text"`
		Explanation  string `json:"explanation"`
		Points       int    `json:"points"`
		Required     bool   `json:"required"`
		Options      []struct {
			ID         string `json:"id"`
			OptionText string `json:"option_text"`
			IsCorrect  bool   `json:"is_correct"`
		} `json:"options"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
	}

	// Update question
	_, err = db.DB.Exec(`
		UPDATE quiz_questions 
		SET question_type = $1, question_text = $2, explanation = $3, points = $4, required = $5
		WHERE id = $6
	`, req.QuestionType, req.QuestionText, req.Explanation, req.Points, req.Required, questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update soal"})
	}

	// Delete old options and insert new ones
	db.DB.Exec(`DELETE FROM quiz_options WHERE question_id = $1`, questionID)
	for _, opt := range req.Options {
		db.DB.Exec(`
			INSERT INTO quiz_options (question_id, option_text, is_correct)
			VALUES ($1, $2, $3)
		`, questionID, opt.OptionText, opt.IsCorrect)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Soal berhasil diupdate"})
}

// InstructorDeleteQuestion deletes a question
func InstructorDeleteQuestion(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	questionID := c.Param("id")

	// Verify ownership
	var exists bool
	err = db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM quiz_questions qs
			JOIN quizzes q ON qs.quiz_id = q.id 
			JOIN lessons l ON q.lesson_id = l.id 
			JOIN courses c ON l.course_id = c.id 
			WHERE qs.id = $1 AND c.instructor_id = $2
		)
	`, questionID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akses ditolak"})
	}

	_, err = db.DB.Exec(`DELETE FROM quiz_questions WHERE id = $1`, questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal hapus soal"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Soal berhasil dihapus"})
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func generateSlugFromTitle(title string) string {
	// Simple slug generation - reuse from existing code if available
	slug := ""
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			slug += string(c)
		} else if c >= 'A' && c <= 'Z' {
			slug += string(c + 32) // lowercase
		} else if c == ' ' || c == '-' {
			if len(slug) > 0 && slug[len(slug)-1] != '-' {
				slug += "-"
			}
		}
	}
	return slug
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}


