package middlewares

import (
	"forum/internal/handlers"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count      int
	lastAccess time.Time
	blocked    bool
	blockUntil time.Time
}

var (
	clients     = make(map[string]*client)
	mu          sync.Mutex
	maxRequests = 5
	interval    = time.Second
	banDuration = 10 * time.Second
)

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()
		mu.Lock()
		defer mu.Unlock()
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{count: 0, lastAccess: now}
		}
		c := clients[ip]
		if c.blocked {
			if now.Before(c.blockUntil) {
				handlers.ErrorsHandler(w, r, http.StatusTooManyRequests, "")
				return
			}
			c.blocked = false
		}
		if now.Sub(c.lastAccess) >= interval {
			c.count = 0
			c.lastAccess = now
		}
		if c.count >= maxRequests {
			c.blocked = true
			c.blockUntil = now.Add(banDuration)
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}
		c.count++
		next.ServeHTTP(w, r)
	}
}
