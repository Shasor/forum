package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// RenderTemplate renders a template with the given data
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("web", "template", tmpl) // Ensure the path is correct

	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not load template: %v", err), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		return
	}
}
