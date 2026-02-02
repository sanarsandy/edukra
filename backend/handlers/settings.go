package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/services"
)

// Settings represents platform settings
type Settings struct {
	SiteName        string `json:"site_name"`
	SiteDescription string `json:"site_description"`
	ContactEmail    string `json:"contact_email"`
	Currency        string `json:"currency"`
	Language        string `json:"language"`
	Logo            string `json:"logo_url"`
	Theme           string `json:"theme"`
	BannerEnabled   bool   `json:"banner_enabled"`
	BannerText      string `json:"banner_text"`
	BannerLink      string `json:"banner_link"`
	BannerBgColor   string `json:"banner_bg_color"`
	BannerTextColor string `json:"banner_text_color"`
	
	// Security Settings
	Require2FA     bool `json:"require_2fa"`
	SessionTimeout int  `json:"session_timeout"`
	
	// Feature Flags (JSON string)
	Features string `json:"features"` 
	
	// SMTP Settings
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFromEmail string `json:"smtp_from_email"`
	SMTPFromName  string `json:"smtp_from_name"`
}

// UpdateSettingsRequest uses pointers to allow partial updates (nil = not provided)
type UpdateSettingsRequest struct {
	SiteName        *string `json:"site_name"`
	SiteDescription *string `json:"site_description"`
	ContactEmail    *string `json:"contact_email"`
	Currency        *string `json:"currency"`
	Language        *string `json:"language"`
	Logo            *string `json:"logo_url"`
	Theme           *string `json:"theme"`
	
	BannerEnabled   *bool   `json:"banner_enabled"`
	BannerText      *string `json:"banner_text"`
	BannerLink      *string `json:"banner_link"`
	BannerBgColor   *string `json:"banner_bg_color"`
	BannerTextColor *string `json:"banner_text_color"`
	
	Require2FA     *bool   `json:"require_2fa"`
	SessionTimeout *int    `json:"session_timeout"`
	Features       *string `json:"features"`
	
	SMTPHost      *string `json:"smtp_host"`
	SMTPPort      *int    `json:"smtp_port"`
	SMTPUsername  *string `json:"smtp_username"`
	SMTPPassword  *string `json:"smtp_password"`
	SMTPFromEmail *string `json:"smtp_from_email"`
	SMTPFromName  *string `json:"smtp_from_name"`
}

// Default settings (used as fallback)
var defaultSettings = Settings{
	SiteName:        "LearnHub",
	SiteDescription: "Platform pembelajaran online",
	ContactEmail:    "admin@learnhub.id",
	Currency:        "IDR",
	Language:        "id",
	Logo:            "",
	Theme:           "default",
	BannerEnabled:   false,
	BannerText:      "",
	BannerLink:      "",
	BannerBgColor:   "#1E3A5F",
	BannerTextColor: "#FFFFFF",
	Require2FA:      false,
	SessionTimeout:  30,
	Features:        "[]",
	SMTPHost:        "",
	SMTPPort:        587,
	SMTPUsername:    "",
	SMTPPassword:    "",
	SMTPFromEmail:   "noreply@learnhub.id",
	SMTPFromName:    "LearnHub",
}

// getSettingValue retrieves a single setting value from database
func getSettingValue(key string, defaultValue string) string {
	var value string
	err := db.DB.QueryRow("SELECT value FROM settings WHERE key = $1", key).Scan(&value)
	if err == sql.ErrNoRows || err != nil {
		return defaultValue
	}
	return value
}

// setSettingValue updates or inserts a setting value in database
func setSettingValue(key, value string) error {
	query := `
		INSERT INTO settings (key, value, updated_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = $3
	`
	_, err := db.DB.Exec(query, key, value, time.Now())
	return err
}

// GetSettings returns platform settings from database
func GetSettings(c echo.Context) error {
	bannerEnabledIdx := getSettingValue("banner_enabled", "false")
	require2FAIdx := getSettingValue("require_2fa", "false")
	sessionTimeoutStr := getSettingValue("session_timeout", "30")
	
	sessionTimeout, _ := strconv.Atoi(sessionTimeoutStr)
	if sessionTimeout < 1 {
		sessionTimeout = 30
	}

	smtpPortStr := getSettingValue("smtp_port", "587")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	
	settings := Settings{
		SiteName:        getSettingValue("site_name", defaultSettings.SiteName),
		SiteDescription: getSettingValue("site_description", defaultSettings.SiteDescription),
		ContactEmail:    getSettingValue("contact_email", defaultSettings.ContactEmail),
		Currency:        getSettingValue("currency", defaultSettings.Currency),
		Language:        getSettingValue("language", defaultSettings.Language),
		Logo:            getSettingValue("logo_url", defaultSettings.Logo),
		Theme:           getSettingValue("theme", defaultSettings.Theme),
		BannerEnabled:   bannerEnabledIdx == "true",
		BannerText:      getSettingValue("banner_text", defaultSettings.BannerText),
		BannerLink:      getSettingValue("banner_link", defaultSettings.BannerLink),
		BannerBgColor:   getSettingValue("banner_bg_color", defaultSettings.BannerBgColor),
		BannerTextColor: getSettingValue("banner_text_color", defaultSettings.BannerTextColor),
		Require2FA:      require2FAIdx == "true",
		SessionTimeout:  sessionTimeout,
		Features:        getSettingValue("features", defaultSettings.Features),
		SMTPHost:        getSettingValue("smtp_host", defaultSettings.SMTPHost),
		SMTPPort:        smtpPort,
		SMTPUsername:    getSettingValue("smtp_username", defaultSettings.SMTPUsername),
		SMTPPassword:    getSettingValue("smtp_password", defaultSettings.SMTPPassword),
		SMTPFromEmail:   getSettingValue("smtp_from_email", defaultSettings.SMTPFromEmail),
		SMTPFromName:    getSettingValue("smtp_from_name", defaultSettings.SMTPFromName),
	}

	return c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates platform settings in database (admin only)
func UpdateSettings(c echo.Context) error {
	var req UpdateSettingsRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update each setting if provided (not nil)
	if req.SiteName != nil {
		setSettingValue("site_name", *req.SiteName)
	}
	if req.SiteDescription != nil {
		setSettingValue("site_description", *req.SiteDescription)
	}
	if req.ContactEmail != nil {
		setSettingValue("contact_email", *req.ContactEmail)
	}
	if req.Currency != nil {
		setSettingValue("currency", *req.Currency)
	}
	if req.Language != nil {
		setSettingValue("language", *req.Language)
	}
	if req.Theme != nil {
		setSettingValue("theme", *req.Theme)
	}
	if req.Logo != nil {
		setSettingValue("logo_url", *req.Logo)
	}
	
	// Banner settings
	if req.BannerText != nil {
		setSettingValue("banner_text", *req.BannerText)
	}
	if req.BannerLink != nil {
		setSettingValue("banner_link", *req.BannerLink)
	}
	if req.BannerBgColor != nil {
		setSettingValue("banner_bg_color", *req.BannerBgColor)
	}
	if req.BannerTextColor != nil {
		setSettingValue("banner_text_color", *req.BannerTextColor)
	}
	
	if req.BannerEnabled != nil {
		val := "false"
		if *req.BannerEnabled {
			val = "true"
		}
		setSettingValue("banner_enabled", val)
	}
	
	if req.Require2FA != nil {
		val := "false"
		if *req.Require2FA {
			val = "true"
		}
		setSettingValue("require_2fa", val)
	}
	
	// Features (store as JSON string)
	if req.Features != nil {
		setSettingValue("features", *req.Features)
	}
	
	// Session Timeout
	if req.SessionTimeout != nil && *req.SessionTimeout > 0 {
		setSettingValue("session_timeout", strconv.Itoa(*req.SessionTimeout))
	}
	
	// SMTP Settings
	if req.SMTPHost != nil {
		setSettingValue("smtp_host", *req.SMTPHost)
	}
	if req.SMTPPort != nil {
		setSettingValue("smtp_port", strconv.Itoa(*req.SMTPPort))
	}
	if req.SMTPUsername != nil {
		setSettingValue("smtp_username", *req.SMTPUsername)
	}
	if req.SMTPPassword != nil {
		setSettingValue("smtp_password", *req.SMTPPassword)
	}
	if req.SMTPFromEmail != nil {
		setSettingValue("smtp_from_email", *req.SMTPFromEmail)
	}
	if req.SMTPFromName != nil {
		setSettingValue("smtp_from_name", *req.SMTPFromName)
	}

	// Return updated settings
	return GetSettings(c)
}

// TestSMTP sends a test email to the current user to verify SMTP settings
func TestSMTP(c echo.Context) error {
	// Get current user email from token or request
	// For now, let's allow passing email in body or use admin default
	type TestRequest struct {
		ToEmail string `json:"to_email"`
	}
	var req TestRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if req.ToEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	// Use services.Email to send test email
	// Needs import "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/services"
	err := services.Email.SendTestEmail(req.ToEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to send test email: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Test email sent successfully"})
}
