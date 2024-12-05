package handlers

import (
	"forum/internal/db"
	"net/http"
	"encoding/json"
	"strconv"
)

func NotificationClearHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost || !IsCookieValid(w, r) {
		Resp = Response{Msg: []string{
			map[bool]string{
				true:  "Method not Allowed",
				false: "You are not connected",
			}[r.Method != http.MethodPost],
		}}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var req struct {
		UserID string `json:"userID"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//fmt.Println(req.UserID)

	userId, err := strconv.Atoi(req.UserID[1:])
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	err = db.MarkAllNotificationsAsRead(userId)
	if err != nil{
		http.Error(w, "Error updating notifications", http.StatusBadRequest)
		return 
	}

	//fmt.Printf("Utilisateur ID reçu : %d\n", userId)
	response := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Notifications cleared successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}