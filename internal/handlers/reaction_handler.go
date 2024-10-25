package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"log"
	"net/http"
)

func ReactToPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !IsCookieValid(w, r) {
		return
	}

	var reqBody ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	postID := reqBody.PostID
	reaction := reqBody.Reaction

	user := GetUserFromCookie(w, r)
	if user == nil { // temporary
		return
	}

	err = db.UpdatePostReaction(user.ID, postID, reaction)
	if err != nil {
		log.Println("Error saving reaction:", err)
		http.Error(w, "Error saving reaction", http.StatusInternalServerError)
		return
	}

	// Fetch updated counts
	likes, dislikes, err := db.GetPostReactions(postID)
	if err != nil {
		http.Error(w, "Error fetching updated reactions", http.StatusInternalServerError)
		log.Println("Error fetching updated reactions:", err)
		return
	}

	// Send updated counts back as JSON response
	response := map[string]int{
		"likes":    likes,
		"dislikes": dislikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
