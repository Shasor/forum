package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
)

type ErrorPageData struct {
	Code    int
	Message string
}

var ErrInvalidCookie = errors.New("invalid cookie")

func ErrorsExit(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
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
