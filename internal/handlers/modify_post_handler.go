package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func ModifyPostHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost || !IsCookieValid(w, r) {
		msg := "Method not Allowed"
		if !IsCookieValid(w, r) {
			msg = "You are not connected"
		}
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	
	// Analyse les données du formulaire
	if err := r.ParseForm(); err != nil {
		fmt.Println("ParseForm error:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limite de 10MB
		fmt.Println("ParseMultipartForm error:", err)
		http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
		return
	}
	
	postIDStr := r.MultipartForm.Value["post-to-modify"][0]
	contentPost := r.MultipartForm.Value["content_post"][0]
	
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}


	newContent := contentPost + " (Modified)"

	if err := db.ModifyContentPostByID(postID, newContent); err != nil{
		http.Error(w, fmt.Sprintf("Failed to update content: %v", err), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w,r, "/", http.StatusSeeOther)
}