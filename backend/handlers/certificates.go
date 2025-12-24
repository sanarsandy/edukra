package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var certificateRepo *postgres.CertificateRepository

func initCertificateRepos() {
	if certificateRepo == nil && db.DB != nil {
		certificateRepo = postgres.NewCertificateRepository(db.DB)
	}
	initCourseRepos()
}

// ListMyCertificates returns all certificates for the current user
func ListMyCertificates(c echo.Context) error {
	initCertificateRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	certificates, err := certificateRepo.GetByUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch certificates"})
	}

	if certificates == nil {
		certificates = []*domain.Certificate{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"certificates": certificates,
		"count":        len(certificates),
	})
}

// GetCertificate returns a specific certificate
func GetCertificate(c echo.Context) error {
	initCertificateRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	id := c.Param("id")
	cert, err := certificateRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch certificate"})
	}
	if cert == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Certificate not found"})
	}

	// Verify ownership
	if cert.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	return c.JSON(http.StatusOK, cert)
}

// DownloadCertificatePDF generates and returns a PDF certificate
func DownloadCertificatePDF(c echo.Context) error {
	initCertificateRepos()

	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	id := c.Param("id")
	cert, err := certificateRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch certificate"})
	}
	if cert == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Certificate not found"})
	}

	// Verify ownership
	if cert.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	// Generate HTML certificate (simple version)
	htmlContent, err := generateCertificateHTML(cert)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate certificate"})
	}

	// Return as HTML (can be printed to PDF by browser)
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"certificate-%s.html\"", cert.CertificateNumber))
	return c.HTML(http.StatusOK, htmlContent)
}

// GenerateCertificateForUser creates a certificate for completing a course
func GenerateCertificateForUser(userID, courseID string) (*domain.Certificate, error) {
	initCertificateRepos()

	// Check if certificate already exists
	exists, _ := certificateRepo.ExistsForCourse(userID, courseID)
	if exists {
		return certificateRepo.GetByCourseAndUser(userID, courseID)
	}

	// Create new certificate
	cert := &domain.Certificate{
		UserID:            userID,
		CourseID:          courseID,
		CertificateNumber: certificateRepo.GenerateCertificateNumber(),
	}

	if err := certificateRepo.Create(cert); err != nil {
		return nil, err
	}

	// Log activity
	if activityRepo != nil {
		course, _ := courseRepo.GetByID(courseID)
		if course != nil {
			activityRepo.LogCourseComplete(userID, courseID, course.Title)
		}
	}

	return cert, nil
}

// generateCertificateHTML creates HTML content for the certificate
func generateCertificateHTML(cert *domain.Certificate) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Sertifikat - {{.CourseName}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Georgia', serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 40px;
        }
        .certificate {
            background: white;
            width: 800px;
            padding: 60px;
            text-align: center;
            border: 3px solid #667eea;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }
        .header {
            margin-bottom: 40px;
        }
        .logo {
            font-size: 28px;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 10px;
        }
        .title {
            font-size: 42px;
            color: #333;
            margin-bottom: 20px;
            text-transform: uppercase;
            letter-spacing: 4px;
        }
        .subtitle {
            font-size: 18px;
            color: #666;
            margin-bottom: 30px;
        }
        .recipient {
            font-size: 32px;
            color: #333;
            margin: 30px 0;
            font-style: italic;
        }
        .course-name {
            font-size: 24px;
            color: #667eea;
            margin: 20px 0;
            font-weight: bold;
        }
        .details {
            margin-top: 40px;
            font-size: 14px;
            color: #888;
        }
        .cert-number {
            font-family: monospace;
            margin-top: 10px;
        }
        @media print {
            body { background: white; padding: 0; }
            .certificate { box-shadow: none; border: 2px solid #667eea; }
        }
    </style>
</head>
<body>
    <div class="certificate">
        <div class="header">
            <div class="logo">🎓 LearnHub</div>
        </div>
        <h1 class="title">Sertifikat</h1>
        <p class="subtitle">Diberikan kepada</p>
        <p class="recipient">{{.UserName}}</p>
        <p class="subtitle">Atas keberhasilan menyelesaikan kursus</p>
        <p class="course-name">{{.CourseName}}</p>
        <div class="details">
            <p>Diterbitkan pada: {{.IssuedAt.Format "02 January 2006"}}</p>
            <p class="cert-number">ID: {{.CertificateNumber}}</p>
        </div>
    </div>
</body>
</html>
`
	t, err := template.New("certificate").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, cert); err != nil {
		return "", err
	}

	return buf.String(), nil
}
