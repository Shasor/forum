package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"net/http"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CategorieID int `json:"categorieId"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	user := GetUserFromCookie(w, r)

	// Vérifie si l'utilisateur suit déjà la catégorie
	alreadyFollowing, err := db.AlreadyFollowingCategory(req.CategorieID, user.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if alreadyFollowing {
		// Retirer le suivi
		err = db.StopFollowingCategory(req.CategorieID, user.ID)
	} else {
		// Commencer à suivre la catégorie
		err = db.StartFollowingCategory(req.CategorieID, user.ID)
	}

	if err != nil {
		http.Error(w, "Error updating follow status", http.StatusInternalServerError)
		return
	}
	// Répond avec le nouvel état de suivi
	json.NewEncoder(w).Encode(map[string]bool{
		"isFollowing": !alreadyFollowing,
	})
}
