package handlers

import (
	"fmt"
	"net/http"
	"forum/internal/db"
)

func NotificationHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet || !IsCookieValid(w, r) {
		msg := "Method not Allowed"
		if !IsCookieValid(w, r) {
			msg = "You are not connected"
		}
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	

	// Simuler une notification

	user := GetUserFromCookie(w,r)

	fmt.Println(user.ID)

	notifs, err := db.FetchNotificationsByUserId(user.ID)
	if err != nil{
		fmt.Println("Error fetching Notifs : ", )
	}

	//fmt.Println("Notifications : ", notifs)
	fmt.Println("Notification envoyée", user.Username )

	//notification := "Nouvelle notification disponible !"

	// Envoyer la notification au client
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, notifs)))

}