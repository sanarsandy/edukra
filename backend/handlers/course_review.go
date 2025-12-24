package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// ========================================
// ADMIN COURSE REVIEW HANDLERS
// ========================================

// AdminListPendingReviews returns courses pending review
func AdminListPendingReviews(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := db.DB.Query(`
		SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.price, 
		       c.status, COALESCE(c.lessons_count, 0), c.submitted_at, 
			   u.id as instructor_id, u.full_name as instructor_name, u.email as instructor_email,
			   cat.name as category_name
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.status = 'pending_review'
		ORDER BY c.submitted_at ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data kursus"})
	}
	defer rows.Close()

	var courses []map[string]interface{}
	for rows.Next() {
		var id, title, slug, description, status string
		var thumbnailURL, instructorID, instructorName, instructorEmail, categoryName sql.NullString
		var price float64
		var lessonsCount int
		var submittedAt sql.NullTime

		err := rows.Scan(&id, &title, &slug, &description, &thumbnailURL, &price,
			&status, &lessonsCount, &submittedAt,
			&instructorID, &instructorName, &instructorEmail, &categoryName)
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
		}
		if thumbnailURL.Valid {
			course["thumbnail_url"] = thumbnailURL.String
		}
		if submittedAt.Valid {
			course["submitted_at"] = submittedAt.Time
		}
		if instructorID.Valid {
			course["instructor"] = map[string]string{
				"id":    instructorID.String,
				"name":  instructorName.String,
				"email": instructorEmail.String,
			}
		}
		if categoryName.Valid {
			course["category_name"] = categoryName.String
		}

		courses = append(courses, course)
	}

	// Get total count
	var total int
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE status = 'pending_review'`).Scan(&total)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"courses": courses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// AdminReviewCourseDetail returns detailed info for reviewing a course
func AdminReviewCourseDetail(c echo.Context) error {
	courseID := c.Param("id")

	var course struct {
		ID             string    `json:"id"`
		Title          string    `json:"title"`
		Slug           string    `json:"slug"`
		Description    string    `json:"description"`
		ThumbnailURL   *string   `json:"thumbnail_url"`
		Price          float64   `json:"price"`
		Currency       string    `json:"currency"`
		Status         string    `json:"status"`
		LessonsCount   int       `json:"lessons_count"`
		Duration       string    `json:"duration"`
		CategoryID     *string   `json:"category_id"`
		CategoryName   *string   `json:"category_name"`
		InstructorID   string    `json:"instructor_id"`
		InstructorName string    `json:"instructor_name"`
		SubmittedAt    *time.Time `json:"submitted_at"`
		CreatedAt      time.Time `json:"created_at"`
	}

	var thumbnailURL, categoryID, categoryName sql.NullString
	var submittedAt sql.NullTime

	err := db.DB.QueryRow(`
		SELECT c.id, c.title, c.slug, COALESCE(c.description, ''), c.thumbnail_url, c.price, c.currency,
		       COALESCE(c.status, 'draft'), COALESCE(c.lessons_count, 0), COALESCE(c.duration, ''),
			   c.category_id, cat.name, c.instructor_id, u.full_name, c.submitted_at, c.created_at
		FROM courses c
		LEFT JOIN categories cat ON c.category_id = cat.id
		JOIN users u ON c.instructor_id = u.id
		WHERE c.id = $1
	`, courseID).Scan(
		&course.ID, &course.Title, &course.Slug, &course.Description, &thumbnailURL,
		&course.Price, &course.Currency, &course.Status, &course.LessonsCount, &course.Duration,
		&categoryID, &categoryName, &course.InstructorID, &course.InstructorName,
		&submittedAt, &course.CreatedAt,
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
	if submittedAt.Valid {
		course.SubmittedAt = &submittedAt.Time
	}

	// Get lessons for this course
	initLessonRepos()
	lessons, _ := lessonRepo.GetTree(courseID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"course":  course,
		"lessons": lessons,
	})
}

// AdminApproveCourse approves a course (sets status to approved/published)
func AdminApproveCourse(c echo.Context) error {
	adminID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	var req struct {
		PublishDirectly bool   `json:"publish_directly"` // true = publish, false = approved only
		Notes           string `json:"notes"`
	}
	c.Bind(&req)

	// Verify course exists and is pending review
	var currentStatus string
	var instructorID string
	err = db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft'), instructor_id FROM courses WHERE id = $1
	`, courseID).Scan(&currentStatus, &instructorID)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	if currentStatus != CourseStatusPendingReview {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Hanya kursus dengan status pending_review yang dapat diapprove",
		})
	}

	newStatus := CourseStatusApproved
	isPublished := false
	if req.PublishDirectly {
		newStatus = CourseStatusPublished
		isPublished = true
	}

	_, err = db.DB.Exec(`
		UPDATE courses 
		SET status = $1, is_published = $2, reviewed_at = NOW(), reviewed_by = $3, review_notes = $4, updated_at = NOW()
		WHERE id = $5
	`, newStatus, isPublished, adminID, req.Notes, courseID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal approve kursus"})
	}

	// Notify instructor
	notifTitle := "Kursus Disetujui"
	notifMessage := "Kursus Anda telah disetujui oleh admin."
	if req.PublishDirectly {
		notifMessage = "Kursus Anda telah dipublish dan sekarang tersedia untuk siswa."
	}
	if req.Notes != "" {
		notifMessage += " Catatan: " + req.Notes
	}

	db.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, reference_type)
		VALUES ($1, 'course_approved', $2, $3, $4, 'course')
	`, instructorID, notifTitle, notifMessage, courseID)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kursus berhasil diapprove",
		"status":  newStatus,
	})
}

// AdminRejectCourse rejects a course with notes
func AdminRejectCourse(c echo.Context) error {
	adminID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	var req struct {
		Notes string `json:"notes"` // Required - reason for rejection
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	if req.Notes == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Alasan penolakan wajib diisi"})
	}

	// Verify course exists and is pending review
	var currentStatus string
	var instructorID string
	err = db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft'), instructor_id FROM courses WHERE id = $1
	`, courseID).Scan(&currentStatus, &instructorID)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	if currentStatus != CourseStatusPendingReview {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Hanya kursus dengan status pending_review yang dapat direject",
		})
	}

	_, err = db.DB.Exec(`
		UPDATE courses 
		SET status = $1, reviewed_at = NOW(), reviewed_by = $2, review_notes = $3, updated_at = NOW()
		WHERE id = $4
	`, CourseStatusRejected, adminID, req.Notes, courseID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal reject kursus"})
	}

	// Notify instructor
	db.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, reference_type)
		VALUES ($1, 'course_rejected', 'Kursus Ditolak', $2, $3, 'course')
	`, instructorID, "Kursus Anda ditolak. Alasan: "+req.Notes, courseID)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kursus berhasil ditolak",
		"status":  CourseStatusRejected,
	})
}

// AdminPublishCourse publishes an approved course
func AdminPublishCourse(c echo.Context) error {
	adminID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	courseID := c.Param("id")

	// Verify course exists and is approved
	var currentStatus string
	var instructorID string
	err = db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft'), instructor_id FROM courses WHERE id = $1
	`, courseID).Scan(&currentStatus, &instructorID)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	if currentStatus != CourseStatusApproved {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Hanya kursus dengan status approved yang dapat dipublish",
		})
	}

	_, err = db.DB.Exec(`
		UPDATE courses 
		SET status = 'published', is_published = true, reviewed_by = $1, updated_at = NOW()
		WHERE id = $2
	`, adminID, courseID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal publish kursus"})
	}

	// Notify instructor
	db.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, reference_type)
		VALUES ($1, 'course_published', 'Kursus Dipublish', 'Kursus Anda telah dipublish dan sekarang tersedia untuk siswa.', $2, 'course')
	`, instructorID, courseID)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kursus berhasil dipublish",
		"status":  CourseStatusPublished,
	})
}

// AdminUnpublishCourse unpublishes a published course (back to approved)
func AdminUnpublishCourse(c echo.Context) error {
	courseID := c.Param("id")

	var req struct {
		Notes string `json:"notes"`
	}
	c.Bind(&req)

	// Verify course exists and is published
	var currentStatus string
	var instructorID string
	err := db.DB.QueryRow(`
		SELECT COALESCE(status, 'draft'), instructor_id FROM courses WHERE id = $1
	`, courseID).Scan(&currentStatus, &instructorID)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kursus tidak ditemukan"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa kursus"})
	}

	if currentStatus != CourseStatusPublished {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Hanya kursus dengan status published yang dapat di-unpublish",
		})
	}

	_, err = db.DB.Exec(`
		UPDATE courses 
		SET status = 'approved', is_published = false, updated_at = NOW()
		WHERE id = $1
	`, courseID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal unpublish kursus"})
	}

	// Notify instructor
	message := "Kursus Anda telah di-unpublish oleh admin."
	if req.Notes != "" {
		message += " Alasan: " + req.Notes
	}
	db.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, reference_id, reference_type)
		VALUES ($1, 'course_unpublished', 'Kursus Di-unpublish', $2, $3, 'course')
	`, instructorID, message, courseID)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kursus berhasil di-unpublish",
		"status":  CourseStatusApproved,
	})
}

// ========================================
// ADMIN REVIEW STATS
// ========================================

// AdminReviewStats returns review statistics
func AdminReviewStats(c echo.Context) error {
	stats := struct {
		PendingReview int `json:"pending_review"`
		ApprovedToday int `json:"approved_today"`
		RejectedToday int `json:"rejected_today"`
		TotalReviewed int `json:"total_reviewed"`
	}{}

	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE status = 'pending_review'`).Scan(&stats.PendingReview)
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE status = 'approved' AND reviewed_at::date = CURRENT_DATE`).Scan(&stats.ApprovedToday)
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE status = 'rejected' AND reviewed_at::date = CURRENT_DATE`).Scan(&stats.RejectedToday)
	db.DB.QueryRow(`SELECT COUNT(*) FROM courses WHERE reviewed_at IS NOT NULL`).Scan(&stats.TotalReviewed)

	return c.JSON(http.StatusOK, stats)
}


