package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/joaofreitas21/waggis/models"
	"github.com/joaofreitas21/waggis/services"
	"github.com/joaofreitas21/waggis/views"
)


type EmailResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	RetryAfter int   `json:"retry_after,omitempty"` // seconds until next request allowed
}

// SendEmail handles POST /api/email requests
func SendEmail(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Get client IP for rate limiting
	clientIP := GetClientIP(r)

	// Check rate limit
	limiter := GetRateLimiter()
	allowed, retryAfter := limiter.Allow(clientIP)
	if !allowed {
		response := EmailResponse{
			Success:    false,
			Error:      "Rate limit exceeded. Please try again later.",
			RetryAfter: int(retryAfter.Seconds()),
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse request body
	var emailReq models.EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		response := EmailResponse{
			Success: false,
			Error:   "Invalid request format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Sanitize input
	emailReq.Sanitize()

	// Validate request
	if err := emailReq.Validate(); err != nil {
		response := EmailResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get Resend configuration
	config := services.GetEmailConfig()

	// Validate SMTP configuration
	if config.APIKey == "" || config.From == "" || config.To == "" {
		log.Println("Email service is not fully configured (missing API key/from/to)")
		response := EmailResponse{
			Success: false,
			Error:   "Email service is not configured. Please contact the administrator.",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Send email via Resend
	if err := services.SendContactEmail(
		config,
		emailReq.Name,
		emailReq.Email,
		emailReq.Subject,
		emailReq.Message,
	); err != nil {
		log.Printf("Error sending email through Resend: %v", err)
		response := EmailResponse{
			Success: false,
			Error:   "Failed to send email. Please try again later.",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := EmailResponse{
		Success: true,
		Message: "Thanks for reaching out, I will be in contact soon!",
		
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ServeEmailSuccessCard(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Email sent successfully!"
	}
	
	views.EmailSuccessCard(message).Render(r.Context(), w)
}

func ServeEmailErrorCard(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Failed to send email. Please try again."
	}
	
	retryAfter := 0
	if retryStr := r.URL.Query().Get("retry_after"); retryStr != "" {
		if retry, err := strconv.Atoi(retryStr); err == nil {
			retryAfter = retry
		}
	}
	
	views.EmailErrorCard(message, retryAfter).Render(r.Context(), w)
}