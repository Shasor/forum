package server

import (
	"fmt"
	"forum/internal/handlers"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func staticMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Liste des extensions de fichiers autorisées
		allowedExt := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".css": true, ".js": true,
		}
		// Vérifiez l'extension du fichier demandé
		ext := strings.ToLower(filepath.Ext(r.URL.Path))
		if !allowedExt[ext] {
			handlers.ErrorsHandler(w, r, http.StatusForbidden, "access prohibited")
			return
		}
		// Si l'extension est autorisée, servez le fichier
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				handlers.ErrorsHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
}
