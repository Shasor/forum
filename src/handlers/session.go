package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"login/src/models"
	"net/http"
	"strings"
	"time"
)

// Signing key (should be kept out of the code)
var signingKey = []byte("a very strong signing key")

// SetSession creates a session and stores it in a signed cookie
func SetSession(w http.ResponseWriter, r *http.Request, username string) {
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
	http.SetCookie(w, &cookie)
}

// ClearSession removes the cookie
func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

// IsCookieExist checks whether the session cookie exists
func IsCookieExist(r *http.Request) bool {
	_, err := r.Cookie("session_token")
	return err == nil
}

// GetSessionUsername returns the username from the session cookie
func GetSessionUsername(db *sql.DB, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	// verify the session cookie and get the session data (username)
	sessionData, err := verifyHash(cookie.Value)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded session data
	usernameBytes, err := base64.URLEncoding.DecodeString(sessionData)
	if err != nil {
		return "", err
	}
	if !models.UserExist(db, string(usernameBytes)) {
		return "", ErrInvalidCookie
	}

	return string(usernameBytes), nil
}

// sign signs the session data using HMAC-SHA256
func sign(data string) string {
	h := hmac.New(sha256.New, signingKey)
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return data + "." + signature
}

// verifyHash checks the session data signature, then returns the username (base64) in the cookie, returns [ErrInvalidCookie] if the hash doesn't match
func verifyHash(signedData string) (string, error) {
	parts := strings.Split(signedData, ".")
	if len(parts) != 2 {
		return "", ErrInvalidCookie
	}

	dataPart := parts[0]
	signaturePart := parts[1]

	// Recreate the HMAC signature for the data part
	h := hmac.New(sha256.New, signingKey)
	h.Write([]byte(dataPart))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Compare the expected signature with the actual signature in the cookie
	if !hmac.Equal([]byte(signaturePart), []byte(expectedSignature)) {
		return "", ErrInvalidCookie
	}
	return dataPart, nil
}
