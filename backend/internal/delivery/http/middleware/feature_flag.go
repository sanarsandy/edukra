package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// FeatureFlag represents available features
type FeatureFlag string

const (
	FeatureQuiz        FeatureFlag = "enable_quiz"
	FeatureCertificate FeatureFlag = "enable_certificate"
	FeatureForum       FeatureFlag = "enable_forum"
)

// FeatureConfig holds the feature flags configuration
type FeatureConfig struct {
	EnableQuiz         bool   `json:"enable_quiz"`
	EnableCertificate  bool   `json:"enable_certificate"`
	EnableForum        bool   `json:"enable_forum"`
	DRMProtectionLevel string `json:"drm_protection_level"`
}

// FeatureResolver defines how to get feature config for a tenant
type FeatureResolver interface {
	GetFeatures(tenantID string) (*FeatureConfig, error)
}

// RequireFeature creates a middleware that checks if a feature is enabled
func RequireFeature(resolver FeatureResolver, feature FeatureFlag) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := GetTenantID(c)
			if tenantID == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Tenant identification required",
				})
			}

			features, err := resolver.GetFeatures(tenantID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to get feature configuration",
				})
			}

			enabled := false
			switch feature {
			case FeatureQuiz:
				enabled = features.EnableQuiz
			case FeatureCertificate:
				enabled = features.EnableCertificate
			case FeatureForum:
				enabled = features.EnableForum
			}

			if !enabled {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "This feature is not available for your plan",
				})
			}

			return next(c)
		}
	}
}

// FeatureMiddleware injects feature config into context
func FeatureMiddleware(resolver FeatureResolver) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := GetTenantID(c)
			if tenantID != "" {
				features, err := resolver.GetFeatures(tenantID)
				if err == nil {
					c.Set("features", features)
				}
			}
			return next(c)
		}
	}
}

// GetFeatures extracts feature config from context
func GetFeatures(c echo.Context) *FeatureConfig {
	features := c.Get("features")
	if features == nil {
		return nil
	}
	return features.(*FeatureConfig)
}

// IsFeatureEnabled checks if a feature is enabled in the current context
func IsFeatureEnabled(c echo.Context, feature FeatureFlag) bool {
	features := GetFeatures(c)
	if features == nil {
		return false
	}

	switch feature {
	case FeatureQuiz:
		return features.EnableQuiz
	case FeatureCertificate:
		return features.EnableCertificate
	case FeatureForum:
		return features.EnableForum
	default:
		return false
	}
}
