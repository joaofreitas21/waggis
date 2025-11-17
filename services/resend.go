package services

import (
	"fmt"
	"os"
	"bytes"
	"net/http"
	"encoding/json"
	"time"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	APIKey string
	From string
	To string
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

// GetEmailConfig loads email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	
	return &EmailConfig{
		APIKey:   getEnvOrDefault("RESEND_API_KEY", ""),
		From:     getEnvOrDefault("RESEND_FROM_EMAIL", ""),
		To:       getEnvOrDefault("RESEND_TO_EMAIL", ""),
	}
}


// SendEmail sends an email using SMTP
func SendEmail(config *EmailConfig, subject, body, replyName, replyEmail string) error {
	payload := resendPayload{
		From:    config.From,
		To:      []string{config.To},
		Subject: subject,
		Text:    body,
	}
	if replyEmail != "" {
		payload.ReplyTo = fmt.Sprintf("%s <%s>", replyName, replyEmail)
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal resend payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("build resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send resend request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("resend responded with %s", res.Status)
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

	return SendEmail(config, subject, body, name, email)
}


func getEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}