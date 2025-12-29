package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var campaignRepo *postgres.CampaignRepository

func initCampaignRepo() {
	if campaignRepo == nil && db.DB != nil {
		campaignRepo = postgres.NewCampaignRepository(db.DB)
	}
}

// ==================== ADMIN ENDPOINTS ====================

// CreateCampaign creates a new campaign (Admin)
// POST /api/admin/campaigns
func CreateCampaign(c echo.Context) error {
	initCampaignRepo()

	var req domain.CreateCampaignRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate required fields
	if req.Slug == "" || req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Slug and title are required"})
	}

	// Normalize slug
	req.Slug = strings.ToLower(strings.ReplaceAll(req.Slug, " ", "-"))

	// Check slug uniqueness
	exists, _ := campaignRepo.SlugExists(req.Slug, "")
	if exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Slug already exists"})
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, *req.StartDate); err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := time.Parse(time.RFC3339, *req.EndDate); err == nil {
			endDate = &t
		}
	}

	// Get default blocks if none provided
	var blocksJSON []byte
	if req.Blocks != nil && len(req.Blocks) > 0 {
		blocksJSON = req.Blocks
	} else {
		defaultBlocks := domain.GetDefaultBlocks()
		blocksJSON, _ = json.Marshal(defaultBlocks)
	}

	campaign := &domain.Campaign{
		Slug:        req.Slug,
		Title:       req.Title,
		IsActive:    req.IsActive,
		CourseID:    req.CourseID,
		MetaDesc:    req.MetaDesc,
		OGImageURL:  req.OGImageURL,
		Blocks:      blocksJSON,
		Styles:      req.Styles,
		HTMLContent: req.HTMLContent,
		CSSContent:  req.CSSContent,
		GJSData:     req.GJSData,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := campaignRepo.Create(campaign); err != nil {
		log.Printf("[Campaign] Failed to create: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create campaign"})
	}

	return c.JSON(http.StatusCreated, campaign)
}

// ListCampaigns lists all campaigns (Admin)
// GET /api/admin/campaigns
func ListCampaigns(c echo.Context) error {
	initCampaignRepo()

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	campaigns, err := campaignRepo.List(limit, offset)
	if err != nil {
		log.Printf("[Campaign] Failed to list: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list campaigns"})
	}

	if campaigns == nil {
		campaigns = []*domain.Campaign{}
	}

	return c.JSON(http.StatusOK, campaigns)
}

// GetCampaign gets a campaign by ID (Admin)
// GET /api/admin/campaigns/:id
func GetCampaign(c echo.Context) error {
	initCampaignRepo()
	id := c.Param("id")

	campaign, err := campaignRepo.GetByID(id)
	if err != nil {
		log.Printf("[Campaign] Error fetching by ID: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get campaign"})
	}
	if campaign == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Campaign not found"})
	}

	return c.JSON(http.StatusOK, campaign)
}

// UpdateCampaign updates a campaign (Admin)
// PUT /api/admin/campaigns/:id
func UpdateCampaign(c echo.Context) error {
	initCampaignRepo()
	id := c.Param("id")

	campaign, err := campaignRepo.GetByID(id)
	if err != nil || campaign == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Campaign not found"})
	}

	var req domain.UpdateCampaignRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Update fields if provided
	if req.Slug != nil {
		slug := strings.ToLower(strings.ReplaceAll(*req.Slug, " ", "-"))
		// Check uniqueness
		exists, _ := campaignRepo.SlugExists(slug, id)
		if exists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Slug already exists"})
		}
		campaign.Slug = slug
	}
	if req.Title != nil {
		campaign.Title = *req.Title
	}
	if req.IsActive != nil {
		campaign.IsActive = *req.IsActive
	}
	if req.CourseID != nil {
		campaign.CourseID = req.CourseID
	}
	if req.MetaDesc != nil {
		campaign.MetaDesc = req.MetaDesc
	}
	if req.OGImageURL != nil {
		campaign.OGImageURL = req.OGImageURL
	}
	if req.Blocks != nil && len(req.Blocks) > 0 {
		campaign.Blocks = req.Blocks
	}
	if req.Styles != nil && len(req.Styles) > 0 {
		campaign.Styles = req.Styles
	}
	if req.HTMLContent != nil {
		campaign.HTMLContent = req.HTMLContent
	}
	if req.CSSContent != nil {
		campaign.CSSContent = req.CSSContent
	}
	if req.GJSData != nil && len(req.GJSData) > 0 {
		campaign.GJSData = req.GJSData
	}
	if req.StartDate != nil {
		if *req.StartDate == "" {
			campaign.StartDate = nil
		} else if t, err := time.Parse(time.RFC3339, *req.StartDate); err == nil {
			campaign.StartDate = &t
		}
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			campaign.EndDate = nil
		} else if t, err := time.Parse(time.RFC3339, *req.EndDate); err == nil {
			campaign.EndDate = &t
		}
	}

	if err := campaignRepo.Update(campaign); err != nil {
		log.Printf("[Campaign] Failed to update: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update campaign"})
	}

	return c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign deletes a campaign (Admin)
// DELETE /api/admin/campaigns/:id
func DeleteCampaign(c echo.Context) error {
	initCampaignRepo()
	id := c.Param("id")

	if err := campaignRepo.Delete(id); err != nil {
		log.Printf("[Campaign] Failed to delete: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete campaign"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// GetCampaignAnalytics gets analytics for a campaign (Admin)
// GET /api/admin/campaigns/:id/analytics
func GetCampaignAnalytics(c echo.Context) error {
	initCampaignRepo()
	id := c.Param("id")

	summary, err := campaignRepo.GetAnalyticsSummary(id)
	if err != nil {
		log.Printf("[Campaign] Failed to get analytics: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get analytics"})
	}

	return c.JSON(http.StatusOK, summary)
}

// ==================== PUBLIC ENDPOINTS ====================

// GetCampaignBySlug gets an active campaign by slug (Public)
// GET /api/c/:slug
func GetCampaignBySlug(c echo.Context) error {
	initCampaignRepo()
	slug := c.Param("slug")

	campaign, err := campaignRepo.GetBySlug(slug)
	if err != nil {
		log.Printf("[Campaign] Error fetching slug %s: %v", slug, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch campaign"})
	}
	if campaign == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Campaign not found"})
	}

	// Track view asynchronously
	go func() {
		// Increment counter
		if err := campaignRepo.IncrementViewCount(campaign.ID); err != nil {
			log.Printf("[Campaign] Failed to increment view count: %v", err)
		}

		// Track detailed event
		event := &domain.CampaignAnalytics{
			CampaignID: campaign.ID,
			EventType:  "view",
		}
		// Get visitor info from context
		if ip := c.RealIP(); ip != "" {
			event.VisitorIP = &ip
		}
		if ua := c.Request().Header.Get("User-Agent"); ua != "" {
			event.UserAgent = &ua
		}
		if ref := c.Request().Header.Get("Referer"); ref != "" {
			event.Referer = &ref
		}
		// UTM params
		if utm := c.QueryParam("utm_source"); utm != "" {
			event.UTMSource = &utm
		}
		if utm := c.QueryParam("utm_medium"); utm != "" {
			event.UTMMedium = &utm
		}
		if utm := c.QueryParam("utm_campaign"); utm != "" {
			event.UTMCampaign = &utm
		}

		if err := campaignRepo.TrackEvent(event); err != nil {
			log.Printf("[Campaign] Failed to track view event: %v", err)
		}
	}()

	return c.JSON(http.StatusOK, campaign)
}

// TrackCampaignClick tracks a click on the CTA (Public)
// POST /api/c/:id/click
func TrackCampaignClick(c echo.Context) error {
	initCampaignRepo()
	id := c.Param("id")

	go func() {
		// Increment counter
		if err := campaignRepo.IncrementClickCount(id); err != nil {
			log.Printf("[Campaign] Failed to increment click count: %v", err)
		}

		// Track detailed event
		event := &domain.CampaignAnalytics{
			CampaignID: id,
			EventType:  "click",
		}
		if ip := c.RealIP(); ip != "" {
			event.VisitorIP = &ip
		}
		if ua := c.Request().Header.Get("User-Agent"); ua != "" {
			event.UserAgent = &ua
		}

		if err := campaignRepo.TrackEvent(event); err != nil {
			log.Printf("[Campaign] Failed to track click event: %v", err)
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// TrackCampaignConversion tracks a purchase conversion (called from payment webhook)
// This is an internal function, not an HTTP endpoint
func TrackCampaignConversion(campaignID string, transactionID string) {
	initCampaignRepo()

	// Increment counter
	if err := campaignRepo.IncrementConversionCount(campaignID); err != nil {
		log.Printf("[Campaign] Failed to increment conversion count: %v", err)
	}

	// Track detailed event
	event := &domain.CampaignAnalytics{
		CampaignID:    campaignID,
		EventType:     "conversion",
		TransactionID: &transactionID,
	}

	if err := campaignRepo.TrackEvent(event); err != nil {
		log.Printf("[Campaign] Failed to track conversion event: %v", err)
	}
}
