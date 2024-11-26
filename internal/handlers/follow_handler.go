package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"net/http"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Resp.Msg = append(Resp.Msg, "Method not Allowed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var req struct {
		CategorieID int `json:"categorieId"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	user := GetUserFromCookie(w, r)

	// Vérifie si l'utilisateur suit déjà la catégorie
	alreadyFollowing, err := db.AlreadyFollowingCategory(req.CategorieID, user.ID)
	if err != nil {
		panic(err)
	}

	if alreadyFollowing {
		// Retirer le suivi
		err = db.StopFollowingCategory(req.CategorieID, user.ID)
	} else {
		// Commencer à suivre la catégorie
		err = db.StartFollowingCategory(req.CategorieID, user.ID)
	}

	if err != nil {
		panic(err)
	}
	// Répond avec le nouvel état de suivi
	json.NewEncoder(w).Encode(map[string]bool{
		"isFollowing": !alreadyFollowing,
	})
}
