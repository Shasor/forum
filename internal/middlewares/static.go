package middlewares

import (
	"forum/internal/handlers"
	"net/http"
	"path/filepath"
	"strings"
)

func StaticMiddleware(next http.Handler) http.Handler {
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
