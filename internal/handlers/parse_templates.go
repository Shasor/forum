package handlers

import (
	"html/template"
	"net/http"
)

func Parse(w http.ResponseWriter, data map[string]interface{}) {
	tmpl, err := template.ParseFiles("web/pages/index.html", "web/templates/header.html", "web/templates/left-bar.html", "web/templates/posts.html", "web/templates/create-post.html", "web/templates/js.html", "web/templates/response.html")
	if err != nil {
		http.Error(w, "Internal Server Error (Error parsing templates)", http.StatusInternalServerError)
		return
	}
	// fmt.Println(data)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error (Error executing template)", http.StatusInternalServerError)
		return
	}
}
