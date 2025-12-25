package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
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
	settings := Settings{
		SiteName:        getSettingValue("site_name", defaultSettings.SiteName),
		SiteDescription: getSettingValue("site_description", defaultSettings.SiteDescription),
		ContactEmail:    getSettingValue("contact_email", defaultSettings.ContactEmail),
		Currency:        getSettingValue("currency", defaultSettings.Currency),
		Language:        getSettingValue("language", defaultSettings.Language),
		Logo:            getSettingValue("logo_url", defaultSettings.Logo),
		Theme:           getSettingValue("theme", defaultSettings.Theme),
	}

	return c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates platform settings in database (admin only)
func UpdateSettings(c echo.Context) error {
	var req Settings

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update each setting if provided
	if req.SiteName != "" {
		if err := setSettingValue("site_name", req.SiteName); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save site_name"})
		}
	}
	if req.SiteDescription != "" {
		if err := setSettingValue("site_description", req.SiteDescription); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save site_description"})
		}
	}
	if req.ContactEmail != "" {
		if err := setSettingValue("contact_email", req.ContactEmail); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save contact_email"})
		}
	}
	if req.Currency != "" {
		if err := setSettingValue("currency", req.Currency); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save currency"})
		}
	}
	if req.Language != "" {
		if err := setSettingValue("language", req.Language); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save language"})
		}
	}
	if req.Theme != "" {
		if err := setSettingValue("theme", req.Theme); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save theme"})
		}
	}
	if req.Logo != "" {
		if err := setSettingValue("logo_url", req.Logo); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save logo_url"})
		}
	}

	// Return updated settings
	settings := Settings{
		SiteName:        getSettingValue("site_name", defaultSettings.SiteName),
		SiteDescription: getSettingValue("site_description", defaultSettings.SiteDescription),
		ContactEmail:    getSettingValue("contact_email", defaultSettings.ContactEmail),
		Currency:        getSettingValue("currency", defaultSettings.Currency),
		Language:        getSettingValue("language", defaultSettings.Language),
		Logo:            getSettingValue("logo_url", defaultSettings.Logo),
		Theme:           getSettingValue("theme", defaultSettings.Theme),
	}

	return c.JSON(http.StatusOK, settings)
}
