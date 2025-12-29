package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// WebinarRepository handles webinar data operations
type WebinarRepository struct {
	db *sqlx.DB
}

// NewWebinarRepository creates a new webinar repository
func NewWebinarRepository(db *sqlx.DB) *WebinarRepository {
	return &WebinarRepository{db: db}
}

// Create creates a new webinar
func (r *WebinarRepository) Create(w *domain.Webinar) error {
	query := `
		INSERT INTO webinars (course_id, title, description, scheduled_at, duration_minutes, 
			meeting_url, meeting_password, max_participants, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		w.CourseID, w.Title, w.Description, w.ScheduledAt, w.DurationMinutes,
		w.MeetingURL, w.MeetingPassword, w.MaxParticipants, w.Status,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

// GetByID retrieves a webinar by ID
func (r *WebinarRepository) GetByID(id string) (*domain.Webinar, error) {
	query := `
		SELECT w.id, w.course_id, w.title, w.description, w.scheduled_at, w.duration_minutes,
			w.meeting_url, w.meeting_password, w.max_participants, w.status, w.recording_url,
			w.created_at, w.updated_at,
			c.id, c.title, c.slug, c.thumbnail_url,
			(SELECT COUNT(*) FROM webinar_registrations WHERE webinar_id = w.id) as registrations_count
		FROM webinars w
		LEFT JOIN courses c ON w.course_id = c.id
		WHERE w.id = $1
	`
	
	var w domain.Webinar
	var courseID, courseTitle, courseSlug sql.NullString
	var courseThumbnail sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&w.ID, &w.CourseID, &w.Title, &w.Description, &w.ScheduledAt, &w.DurationMinutes,
		&w.MeetingURL, &w.MeetingPassword, &w.MaxParticipants, &w.Status, &w.RecordingURL,
		&w.CreatedAt, &w.UpdatedAt,
		&courseID, &courseTitle, &courseSlug, &courseThumbnail,
		&w.RegistrationsCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if courseID.Valid {
		w.Course = &domain.Course{
			ID:           courseID.String,
			Title:        courseTitle.String,
			Slug:         courseSlug.String,
			ThumbnailURL: &courseThumbnail.String,
		}
	}
	
	return &w, nil
}

// GetByCourseID retrieves all webinars for a course
func (r *WebinarRepository) GetByCourseID(courseID string) ([]*domain.Webinar, error) {
	query := `
		SELECT w.id, w.course_id, w.title, w.description, w.scheduled_at, w.duration_minutes,
			w.meeting_url, w.meeting_password, w.max_participants, w.status, w.recording_url,
			w.created_at, w.updated_at,
			(SELECT COUNT(*) FROM webinar_registrations WHERE webinar_id = w.id) as registrations_count
		FROM webinars w
		WHERE w.course_id = $1
		ORDER BY w.scheduled_at ASC
	`
	
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var webinars []*domain.Webinar
	for rows.Next() {
		var w domain.Webinar
		err := rows.Scan(
			&w.ID, &w.CourseID, &w.Title, &w.Description, &w.ScheduledAt, &w.DurationMinutes,
			&w.MeetingURL, &w.MeetingPassword, &w.MaxParticipants, &w.Status, &w.RecordingURL,
			&w.CreatedAt, &w.UpdatedAt, &w.RegistrationsCount,
		)
		if err != nil {
			return nil, err
		}
		webinars = append(webinars, &w)
	}
	
	return webinars, nil
}

// GetUpcomingByCourse retrieves upcoming webinars for a course
func (r *WebinarRepository) GetUpcomingByCourse(courseID string) ([]*domain.Webinar, error) {
	query := `
		SELECT w.id, w.course_id, w.title, w.description, w.scheduled_at, w.duration_minutes,
			w.meeting_url, w.meeting_password, w.max_participants, w.status, w.recording_url,
			w.created_at, w.updated_at,
			(SELECT COUNT(*) FROM webinar_registrations WHERE webinar_id = w.id) as registrations_count
		FROM webinars w
		WHERE w.course_id = $1 
			AND w.status = 'upcoming' 
			AND w.scheduled_at > NOW()
		ORDER BY w.scheduled_at ASC
	`
	
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var webinars []*domain.Webinar
	for rows.Next() {
		var w domain.Webinar
		err := rows.Scan(
			&w.ID, &w.CourseID, &w.Title, &w.Description, &w.ScheduledAt, &w.DurationMinutes,
			&w.MeetingURL, &w.MeetingPassword, &w.MaxParticipants, &w.Status, &w.RecordingURL,
			&w.CreatedAt, &w.UpdatedAt, &w.RegistrationsCount,
		)
		if err != nil {
			return nil, err
		}
		webinars = append(webinars, &w)
	}
	
	return webinars, nil
}

// List retrieves all webinars with pagination
func (r *WebinarRepository) List(limit, offset int) ([]*domain.Webinar, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM webinars`
	r.db.QueryRow(countQuery).Scan(&total)
	
	query := `
		SELECT w.id, w.course_id, w.title, w.description, w.scheduled_at, w.duration_minutes,
			w.meeting_url, w.meeting_password, w.max_participants, w.status, w.recording_url,
			w.created_at, w.updated_at,
			c.id, c.title, c.slug,
			(SELECT COUNT(*) FROM webinar_registrations WHERE webinar_id = w.id) as registrations_count
		FROM webinars w
		LEFT JOIN courses c ON w.course_id = c.id
		ORDER BY w.scheduled_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var webinars []*domain.Webinar
	for rows.Next() {
		var w domain.Webinar
		var courseID, courseTitle, courseSlug sql.NullString
		
		err := rows.Scan(
			&w.ID, &w.CourseID, &w.Title, &w.Description, &w.ScheduledAt, &w.DurationMinutes,
			&w.MeetingURL, &w.MeetingPassword, &w.MaxParticipants, &w.Status, &w.RecordingURL,
			&w.CreatedAt, &w.UpdatedAt,
			&courseID, &courseTitle, &courseSlug,
			&w.RegistrationsCount,
		)
		if err != nil {
			return nil, 0, err
		}
		
		if courseID.Valid {
			w.Course = &domain.Course{
				ID:    courseID.String,
				Title: courseTitle.String,
				Slug:  courseSlug.String,
			}
		}
		
		webinars = append(webinars, &w)
	}
	
	return webinars, total, nil
}

// Update updates a webinar
func (r *WebinarRepository) Update(w *domain.Webinar) error {
	query := `
		UPDATE webinars SET
			title = $2,
			description = $3,
			scheduled_at = $4,
			duration_minutes = $5,
			meeting_url = $6,
			meeting_password = $7,
			max_participants = $8,
			status = $9,
			recording_url = $10,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		w.ID, w.Title, w.Description, w.ScheduledAt, w.DurationMinutes,
		w.MeetingURL, w.MeetingPassword, w.MaxParticipants, w.Status, w.RecordingURL,
	)
	return err
}

// UpdateStatus updates webinar status
func (r *WebinarRepository) UpdateStatus(id, status string) error {
	query := `UPDATE webinars SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, status)
	return err
}

// Delete deletes a webinar
func (r *WebinarRepository) Delete(id string) error {
	query := `DELETE FROM webinars WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// RegisterUser registers a user for a webinar
func (r *WebinarRepository) RegisterUser(webinarID, userID, source string) error {
	query := `
		INSERT INTO webinar_registrations (webinar_id, user_id, registration_source)
		VALUES ($1, $2, $3)
		ON CONFLICT (webinar_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(query, webinarID, userID, source)
	return err
}

// IsUserRegistered checks if user is registered for webinar
func (r *WebinarRepository) IsUserRegistered(webinarID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM webinar_registrations WHERE webinar_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, webinarID, userID).Scan(&exists)
	return exists, err
}

// GetRegistrations retrieves all registrations for a webinar
func (r *WebinarRepository) GetRegistrations(webinarID string) ([]*domain.WebinarRegistration, error) {
	query := `
		SELECT wr.id, wr.webinar_id, wr.user_id, wr.registered_at, wr.registration_source,
			wr.attended, wr.attended_at, wr.created_at,
			u.full_name, u.email, u.phone
		FROM webinar_registrations wr
		LEFT JOIN users u ON wr.user_id = u.id
		WHERE wr.webinar_id = $1
		ORDER BY wr.registered_at DESC
	`
	
	rows, err := r.db.Query(query, webinarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var registrations []*domain.WebinarRegistration
	for rows.Next() {
		var reg domain.WebinarRegistration
		var userName, userEmail sql.NullString
		var userPhone sql.NullString
		
		err := rows.Scan(
			&reg.ID, &reg.WebinarID, &reg.UserID, &reg.RegisteredAt, &reg.RegistrationSource,
			&reg.Attended, &reg.AttendedAt, &reg.CreatedAt,
			&userName, &userEmail, &userPhone,
		)
		if err != nil {
			return nil, err
		}
		
		if userName.Valid {
			phone := userPhone.String
			reg.User = &domain.User{
				ID:       reg.UserID,
				FullName: userName.String,
				Email:    userEmail.String,
				Phone:    &phone,
			}
		}
		
		registrations = append(registrations, &reg)
	}
	
	return registrations, nil
}

// MarkAttendance marks user attendance for webinar
func (r *WebinarRepository) MarkAttendance(webinarID, userID string) error {
	query := `
		UPDATE webinar_registrations SET attended = true, attended_at = NOW()
		WHERE webinar_id = $1 AND user_id = $2
	`
	_, err := r.db.Exec(query, webinarID, userID)
	return err
}

// GetUserWebinars retrieves webinars a user is registered for
func (r *WebinarRepository) GetUserWebinars(userID string, includeCompleted bool) ([]*domain.Webinar, error) {
	statusFilter := "AND w.status IN ('upcoming', 'live')"
	if includeCompleted {
		statusFilter = ""
	}
	
	query := `
		SELECT w.id, w.course_id, w.title, w.description, w.scheduled_at, w.duration_minutes,
			w.meeting_url, w.meeting_password, w.max_participants, w.status, w.recording_url,
			w.created_at, w.updated_at,
			c.id, c.title, c.slug
		FROM webinars w
		INNER JOIN webinar_registrations wr ON w.id = wr.webinar_id
		LEFT JOIN courses c ON w.course_id = c.id
		WHERE wr.user_id = $1 ` + statusFilter + `
		ORDER BY w.scheduled_at ASC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var webinars []*domain.Webinar
	for rows.Next() {
		var w domain.Webinar
		var courseID, courseTitle, courseSlug sql.NullString
		
		err := rows.Scan(
			&w.ID, &w.CourseID, &w.Title, &w.Description, &w.ScheduledAt, &w.DurationMinutes,
			&w.MeetingURL, &w.MeetingPassword, &w.MaxParticipants, &w.Status, &w.RecordingURL,
			&w.CreatedAt, &w.UpdatedAt,
			&courseID, &courseTitle, &courseSlug,
		)
		if err != nil {
			return nil, err
		}
		
		if courseID.Valid {
			w.Course = &domain.Course{
				ID:    courseID.String,
				Title: courseTitle.String,
				Slug:  courseSlug.String,
			}
		}
		
		webinars = append(webinars, &w)
	}
	
	return webinars, nil
}

// CreateReminders creates reminder entries for a user's webinar registration
func (r *WebinarRepository) CreateReminders(webinarID, userID string, scheduledAt time.Time) error {
	// Create 1 day reminder
	reminder1Day := scheduledAt.Add(-24 * time.Hour)
	if reminder1Day.After(time.Now()) {
		_, err := r.db.Exec(`
			INSERT INTO webinar_reminders (webinar_id, user_id, reminder_type, scheduled_at, channel, status)
			VALUES ($1, $2, $3, $4, 'whatsapp', 'pending')
			ON CONFLICT DO NOTHING
		`, webinarID, userID, domain.ReminderType1Day, reminder1Day)
		if err != nil {
			log.Printf("[WebinarRepo] Failed to create 1-day reminder: %v", err)
		}
	}
	
	// Create 3 hours reminder
	reminder3Hours := scheduledAt.Add(-3 * time.Hour)
	if reminder3Hours.After(time.Now()) {
		_, err := r.db.Exec(`
			INSERT INTO webinar_reminders (webinar_id, user_id, reminder_type, scheduled_at, channel, status)
			VALUES ($1, $2, $3, $4, 'whatsapp', 'pending')
			ON CONFLICT DO NOTHING
		`, webinarID, userID, domain.ReminderType3Hours, reminder3Hours)
		if err != nil {
			log.Printf("[WebinarRepo] Failed to create 3-hours reminder: %v", err)
		}
	}
	
	// Create 30 minutes reminder
	reminder30Min := scheduledAt.Add(-30 * time.Minute)
	if reminder30Min.After(time.Now()) {
		_, err := r.db.Exec(`
			INSERT INTO webinar_reminders (webinar_id, user_id, reminder_type, scheduled_at, channel, status)
			VALUES ($1, $2, $3, $4, 'whatsapp', 'pending')
			ON CONFLICT DO NOTHING
		`, webinarID, userID, domain.ReminderType30Min, reminder30Min)
		if err != nil {
			log.Printf("[WebinarRepo] Failed to create 30-min reminder: %v", err)
		}
	}
	
	return nil
}

// GetPendingReminders retrieves pending reminders that are due
func (r *WebinarRepository) GetPendingReminders(before time.Time) ([]*domain.WebinarReminder, error) {
	query := `
		SELECT r.id, r.webinar_id, r.user_id, r.reminder_type, r.scheduled_at, 
			r.sent_at, r.channel, r.status, r.error_message, r.created_at,
			w.title, w.scheduled_at, w.meeting_url, w.meeting_password,
			u.full_name, u.email, u.phone
		FROM webinar_reminders r
		INNER JOIN webinars w ON r.webinar_id = w.id
		INNER JOIN users u ON r.user_id = u.id
		WHERE r.status = 'pending' AND r.scheduled_at <= $1
		ORDER BY r.scheduled_at ASC
		LIMIT 100
	`
	
	rows, err := r.db.Query(query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reminders []*domain.WebinarReminder
	for rows.Next() {
		var rem domain.WebinarReminder
		var webinarTitle string
		var webinarScheduledAt time.Time
		var meetingURL, meetingPassword sql.NullString
		var userName, userEmail sql.NullString
		var userPhone sql.NullString
		
		err := rows.Scan(
			&rem.ID, &rem.WebinarID, &rem.UserID, &rem.ReminderType, &rem.ScheduledAt,
			&rem.SentAt, &rem.Channel, &rem.Status, &rem.ErrorMessage, &rem.CreatedAt,
			&webinarTitle, &webinarScheduledAt, &meetingURL, &meetingPassword,
			&userName, &userEmail, &userPhone,
		)
		if err != nil {
			return nil, err
		}
		
		rem.Webinar = &domain.Webinar{
			ID:              rem.WebinarID,
			Title:           webinarTitle,
			ScheduledAt:     webinarScheduledAt,
			MeetingURL:      &meetingURL.String,
			MeetingPassword: &meetingPassword.String,
		}
		
		phone := userPhone.String
		rem.User = &domain.User{
			ID:       rem.UserID,
			FullName: userName.String,
			Email:    userEmail.String,
			Phone:    &phone,
		}
		
		reminders = append(reminders, &rem)
	}
	
	return reminders, nil
}

// MarkReminderSent marks a reminder as sent
func (r *WebinarRepository) MarkReminderSent(reminderID string) error {
	query := `UPDATE webinar_reminders SET status = 'sent', sent_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, reminderID)
	return err
}

// MarkReminderFailed marks a reminder as failed
func (r *WebinarRepository) MarkReminderFailed(reminderID, errorMsg string) error {
	query := `UPDATE webinar_reminders SET status = 'failed', error_message = $2 WHERE id = $1`
	_, err := r.db.Exec(query, reminderID, errorMsg)
	return err
}
