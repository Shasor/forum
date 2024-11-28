package handlers

import (
	"fmt"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
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

	// Debug : Affiche toutes les données du formulaire
	//fmt.Println("Form data:", r.Form)

	// Récupère l'ID du post à supprimer
	postIDStr := r.FormValue("id-post-to-delete")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	//fmt.Printf("Deleting post ID: %d\n", postID)

	// Supprime le post en appelant la fonction appropriée
	if err := db.DeletePostByID(postID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete post: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirige après la suppression
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}
