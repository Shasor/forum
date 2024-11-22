package router

import (
	"database/sql"
	"forum/internal/handlers"
	"forum/internal/security"
	"net/http"
)

// SetupRoutes configures the application routes
func SetupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize rate limiter
	rateLimiter := security.NewRateLimiterMap()

	// Public routes
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.Handle("/signup", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.SignupHandler(w, r)
	})))

	mux.Handle("/login", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r)
	})))

	// Post-related routes
	mux.Handle("/create-post", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePostHandler(w, r)
	})))

	mux.Handle("/react", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ReactToPost(w, r)
	})))

	// User-related routes
	mux.Handle("/follow", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.FollowHandler(w, r)
	})))

	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// Account management routes
	mux.Handle("/delete", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteHandler(w, r)
	})))

	mux.Handle("/edit", rateLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.EditProfileHandler(w, r)
	})))

	// Static files (CSS, images, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
