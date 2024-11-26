package middlewares

import (
	"forum/internal/handlers"
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter stores rate-limiting information for clients
type RateLimiter struct {
	limit      int           // Max tokens available (bucket size)
	refillRate time.Duration // How often tokens are refilled
	tokens     int           // Current available tokens
	lastRefill time.Time     // Last time tokens were refilled
	mu         sync.Mutex    // To synchronize access to tokens
}

// NewRateLimiter creates a new RateLimiter
func newRateLimiter(limit int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:      limit,
		refillRate: refillRate,
		tokens:     limit,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed or rate-limited
func (rl *RateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate time since the last refill
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed time
	refillTokens := int(elapsed / rl.refillRate)
	if refillTokens > 0 {
		rl.tokens = min(rl.limit, rl.tokens+refillTokens) // Refill the tokens
		rl.lastRefill = now                               // Reset the refill time
	}

	// If tokens are available, consume one and allow the request
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	// No tokens available, reject the request
	return false
}

// Helper function to find the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimiterMap keeps track of rate limiters for each client (IP address)
type RateLimiterMap struct {
	clients map[string]*RateLimiter
	mu      sync.Mutex
}

// NewRateLimiterMap initializes a new map for tracking rate limiters
func NewRateLimiterMap() *RateLimiterMap {
	return &RateLimiterMap{
		clients: make(map[string]*RateLimiter),
	}
}

// GetRateLimiter retrieves or creates a rate limiter for the given IP address
func (rlm *RateLimiterMap) getRateLimiter(ip string) *RateLimiter {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	// Check if there's already a rate limiter for this IP
	limiter, exists := rlm.clients[ip]
	if !exists {
		// Create a new rate limiter (e.g., 5 requests per second)
		limiter = newRateLimiter(5, time.Second)
		rlm.clients[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware applies rate limiting based on the client's IP address
func (rlm *RateLimiterMap) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the client's IP address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			panic(err)
		}

		// Get or create a rate limiter for this IP address
		limiter := rlm.getRateLimiter(ip)

		// Check if the request is allowed
		if !limiter.allow() {
			handlers.ErrorsHandler(w, r, http.StatusTooManyRequests, "Too Many Requests")
			return
		}

		// If allowed, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
