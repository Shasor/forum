package handlers

import (
	"encoding/json"
	"forum/internal/db"
	"net/http"
)

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !IsCookieValid(w, r) {
		msg := "Method not Allowed"
		if !IsCookieValid(w, r) {
			msg = "You are not connected"
		}
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	user := GetUserFromCookie(w, r)
	notifs, _ := db.FetchNotificationsByUserId(user.ID)

	dataNotifs := map[string]interface{}{
		"notifData": notifs,
	}

	jsonData, err := json.Marshal(dataNotifs)
	if err != nil {
		http.Error(w, "Erreur lors de la sérialisation JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
