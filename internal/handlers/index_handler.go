package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Here get data

	//
	tmpl, err := template.ParseFiles("web/pages/index.html", "web/templates/header.html", "web/templates/left-bar.html", "web/templates/posts.html", "web/templates/create-post.html")
	if err != nil {
		http.Error(w, "Internal Server Error (Error parsing templates)", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error (Error executing template)", http.StatusInternalServerError)
		return
	}
}
