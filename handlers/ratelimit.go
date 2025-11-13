package handlers

import (
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	clients map[string]*clientLimiter
	mu      sync.RWMutex
	limit   rate.Limit
	burst   int
}


type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	globalLimiter *RateLimiter
	once          sync.Once
)

// GetRateLimiter returns the single rate limiter instance
func GetRateLimiter() *RateLimiter {
	once.Do(func() {
		// Get rate limit config from env
		limitPerWindow := getEnvInt("EMAIL_RATE_LIMIT", 3)
		windowMinutes := getEnvInt("EMAIL_RATE_WINDOW_MINUTES", 15)

		// Calculate rate: requests per second
		requestsPerSecond := float64(limitPerWindow) / float64(windowMinutes*60)
		burst := limitPerWindow 

		globalLimiter = &RateLimiter{
			clients: make(map[string]*clientLimiter),
			limit:   rate.Limit(requestsPerSecond),
			burst:   burst,
		}

		// Start cleanup goroutine to remove old clients
		go globalLimiter.cleanupClients()
	})
	return globalLimiter
}

// Allow checks if the request from the given IP is allowed
// Also returns the time until next token is available if rate limited
func (rl *RateLimiter) Allow(ip string) (allowed bool, retryAfter time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[ip]
	if !exists {
		// Create new limiter for this IP
		client = &clientLimiter{
			limiter:  rate.NewLimiter(rl.limit, rl.burst),
			lastSeen: time.Now(),
		}
		rl.clients[ip] = client
	}

	client.lastSeen = time.Now()

	// Check if request is allowed
	if client.limiter.Allow() {
		return true, 0
	}

	// Calculate retry after time (estimate based on rate limit)
	reservation := client.limiter.Reserve()
	if reservation.OK() {
		retryAfter = reservation.Delay()
		reservation.Cancel()
	} else {
		// If reservation failed, estimate based on rate
		retryAfter = time.Duration(float64(time.Second) / float64(rl.limit))
	}

	return false, retryAfter
}

// cleanupClients removes clients that haven't been seen in the last hour
func (rl *RateLimiter) cleanupClients() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.lastSeen) > 1*time.Hour {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}
