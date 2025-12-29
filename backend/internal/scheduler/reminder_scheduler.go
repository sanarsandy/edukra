package scheduler

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/service"
)

// ReminderScheduler handles sending webinar reminders
type ReminderScheduler struct {
	db           *sqlx.DB
	webinarRepo  *postgres.WebinarRepository
	waService    *service.WhatsAppService
	ticker       *time.Ticker
	done         chan bool
	isRunning    bool
	checkInterval time.Duration
}

// NewReminderScheduler creates a new reminder scheduler
func NewReminderScheduler(db *sqlx.DB) *ReminderScheduler {
	return &ReminderScheduler{
		db:            db,
		webinarRepo:   postgres.NewWebinarRepository(db),
		waService:     service.GetWhatsAppService(),
		done:          make(chan bool),
		isRunning:     false,
		checkInterval: 1 * time.Minute, // Check every minute
	}
}

// Start begins the scheduler loop
func (s *ReminderScheduler) Start() {
	if s.isRunning {
		log.Println("[Scheduler] Already running")
		return
	}

	s.ticker = time.NewTicker(s.checkInterval)
	s.isRunning = true

	go func() {
		log.Println("[Scheduler] Reminder scheduler started")
		
		// Process immediately on start
		s.processReminders()
		
		for {
			select {
			case <-s.done:
				log.Println("[Scheduler] Reminder scheduler stopped")
				return
			case <-s.ticker.C:
				s.processReminders()
			}
		}
	}()
}

// Stop stops the scheduler
func (s *ReminderScheduler) Stop() {
	if !s.isRunning {
		return
	}
	
	s.ticker.Stop()
	s.done <- true
	s.isRunning = false
}

// IsRunning returns whether the scheduler is running
func (s *ReminderScheduler) IsRunning() bool {
	return s.isRunning
}

// processReminders processes pending reminders
func (s *ReminderScheduler) processReminders() {
	// Get pending reminders that are due
	reminders, err := s.webinarRepo.GetPendingReminders(time.Now())
	if err != nil {
		log.Printf("[Scheduler] Error fetching pending reminders: %v", err)
		return
	}

	if len(reminders) == 0 {
		return // Nothing to process
	}

	log.Printf("[Scheduler] Processing %d pending reminders", len(reminders))

	for _, reminder := range reminders {
		s.sendReminder(reminder)
	}
}

// sendReminder sends a single reminder
func (s *ReminderScheduler) sendReminder(reminder *domain.WebinarReminder) {
	// Validate we have required data
	if reminder.User == nil || reminder.User.Phone == nil || *reminder.User.Phone == "" {
		log.Printf("[Scheduler] Skipping reminder %s: no phone number", reminder.ID)
		s.webinarRepo.MarkReminderFailed(reminder.ID, "No phone number")
		return
	}

	if reminder.Webinar == nil {
		log.Printf("[Scheduler] Skipping reminder %s: no webinar data", reminder.ID)
		s.webinarRepo.MarkReminderFailed(reminder.ID, "No webinar data")
		return
	}

	// Send the appropriate reminder based on type
	var err error
	phone := *reminder.User.Phone
	userName := reminder.User.FullName

	switch reminder.ReminderType {
	case domain.ReminderType1Day:
		log.Printf("[Scheduler] Sending H-1 reminder to %s for webinar %s", phone, reminder.WebinarID)
		err = s.waService.SendReminder1Day(phone, reminder.Webinar, userName)
		
	case domain.ReminderType3Hours:
		log.Printf("[Scheduler] Sending H-3hours reminder to %s for webinar %s", phone, reminder.WebinarID)
		err = s.waService.SendReminder3Hours(phone, reminder.Webinar, userName)
		
	case domain.ReminderType30Min:
		log.Printf("[Scheduler] Sending H-30min reminder to %s for webinar %s", phone, reminder.WebinarID)
		err = s.waService.SendReminder30Min(phone, reminder.Webinar, userName)
		
	default:
		log.Printf("[Scheduler] Unknown reminder type: %s", reminder.ReminderType)
		s.webinarRepo.MarkReminderFailed(reminder.ID, "Unknown reminder type")
		return
	}

	// Update reminder status
	if err != nil {
		log.Printf("[Scheduler] Failed to send reminder %s: %v", reminder.ID, err)
		s.webinarRepo.MarkReminderFailed(reminder.ID, err.Error())
	} else {
		log.Printf("[Scheduler] Reminder %s sent successfully", reminder.ID)
		s.webinarRepo.MarkReminderSent(reminder.ID)
	}
}

// UpdateWebinarStatuses updates webinar statuses based on time
func (s *ReminderScheduler) UpdateWebinarStatuses() {
	now := time.Now()
	
	// Update webinars that should be live (within duration window)
	_, err := s.db.Exec(`
		UPDATE webinars 
		SET status = 'live', updated_at = NOW()
		WHERE status = 'upcoming' 
			AND scheduled_at <= $1 
			AND scheduled_at + (duration_minutes || ' minutes')::interval > $1
	`, now)
	if err != nil {
		log.Printf("[Scheduler] Error updating live webinars: %v", err)
	}
	
	// Update webinars that have completed
	_, err = s.db.Exec(`
		UPDATE webinars 
		SET status = 'completed', updated_at = NOW()
		WHERE status IN ('upcoming', 'live') 
			AND scheduled_at + (duration_minutes || ' minutes')::interval <= $1
	`, now)
	if err != nil {
		log.Printf("[Scheduler] Error updating completed webinars: %v", err)
	}
}

// ===== Singleton for global access =====

var defaultScheduler *ReminderScheduler

// InitScheduler initializes the global scheduler with the given database
func InitScheduler(db *sqlx.DB) {
	if defaultScheduler != nil {
		return // Already initialized
	}
	defaultScheduler = NewReminderScheduler(db)
}

// StartScheduler starts the global scheduler
func StartScheduler() {
	if defaultScheduler == nil {
		log.Println("[Scheduler] Scheduler not initialized")
		return
	}
	defaultScheduler.Start()
}

// StopScheduler stops the global scheduler
func StopScheduler() {
	if defaultScheduler != nil {
		defaultScheduler.Stop()
	}
}

// GetScheduler returns the global scheduler instance
func GetScheduler() *ReminderScheduler {
	return defaultScheduler
}
