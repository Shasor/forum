package handlers

import (
	"html/template"
	"net/http"
)

func ErrorsHandler(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)

	data := ErrorPageData{
		StatusCode:   status,
		StatusText:   http.StatusText(status),
		ErrorMessage: msg,
		ErrorDetails: "", // Vous pouvez ajouter des détails spécifiques ici si nécessaire
	}

	tmpl := template.Must(template.ParseFiles("web/pages/errors.html"))
	tmpl.Execute(w, data)
}
