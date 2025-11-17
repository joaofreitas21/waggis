package models

import (
	"errors"
	"net/mail"
	"strings"
)

// EmailRequest represents the incoming email form data
type EmailRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Validate performs validation on the email request
func (e *EmailRequest) Validate() error {
	// Trim whitespace
	e.Name = strings.TrimSpace(e.Name)
	e.Email = strings.TrimSpace(e.Email)
	e.Subject = strings.TrimSpace(e.Subject)
	e.Message = strings.TrimSpace(e.Message)

	// Required field checks
	if e.Name == "" {
		return errors.New("name is required")
	}
	if len(e.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}

	if e.Email == "" {
		return errors.New("email is required")
	}
	if len(e.Email) > 255 {
		return errors.New("email must be 255 characters or less")
	}

	// Validate email format
	if _, err := mail.ParseAddress(e.Email); err != nil {
		return errors.New("invalid email format")
	}

	if e.Subject == "" {
		return errors.New("subject is required")
	}
	if len(e.Subject) > 200 {
		return errors.New("subject must be 200 characters or less")
	}

	if e.Message == "" {
		return errors.New("message is required")
	}
	if len(e.Message) > 5000 {
		return errors.New("message must be 5000 characters or less")
	}

	return nil
}

func (e *EmailRequest) Sanitize() {
	// Remove null bytes and control characters
	e.Name = strings.ReplaceAll(e.Name, "\x00", "")
	e.Email = strings.ReplaceAll(e.Email, "\x00", "")
	e.Subject = strings.ReplaceAll(e.Subject, "\x00", "")
	e.Message = strings.ReplaceAll(e.Message, "\x00", "")
}