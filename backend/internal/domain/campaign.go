package domain

import (
	"encoding/json"
	"time"
)

// Campaign represents a promotional landing page
type Campaign struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CourseID    *string   `json:"course_id,omitempty" db:"course_id"`
	Title       string    `json:"title" db:"title"`
	MetaDesc    *string   `json:"meta_description,omitempty" db:"meta_description"`
	OGImageURL  *string   `json:"og_image_url,omitempty" db:"og_image_url"`
	Blocks      json.RawMessage `json:"blocks" db:"blocks"` // JSONB (legacy)
	Styles      json.RawMessage `json:"styles,omitempty" db:"styles"` // JSONB (legacy)
	HTMLContent *string   `json:"html_content,omitempty" db:"html_content"` // GrapeJS HTML output
	CSSContent  *string   `json:"css_content,omitempty" db:"css_content"` // GrapeJS CSS output
	GJSData     json.RawMessage `json:"gjs_data,omitempty" db:"gjs_data"` // GrapeJS project data
	StartDate   *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	ViewCount   int       `json:"view_count" db:"view_count"`
	ClickCount  int       `json:"click_count" db:"click_count"`
	ConversionCount int   `json:"conversion_count" db:"conversion_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joined data (not in DB)
	Course      *Course   `json:"course,omitempty" db:"-"`
	ParsedBlocks []CampaignBlock `json:"parsed_blocks,omitempty" db:"-"`
}

// CampaignBlock represents a single content block in the landing page
type CampaignBlock struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"` // hero, countdown, benefits, pricing, testimonials, faq, instructor
	Enabled bool            `json:"enabled"`
	Order   int             `json:"order"`
	Data    json.RawMessage `json:"data"` // Block-specific data
}

// Block Data Types

type HeroBlockData struct {
	Headline        string  `json:"headline"`
	Subheadline     string  `json:"subheadline"`
	BackgroundImage *string `json:"background_image,omitempty"`
	VideoURL        *string `json:"video_url,omitempty"`
	CTAText         string  `json:"cta_text"`
	CTALink         string  `json:"cta_link"`
}

type CountdownBlockData struct {
	EndDate string `json:"end_date"` // ISO string
	Label   string `json:"label"`
}

type BenefitsBlockData struct {
	Title string        `json:"title"`
	Items []BenefitItem `json:"items"`
}

type BenefitItem struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

type PricingBlockData struct {
	OriginalPrice float64 `json:"original_price"`
	DiscountPrice float64 `json:"discount_price"`
	Currency      string  `json:"currency"`
	CTAText       string  `json:"cta_text"`
	ShowTimer     bool    `json:"show_timer"`
}

type TestimonialsBlockData struct {
	Title string            `json:"title"`
	Items []TestimonialItem `json:"items"`
}

type TestimonialItem struct {
	Name   string  `json:"name"`
	Text   string  `json:"text"`
	Avatar *string `json:"avatar,omitempty"`
	Role   *string `json:"role,omitempty"`
}

type FAQBlockData struct {
	Title string    `json:"title"`
	Items []FAQItem `json:"items"`
}

type FAQItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type InstructorBlockData struct {
	AutoFill bool    `json:"auto_fill"` // If true, use course instructor
	Name     *string `json:"name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// CampaignAnalytics tracks individual events
type CampaignAnalytics struct {
	ID            string    `json:"id" db:"id"`
	CampaignID    string    `json:"campaign_id" db:"campaign_id"`
	EventType     string    `json:"event_type" db:"event_type"` // view, click, conversion
	TransactionID *string   `json:"transaction_id,omitempty" db:"transaction_id"`
	VisitorIP     *string   `json:"visitor_ip,omitempty" db:"visitor_ip"`
	UserAgent     *string   `json:"user_agent,omitempty" db:"user_agent"`
	Referer       *string   `json:"referer,omitempty" db:"referer"`
	UTMSource     *string   `json:"utm_source,omitempty" db:"utm_source"`
	UTMMedium     *string   `json:"utm_medium,omitempty" db:"utm_medium"`
	UTMCampaign   *string   `json:"utm_campaign,omitempty" db:"utm_campaign"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// Request/Response structs

type CreateCampaignRequest struct {
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	CourseID    *string         `json:"course_id,omitempty"`
	MetaDesc    *string         `json:"meta_description,omitempty"`
	OGImageURL  *string         `json:"og_image_url,omitempty"`
	Blocks      json.RawMessage `json:"blocks,omitempty"`
	Styles      json.RawMessage `json:"styles,omitempty"`
	HTMLContent *string         `json:"html_content,omitempty"`
	CSSContent  *string         `json:"css_content,omitempty"`
	GJSData     json.RawMessage `json:"gjs_data,omitempty"`
	StartDate   *string         `json:"start_date,omitempty"`
	EndDate     *string         `json:"end_date,omitempty"`
	IsActive    bool            `json:"is_active"`
}

type UpdateCampaignRequest struct {
	Slug        *string         `json:"slug,omitempty"`
	Title       *string         `json:"title,omitempty"`
	CourseID    *string         `json:"course_id,omitempty"`
	MetaDesc    *string         `json:"meta_description,omitempty"`
	OGImageURL  *string         `json:"og_image_url,omitempty"`
	Blocks      json.RawMessage `json:"blocks,omitempty"`
	Styles      json.RawMessage `json:"styles,omitempty"`
	HTMLContent *string         `json:"html_content,omitempty"`
	CSSContent  *string         `json:"css_content,omitempty"`
	GJSData     json.RawMessage `json:"gjs_data,omitempty"`
	StartDate   *string         `json:"start_date,omitempty"`
	EndDate     *string         `json:"end_date,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

// Default blocks for new campaign
func GetDefaultBlocks() []CampaignBlock {
	return []CampaignBlock{
		{ID: "hero_1", Type: "hero", Enabled: true, Order: 1, Data: json.RawMessage(`{"headline":"Judul Kursus Anda","subheadline":"Deskripsi singkat yang menarik","cta_text":"Daftar Sekarang","cta_link":"#pricing"}`)},
		{ID: "countdown_1", Type: "countdown", Enabled: true, Order: 2, Data: json.RawMessage(`{"label":"Promo Berakhir Dalam"}`)},
		{ID: "benefits_1", Type: "benefits", Enabled: true, Order: 3, Data: json.RawMessage(`{"title":"Apa yang Akan Anda Dapatkan?","items":[{"icon":"check","text":"Akses materi selamanya"},{"icon":"video","text":"Video berkualitas tinggi"},{"icon":"certificate","text":"Sertifikat kelulusan"}]}`)},
		{ID: "pricing_1", Type: "pricing", Enabled: true, Order: 4, Data: json.RawMessage(`{"original_price":500000,"discount_price":199000,"currency":"IDR","cta_text":"Beli Sekarang","show_timer":true}`)},
		{ID: "testimonials_1", Type: "testimonials", Enabled: false, Order: 5, Data: json.RawMessage(`{"title":"Kata Mereka","items":[]}`)},
		{ID: "faq_1", Type: "faq", Enabled: false, Order: 6, Data: json.RawMessage(`{"title":"Pertanyaan Umum","items":[]}`)},
		{ID: "instructor_1", Type: "instructor", Enabled: false, Order: 7, Data: json.RawMessage(`{"auto_fill":true}`)},
	}
}
