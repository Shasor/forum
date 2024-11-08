package handlers

import (
	"html/template"
	"net/http"
)

func ErrorsHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	data := ErrorPageData{
		StatusCode:   status,
		StatusText:   http.StatusText(status),
		ErrorMessage: "An unexpected error has occurred",
		ErrorDetails: "", // Vous pouvez ajouter des détails spécifiques ici si nécessaire
	}

	tmpl := template.Must(template.ParseFiles("web/pages/errors.html"))
	tmpl.Execute(w, data)
}
