package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func ErrorsHandler(w http.ResponseWriter, r *http.Request, status int, msg string) {
	// Réinitialiser le ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.(http.Flusher).Flush()

	// Définir le statut après la réinitialisation
	w.WriteHeader(status)

	data := ErrorPageData{
		StatusCode:   status,
		StatusText:   http.StatusText(status),
		ErrorMessage: msg,
		ErrorDetails: "",
	}

	tmpl, err := template.ParseFiles("web/pages/errors.html")
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
