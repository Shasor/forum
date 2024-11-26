package middlewares

import (
	"forum/internal/handlers"
	"net/http"
)

var Paths []string

// NotFoundMiddleware checks if the requested path matches any registered routes and returns 404 if not.
func NotFoundMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path matches any registered routes
		matched := false
		for _, path := range Paths { // Assuming `routes` is accessible here or passed as an argument
			if r.URL.Path == path {
				matched = true
				break
			}
		}

		// If no match is found, respond with 404 Not Found
		if !matched {
			handlers.ErrorsHandler(w, r, http.StatusNotFound, "The page you're looking for doesn't exist")
			return
		}

		// Call the next handler in case of a match (this should not happen here due to prior registration)
		next.ServeHTTP(w, r)
	}
}
