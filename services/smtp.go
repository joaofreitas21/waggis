package services

import (
	"fmt"
	"os"

	"gopkg.in/mail.v2"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
}

// GetEmailConfig loads email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	port := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	return &EmailConfig{
		Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		Port:     port,
		Username: getEnvOrDefault("SMTP_USER", ""),
		Password: getEnvOrDefault("SMTP_PASSWORD", ""),
		From:     getEnvOrDefault("SMTP_FROM_EMAIL", ""),
		To:       getEnvOrDefault("SMTP_TO_EMAIL", ""),
	}
}

// SendEmail sends an email using SMTP
func SendEmail(config *EmailConfig, subject, body, fromName, fromEmail string) error {
	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, fromEmail))
	m.SetHeader("To", config.To)
	m.SetHeader("Subject", subject)

	// Set body
	m.SetBody("text/plain", body)

	// Create dialer
	d := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendContactEmail is a convenience function to send a contact form email
func SendContactEmail(config *EmailConfig, name, email, subject, message string) error {
	// Format email body
	body := fmt.Sprintf(`New contact form submission via waggis website:

NAME: %s
EMAIL: %s
SUBJECT: %s

MESSAGE:
%s
`, name, email, subject, message)

	return SendEmail(config, subject, body, name, config.From)
}


func getEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}