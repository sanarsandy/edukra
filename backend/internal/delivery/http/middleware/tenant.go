package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	TenantIDKey     = "tenant_id"
	TenantConfigKey = "tenant_config"
)

// TenantContext holds tenant information for the request
type TenantContext struct {
	TenantID     string `json:"tenant_id"`
	Subdomain    string `json:"subdomain"`
	CustomDomain string `json:"custom_domain"`
}

// TenantResolver defines how to resolve tenant from host
type TenantResolver interface {
	ResolveBySubdomain(subdomain string) (*TenantContext, error)
	ResolveByCustomDomain(domain string) (*TenantContext, error)
}

// TenantMiddleware creates a middleware that detects tenant from the request
func TenantMiddleware(resolver TenantResolver, baseDomain string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := c.Request().Host

			// Remove port if present
			if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
				host = host[:colonIdx]
			}

			var tenant *TenantContext
			var err error

			// Check if it's a subdomain of baseDomain
			if baseDomain != "" && strings.HasSuffix(host, baseDomain) {
				subdomain := strings.TrimSuffix(host, "."+baseDomain)
				if subdomain != "" && subdomain != host {
					tenant, err = resolver.ResolveBySubdomain(subdomain)
				}
			}

			// If not resolved via subdomain, try custom domain
			if tenant == nil && err == nil {
				tenant, err = resolver.ResolveByCustomDomain(host)
			}

			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": "Tenant not found",
				})
			}

			// If still no tenant resolved, allow request (for main domain)
			if tenant != nil {
				c.Set(TenantIDKey, tenant.TenantID)
				c.Set(TenantConfigKey, tenant)
			}

			return next(c)
		}
	}
}

// RequireTenant ensures a tenant is present in the context
func RequireTenant() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := c.Get(TenantIDKey)
			if tenantID == nil || tenantID.(string) == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Tenant identification required",
				})
			}
			return next(c)
		}
	}
}

// GetTenantID extracts tenant ID from context
func GetTenantID(c echo.Context) string {
	tenantID := c.Get(TenantIDKey)
	if tenantID == nil {
		return ""
	}
	return tenantID.(string)
}

// GetTenantContext extracts full tenant context from request
func GetTenantContext(c echo.Context) *TenantContext {
	tenant := c.Get(TenantConfigKey)
	if tenant == nil {
		return nil
	}
	return tenant.(*TenantContext)
}
