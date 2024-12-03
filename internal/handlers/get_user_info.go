package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"net/http"
	"strconv"
	"strings"
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

	purgedOtherID := strings.ReplaceAll(data["otherID"], " ", "")
	otherID, err := strconv.Atoi(purgedOtherID) // to do: error handle
	if err != nil {
		http.Error(w, "error while processing data ", http.StatusInternalServerError)
		return
	}
	other, _ := db.SelectUserByID(otherID)
	otherActivities := db.GetUserActivitiesByID(other.ID)

	dataUser := map[string]interface{}{
		"otherData":     other,
		"otherActivity": otherActivities,
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
