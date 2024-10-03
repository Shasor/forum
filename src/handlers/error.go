package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type ErrorPageData struct {
	Code    int
	Message string
}

func RenderErrorPage(w http.ResponseWriter, code int, message string) {
	tmpl := template.Must(template.ParseFiles("./web/template/error.html"))

	data := ErrorPageData{
		Code:    code,
		Message: message,
	}

	log.Printf("Erreur %d: %s", code, message)

	w.WriteHeader(code)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}

}
