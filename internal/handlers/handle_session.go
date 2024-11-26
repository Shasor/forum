package handlers

import (
	"forum/internal/db"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func SetSession(w http.ResponseWriter, username string) {
	sessionUUID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	// Fetch the user by username
	user, err := db.SelectUserByUsername(username)
	if err != nil {
		panic(err)
	}

	if connected, _ := db.IsUserConnected(user.ID); connected {
		Resp.Msg = append(Resp.Msg, "You were already connected elsewhere")
		id, _ := db.GetUUIDByUserID(user.ID)
		db.DeleteConnectedUser(id)
		SetSession(w, username)
		return
	}

	err = db.AddConnectedUser(user.ID, sessionUUID.String())
	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionUUID.String(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		panic(err)
	}

	err = db.DeleteConnectedUser(cookie.Value)
	if err != nil {
		panic(err)
	}

	cookie = &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func GetUserFromCookie(w http.ResponseWriter, r *http.Request) *db.User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}

	userID, err := db.GetUserIDBySessionUUID(cookie.Value)
	if err != nil {
		ClearSession(w, r)
		return nil
	}

	user, err := db.SelectUserByID(userID)
	if err != nil {
		ClearSession(w, r)
		return nil
	}

	return &user
}

func IsCookieValid(w http.ResponseWriter, r *http.Request) bool {
	return GetUserFromCookie(w, r) != nil
}
