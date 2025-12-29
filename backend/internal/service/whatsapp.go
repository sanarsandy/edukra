package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
)

// WhatsAppService handles WhatsApp message sending
type WhatsAppService struct {
	apiURL    string
	apiKey    string
	client    *http.Client
	isEnabled bool
}

// NewWhatsAppService creates a new WhatsApp service instance
func NewWhatsAppService() *WhatsAppService {
	apiURL := os.Getenv("WA_GATEWAY_URL")
	apiKey := os.Getenv("WA_GATEWAY_API_KEY")

	if apiURL == "" {
		apiURL = "https://lman.id/akila-services/wa-gateway/whatsapp/"
	}

	return &WhatsAppService{
		apiURL:    apiURL,
		apiKey:    apiKey,
		client:    &http.Client{Timeout: 30 * time.Second},
		isEnabled: apiKey != "",
	}
}

// IsEnabled returns whether WhatsApp service is configured
func (s *WhatsAppService) IsEnabled() bool {
	return s.isEnabled
}

// SendMessage sends a WhatsApp message to the specified phone number
func (s *WhatsAppService) SendMessage(phone, message string) error {
	if !s.isEnabled {
		log.Printf("[WhatsApp] Service not enabled, skipping message to %s", phone)
		return nil
	}

	// Normalize phone number
	phone = normalizePhoneNumber(phone)

	// Prepare form data
	data := url.Values{}
	data.Set("api_key", s.apiKey)
	data.Set("tujuan", phone)
	data.Set("pesan", message)

	// Create request
	req, err := http.NewRequest("POST", s.apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WA Gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[WhatsApp] Message sent to %s (response: %s)", phone, string(body))
	return nil
}

// normalizePhoneNumber converts phone to format expected by WA Gateway
func normalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Remove + prefix if present
	phone = strings.TrimPrefix(phone, "+")

	// Convert 08xxx to 628xxx
	if strings.HasPrefix(phone, "08") {
		phone = "62" + phone[1:]
	}

	// Add 62 if no country code
	if !strings.HasPrefix(phone, "62") {
		phone = "62" + phone
	}

	return phone
}

// WebinarConfirmationData holds data for webinar confirmation message
type WebinarConfirmationData struct {
	UserName        string
	UserEmail       string
	WebinarTitle    string
	WebinarDate     string
	WebinarTime     string
	DurationMinutes int
	MeetingURL      string
	MeetingPassword string
	IsNewUser       bool
	TempPassword    string
	LMSUrl          string
}

// SendWebinarConfirmation sends webinar registration confirmation
func (s *WhatsAppService) SendWebinarConfirmation(phone string, data WebinarConfirmationData) error {
	var passwordInfo string
	if data.IsNewUser && data.TempPassword != "" {
		passwordInfo = fmt.Sprintf("\n🔑 *Password Sementara:* %s\n(Silakan ganti password setelah login)", data.TempPassword)
	}

	var meetingPasswordInfo string
	if data.MeetingPassword != "" {
		meetingPasswordInfo = fmt.Sprintf("\n🔐 Password Meeting: %s", data.MeetingPassword)
	}

	message := fmt.Sprintf(`🎉 *Selamat! Pendaftaran Berhasil*

Halo %s! 👋

Terima kasih telah mendaftar webinar:
📚 *%s*

📅 *Jadwal:*
Tanggal: %s
Waktu: %s WIB
Durasi: %d menit

🔐 *Akses LMS Anda:*
Email: %s%s
Link: %s

📍 *Link Webinar:*
%s%s

Sampai jumpa di webinar! 🚀

---
EDUKRA Learning Platform`,
		data.UserName,
		data.WebinarTitle,
		data.WebinarDate,
		data.WebinarTime,
		data.DurationMinutes,
		data.UserEmail,
		passwordInfo,
		data.LMSUrl,
		data.MeetingURL,
		meetingPasswordInfo,
	)

	return s.SendMessage(phone, message)
}

// SendReminder1Day sends H-1 reminder
func (s *WhatsAppService) SendReminder1Day(phone string, webinar *domain.Webinar, userName string) error {
	message := fmt.Sprintf(`⏰ *Reminder: Webinar Besok!*

Halo %s! 👋

Jangan lupa, besok ada webinar:
📚 *%s*

📅 %s | %s WIB

Pastikan Anda sudah:
✅ Menyiapkan koneksi internet stabil
✅ Menyiapkan headset/speaker
✅ Menyimpan link webinar

See you tomorrow! 🎯

---
EDUKRA Learning Platform`,
		userName,
		webinar.Title,
		webinar.ScheduledAt.Format("02 January 2006"),
		webinar.ScheduledAt.Format("15:04"),
	)

	return s.SendMessage(phone, message)
}

// SendReminder3Hours sends H-3 hours reminder
func (s *WhatsAppService) SendReminder3Hours(phone string, webinar *domain.Webinar, userName string) error {
	var meetingInfo string
	if webinar.MeetingURL != nil && *webinar.MeetingURL != "" {
		meetingInfo = fmt.Sprintf("\n📍 Link: %s", *webinar.MeetingURL)
		if webinar.MeetingPassword != nil && *webinar.MeetingPassword != "" {
			meetingInfo += fmt.Sprintf("\n🔑 Password: %s", *webinar.MeetingPassword)
		}
	}

	message := fmt.Sprintf(`🔔 *3 Jam Lagi!*

Halo %s! 👋

Webinar *%s* dimulai dalam 3 jam!

⏰ %s WIB%s

Siapkan diri Anda! 💪

---
EDUKRA Learning Platform`,
		userName,
		webinar.Title,
		webinar.ScheduledAt.Format("15:04"),
		meetingInfo,
	)

	return s.SendMessage(phone, message)
}

// SendReminder30Min sends H-30 min reminder with meeting link
func (s *WhatsAppService) SendReminder30Min(phone string, webinar *domain.Webinar, userName string) error {
	var meetingURL string
	var meetingPassword string
	if webinar.MeetingURL != nil {
		meetingURL = *webinar.MeetingURL
	}
	if webinar.MeetingPassword != nil {
		meetingPassword = *webinar.MeetingPassword
	}

	var passwordInfo string
	if meetingPassword != "" {
		passwordInfo = fmt.Sprintf("\n🔑 Password: %s", meetingPassword)
	}

	message := fmt.Sprintf(`🚀 *30 Menit Lagi - SIAP JOIN!*

Halo %s! 👋

Webinar dimulai sebentar lagi!

📚 *%s*
⏰ %s WIB

👇 *KLIK UNTUK JOIN:*
%s%s

Kami tunggu kehadirannya! 🎬

---
EDUKRA Learning Platform`,
		userName,
		webinar.Title,
		webinar.ScheduledAt.Format("15:04"),
		meetingURL,
		passwordInfo,
	)

	return s.SendMessage(phone, message)
}

// SendPaymentSuccess sends payment success notification (for non-webinar courses)
func (s *WhatsAppService) SendPaymentSuccess(phone, userName, courseName, lmsURL string) error {
	message := fmt.Sprintf(`✅ *Pembayaran Berhasil!*

Halo %s! 👋

Pembayaran Anda untuk kursus berikut telah berhasil:
📚 *%s*

🎓 Anda sekarang bisa mengakses semua materi kursus di:
🔗 %s

Selamat belajar! 🚀

---
EDUKRA Learning Platform`,
		userName,
		courseName,
		lmsURL,
	)

	return s.SendMessage(phone, message)
}

// SendCredentials sends new user credentials
func (s *WhatsAppService) SendCredentials(phone, userName, email, password, lmsURL string) error {
	message := fmt.Sprintf(`🔐 *Akun EDUKRA Anda*

Halo %s! 👋

Akun EDUKRA Anda telah dibuat:

📧 Email: %s
🔑 Password: %s

Login di: %s

⚠️ Untuk keamanan, segera ganti password Anda setelah login.

Selamat belajar! 🎓

---
EDUKRA Learning Platform`,
		userName,
		email,
		password,
		lmsURL,
	)

	return s.SendMessage(phone, message)
}

// LogNotification logs WA notification to database
func LogNotification(db interface{}, userID *string, phone, messageType string, data interface{}, status string, errMsg *string) {
	// This would be implemented to log to wa_notifications table
	// For now, just log to console
	dataJSON, _ := json.Marshal(data)
	log.Printf("[WhatsApp] Notification logged: type=%s, phone=%s, status=%s, data=%s", 
		messageType, phone, status, string(dataJSON))
}

// WANotificationLog for structured logging
type WANotificationLog struct {
	UserID       *string
	Phone        string
	MessageType  string
	TemplateData []byte
	Status       string
	SentAt       *time.Time
	ErrorMessage *string
}

// SaveNotificationLog saves notification to database
func SaveNotificationLog(db interface{}, log *WANotificationLog) error {
	// Implementation would insert into wa_notifications table
	return nil
}

// ===== Singleton for global access =====

var defaultWAService *WhatsAppService

// InitWhatsAppService initializes the global WA service
func InitWhatsAppService() {
	defaultWAService = NewWhatsAppService()
	if defaultWAService.IsEnabled() {
		log.Printf("[WhatsApp] Service initialized with API URL: %s", defaultWAService.apiURL)
	} else {
		log.Println("[WhatsApp] Service not configured (WA_GATEWAY_API_KEY not set)")
	}
}

// GetWhatsAppService returns the global WA service instance
func GetWhatsAppService() *WhatsAppService {
	if defaultWAService == nil {
		InitWhatsAppService()
	}
	return defaultWAService
}

// SendWebinarConfirmationAsync sends confirmation in background goroutine
func SendWebinarConfirmationAsync(phone string, data WebinarConfirmationData) {
	go func() {
		if err := GetWhatsAppService().SendWebinarConfirmation(phone, data); err != nil {
			log.Printf("[WhatsApp] Failed to send webinar confirmation to %s: %v", phone, err)
		}
	}()
}

// Buffer for batch processing if needed
var messageBuffer *bytes.Buffer

func init() {
	messageBuffer = new(bytes.Buffer)
}
