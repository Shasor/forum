package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"net/http"
	"strconv"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Resp.Msg = append(Resp.Msg, "Method not Allowed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ErrorsHandler(w, r, http.StatusBadRequest, "JSON decoding error")
		return
	}

	userId, _ := strconv.Atoi(data["userId"]) // to do: error handle
	user, _ := db.SelectUserById(userId)

	userActivities := db.GetUserActivitiesByID(user.ID)

	dataUser := map[string]interface{}{
		"userData":     user,
		"userActivity": userActivities,
	}

	jsonData, err := json.Marshal(dataUser)
	if err != nil {
		http.Error(w, "Erreur lors de la sérialisation JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
