package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Parse(w http.ResponseWriter, data map[string]interface{}) {
	// Parse the HTML templates
	tmpl, err := template.ParseFiles(
		"web/pages/index.html",
		"web/templates/header.html",
		"web/templates/left-bar.html",
		"web/templates/posts.html",
		"web/templates/create-post.html",
		"web/templates/js.html",
		"web/templates/response.html",
	)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error parsing templates:", err)
		http.Error(w, "Internal Server Error (Error parsing templates)", http.StatusInternalServerError)
		return
	}

	// Execute the template with data, including user and posts
	err = tmpl.Execute(w, data)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error executing template:", err)
		// Only call http.Error if nothing has been written to the response yet
		if w.Header().Get("Content-Type") == "" {
			http.Error(w, "Internal Server Error (Error executing template)", http.StatusInternalServerError)
		}
		return
	}
}
