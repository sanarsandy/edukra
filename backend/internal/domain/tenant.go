package domain

import (
	"time"
)

// Tenant represents a whitelabel client with their own branding and features
type Tenant struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Subdomain     *string        `json:"subdomain,omitempty"`
	CustomDomain  *string        `json:"custom_domain,omitempty"`
	UIConfig      UIConfig       `json:"ui_config"`
	FeatureConfig FeatureConfig  `json:"feature_config"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// UIConfig holds the whitelabel UI configuration
type UIConfig struct {
	PrimaryColor   string  `json:"primary_color"`
	SecondaryColor string  `json:"secondary_color"`
	LogoURL        *string `json:"logo_url,omitempty"`
	FaviconURL     *string `json:"favicon_url,omitempty"`
	FontFamily     string  `json:"font_family"`
}

// FeatureConfig holds the feature flags for a tenant
type FeatureConfig struct {
	EnableQuiz          bool   `json:"enable_quiz"`
	EnableCertificate   bool   `json:"enable_certificate"`
	EnableForum         bool   `json:"enable_forum"`
	DRMProtectionLevel  string `json:"drm_protection_level"` // basic, standard, advanced
}

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetByID(id string) (*Tenant, error)
	GetBySubdomain(subdomain string) (*Tenant, error)
	GetByCustomDomain(domain string) (*Tenant, error)
	Create(tenant *Tenant) error
	Update(tenant *Tenant) error
	Delete(id string) error
	List(limit, offset int) ([]*Tenant, error)
}

// TenantService defines the interface for tenant business logic
type TenantService interface {
	GetTenant(id string) (*Tenant, error)
	ResolveTenant(host string) (*Tenant, error)
	CreateTenant(tenant *Tenant) error
	UpdateTenant(tenant *Tenant) error
	DeleteTenant(id string) error
	ListTenants(limit, offset int) ([]*Tenant, error)
}
