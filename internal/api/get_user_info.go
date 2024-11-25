package api

import (
	"forum/internal/db"
	"strconv"
	"net/http"
	"encoding/json"
	"fmt"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'utilisateur depuis l'URL
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	user, _ := db.SelectUserById(id)
	userPosts := db.GetPostFromUserById(id)

	dataUser := map[string]interface{}{
		"userData": user,
		"userPosts": userPosts,
	}

	fmt.Println(len(userPosts))

	// Retourner les informations de l'utilisateur en format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataUser)
}