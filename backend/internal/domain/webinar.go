package domain

import "time"

// Webinar represents a live webinar session linked to a course
type Webinar struct {
	ID              string     `json:"id" db:"id"`
	CourseID        string     `json:"course_id" db:"course_id"`
	Title           string     `json:"title" db:"title"`
	Description     *string    `json:"description,omitempty" db:"description"`
	ScheduledAt     time.Time  `json:"scheduled_at" db:"scheduled_at"`
	DurationMinutes int        `json:"duration_minutes" db:"duration_minutes"`
	MeetingURL      *string    `json:"meeting_url,omitempty" db:"meeting_url"`
	MeetingPassword *string    `json:"meeting_password,omitempty" db:"meeting_password"`
	MaxParticipants *int       `json:"max_participants,omitempty" db:"max_participants"`
	Status          string     `json:"status" db:"status"`
	RecordingURL    *string    `json:"recording_url,omitempty" db:"recording_url"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Joined data (not stored in DB)
	Course             *Course `json:"course,omitempty"`
	RegistrationsCount int     `json:"registrations_count,omitempty"`
}

// WebinarRegistration tracks user registration for webinar
type WebinarRegistration struct {
	ID                 string     `json:"id" db:"id"`
	WebinarID          string     `json:"webinar_id" db:"webinar_id"`
	UserID             string     `json:"user_id" db:"user_id"`
	RegisteredAt       time.Time  `json:"registered_at" db:"registered_at"`
	RegistrationSource string     `json:"registration_source" db:"registration_source"`
	Attended           bool       `json:"attended" db:"attended"`
	AttendedAt         *time.Time `json:"attended_at,omitempty" db:"attended_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`

	// Joined data
	User    *User    `json:"user,omitempty"`
	Webinar *Webinar `json:"webinar,omitempty"`
}

// WebinarReminder tracks scheduled reminders for webinar
type WebinarReminder struct {
	ID           string     `json:"id" db:"id"`
	WebinarID    string     `json:"webinar_id" db:"webinar_id"`
	UserID       string     `json:"user_id" db:"user_id"`
	ReminderType string     `json:"reminder_type" db:"reminder_type"`
	ScheduledAt  time.Time  `json:"scheduled_at" db:"scheduled_at"`
	SentAt       *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	Channel      string     `json:"channel" db:"channel"`
	Status       string     `json:"status" db:"status"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`

	// Joined data
	User    *User    `json:"user,omitempty"`
	Webinar *Webinar `json:"webinar,omitempty"`
}

// WANotification tracks WhatsApp notifications sent
type WANotification struct {
	ID           string     `json:"id" db:"id"`
	UserID       *string    `json:"user_id,omitempty" db:"user_id"`
	Phone        string     `json:"phone" db:"phone"`
	MessageType  string     `json:"message_type" db:"message_type"`
	TemplateData []byte     `json:"template_data,omitempty" db:"template_data"`
	Status       string     `json:"status" db:"status"`
	SentAt       *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Webinar status constants
const (
	WebinarStatusDraft     = "draft"
	WebinarStatusUpcoming  = "upcoming"
	WebinarStatusLive      = "live"
	WebinarStatusCompleted = "completed"
	WebinarStatusCancelled = "cancelled"
)

// Reminder type constants
const (
	ReminderType1Day   = "1_day"
	ReminderType3Hours = "3_hours"
	ReminderType30Min  = "30_min"
)

// Reminder status constants
const (
	ReminderStatusPending = "pending"
	ReminderStatusSent    = "sent"
	ReminderStatusFailed  = "failed"
)

// WA message type constants
const (
	WAMessageWebinarConfirmation = "webinar_confirmation"
	WAMessageReminder1Day        = "webinar_reminder_1d"
	WAMessageReminder3Hours      = "webinar_reminder_3h"
	WAMessageReminder30Min       = "webinar_reminder_30m"
	WAMessagePaymentSuccess      = "payment_success"
	WAMessageCredentials         = "credentials"
)

// Course type constants
const (
	CourseTypeSelfPaced = "self_paced"
	CourseTypeWebinar   = "webinar"
	CourseTypeHybrid    = "hybrid"
)

// CreateWebinarRequest represents request to create a webinar
type CreateWebinarRequest struct {
	CourseID        string  `json:"course_id" validate:"required"`
	Title           string  `json:"title" validate:"required"`
	Description     *string `json:"description,omitempty"`
	ScheduledAt     string  `json:"scheduled_at" validate:"required"` // ISO 8601 format
	DurationMinutes int     `json:"duration_minutes"`
	MeetingURL      *string `json:"meeting_url,omitempty"`
	MeetingPassword *string `json:"meeting_password,omitempty"`
	MaxParticipants *int    `json:"max_participants,omitempty"`
	Status          string  `json:"status,omitempty"`
}

// UpdateWebinarRequest represents request to update a webinar
type UpdateWebinarRequest struct {
	Title           *string `json:"title,omitempty"`
	Description     *string `json:"description,omitempty"`
	ScheduledAt     *string `json:"scheduled_at,omitempty"`
	DurationMinutes *int    `json:"duration_minutes,omitempty"`
	MeetingURL      *string `json:"meeting_url,omitempty"`
	MeetingPassword *string `json:"meeting_password,omitempty"`
	MaxParticipants *int    `json:"max_participants,omitempty"`
	Status          *string `json:"status,omitempty"`
	RecordingURL    *string `json:"recording_url,omitempty"`
}
