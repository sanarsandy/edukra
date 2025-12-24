package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

var platformSettings = Settings{
	SiteName:        "LearnHub",
	SiteDescription: "Platform pembelajaran online",
	ContactEmail:    "admin@learnhub.id",
	Currency:        "IDR",
	Language:        "id",
	Logo:            "",
	Theme:           "default",
}

// GetSettings returns platform settings
func GetSettings(c echo.Context) error {
	return c.JSON(http.StatusOK, platformSettings)
}

// UpdateSettings updates platform settings (admin only)
func UpdateSettings(c echo.Context) error {
	var req Settings
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	// Update settings
	if req.SiteName != "" {
		platformSettings.SiteName = req.SiteName
	}
	if req.SiteDescription != "" {
		platformSettings.SiteDescription = req.SiteDescription
	}
	if req.ContactEmail != "" {
		platformSettings.ContactEmail = req.ContactEmail
	}
	if req.Currency != "" {
		platformSettings.Currency = req.Currency
	}
	if req.Language != "" {
		platformSettings.Language = req.Language
	}
	if req.Theme != "" {
		platformSettings.Theme = req.Theme
	}
	
	return c.JSON(http.StatusOK, platformSettings)
}
