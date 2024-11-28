package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func Parse(w http.ResponseWriter, data map[string]interface{}) {
	funcMap := template.FuncMap{
		"dict": dict,
	}

	tmpl := template.New("").Funcs(funcMap)

	tmpl, err := tmpl.ParseFiles(
		"web/pages/index.html",
		"web/templates/header.html",
		"web/templates/left-bar.html",
		"web/templates/posts.html",
		"web/templates/create-post.html",
		"web/templates/js.html",
		"web/templates/response.html",
		"web/templates/post.html",
		"web/templates/comment.html",
	)
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		panic(err)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		if w.Header().Get("Content-Type") == "" {
			panic(err)
		}
		return
	}
}
