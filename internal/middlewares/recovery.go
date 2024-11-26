package middlewares

import (
	"fmt"
	"forum/internal/handlers"
	"log"
	"net/http"
)

func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				handlers.ErrorsHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
}
