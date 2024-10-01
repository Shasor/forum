package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

var sessionStore = make(map[string]string)

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SetSession(w http.ResponseWriter, r *http.Request, username string) {
	sessionID := generateSessionID()
	sessionStore[sessionID] = username

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		delete(sessionStore, cookie.Value)
		cookie.Expires = time.Now().Add(-1 * time.Hour)
		http.SetCookie(w, cookie)
	}
}

func IsLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	_, exists := sessionStore[cookie.Value]
	return exists
}

func GetSessionUsername(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	username, exists := sessionStore[cookie.Value]
	return username, exists
}
