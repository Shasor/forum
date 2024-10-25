package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"forum/internal/db"
	"net/http"
	"os"
	"strings"
	"time"
)

var signingKey = os.Getenv("signingKeyForum")

// SetSession creates a session and stores it in a signed cookie
func SetSession(w http.ResponseWriter, username string) {
	// Base64 encode the session data (in this case, just the username)
	sessionData := base64.URLEncoding.EncodeToString([]byte(username))

	// Sign the session data
	signedData := sign(sessionData)

	// Create the cookie with the signed session data
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    signedData,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true, // Cookie is not accessible via JavaScript for security
		MaxAge:   int(1 * time.Hour),
	}
	// fmt.Println(cookie)
	http.SetCookie(w, &cookie)
}

// ClearSession removes the cookie
func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

// sign signs the session data using HMAC-SHA256
func sign(data string) string {
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return data + "." + signature
}

func GetUserFromCookie(w http.ResponseWriter, r *http.Request) *db.User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	parts := strings.Split(cookie.Value, ".")
	dataPart := parts[0]
	signaturePart := parts[1]

	// Recreate the HMAC signature for the data part
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(dataPart))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	username, _ := base64.URLEncoding.DecodeString(dataPart)
	user, err := db.SelectUserByUsername(string(username))

	// Compare the expected signature with the actual signature in the cookie
	if !hmac.Equal([]byte(signaturePart), []byte(expectedSignature)) || len(parts) != 2 || err != nil {
		ClearSession(w)
		return nil
	}
	return &user
}

func IsCookieValid(w http.ResponseWriter, r *http.Request) bool {
	if user := GetUserFromCookie(w, r); user != nil {
		return true
	}
	return false
}
