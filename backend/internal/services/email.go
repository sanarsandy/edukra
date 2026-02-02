package services

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"gopkg.in/gomail.v2"
)

type EmailService struct{}

var Email = &EmailService{}

// getSetting retrieves a setting value from the database
func getSetting(key string) string {
	var value string
	err := db.DB.QueryRow("SELECT value FROM settings WHERE key = $1", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

// SendEmail sends an email using the configured SMTP settings
func (s *EmailService) SendEmail(to, subject, body string) error {
	host := getSetting("smtp_host")
	portStr := getSetting("smtp_port")
	username := getSetting("smtp_username")
	password := getSetting("smtp_password")
	fromEmail := getSetting("smtp_from_email")
	fromName := getSetting("smtp_from_name")

	if host == "" || portStr == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP settings are not configured")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	if fromName != "" {
		m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	} else {
		m.SetHeader("From", fromEmail)
	}
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Allow skip verify for self-signed or common dev setups

	return d.DialAndSend(m)
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to LearnHub!"
	body := fmt.Sprintf(`
		<h2>Welcome, %s!</h2>
		<p>Thank you for joining LearnHub. We are excited to have you on board.</p>
		<p>Start exploring our courses and happy learning!</p>
		<br>
		<p>Best regards,</p>
		<p>The LearnHub Team</p>
	`, name)
	return s.SendEmail(to, subject, body)
}

// SendResetPasswordEmail sends a password reset link
func (s *EmailService) SendResetPasswordEmail(to, token string) error {
	frontendURL := getSetting("frontend_url")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default
	}
	
	// Remove trailing slash if present
	if len(frontendURL) > 0 && frontendURL[len(frontendURL)-1] == '/' {
		frontendURL = frontendURL[:len(frontendURL)-1]
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
		<h2>Reset Password Request</h2>
		<p>You requested a password reset. Click the link below to reset your password:</p>
		<p><a href="%s" style="background-color: #4F46E5; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a></p>
		<p>Link is valid for 15 minutes.</p>
		<p>If you didn't request this, please ignore this email.</p>
		<br>
		<p>Best regards,</p>
		<p>The LearnHub Team</p>
	`, resetLink)
	return s.SendEmail(to, subject, body)
}

// SendTestEmail sends a test email to verify SMTP configuration
func (s *EmailService) SendTestEmail(to string) error {
	subject := "Test Email from LearnHub"
	body := `
		<h2>SMTP Configuration Verified!</h2>
		<p>This is a test email from your LearnHub instance.</p>
		<p>If you see this, your SMTP settings are working correctly.</p>
	`
	return s.SendEmail(to, subject, body)
}
