package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"forum/internal/db"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

var signingKey = os.Getenv("signingKeyForum")

// SetSession creates a session and stores it in a signed cookie
func SetSession(w http.ResponseWriter, username string) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}
	usr, _ := db.SelectUserByUsername(username)

	if connected, _ := db.IsUserConnected(usr.ID); connected {
		Resp = Response{}
		Resp.Msg = append(Resp.Msg, "You're already login in an other browser")
		return
	} else {
		db.AddConnectedUser(usr.ID, sessionID.String())
	}

	// Create the cookie with the signed session data
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionID.String(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true, // Cookie is not accessible via JavaScript for security
		MaxAge:   int(1 * time.Hour),
	}

	http.SetCookie(w, &cookie)
}

// ClearSession removes the cookie
func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	usr, _ := db.SelectUserByUsername(cookie.Name)
	db.DeleteConnectedUser(usr.ID)
}

// sign signs the session data using HMAC-SHA256
func sign(data string) string {
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return data + "." + signature
}

// GetUserFromCookie retrieves a user object based on the session token in the cookie.
func GetUserFromCookie(w http.ResponseWriter, r *http.Request) *db.User {
	// Retrieve the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Cookie not found or invalid
		return nil
	}

	sessionUUID := cookie.Value
	DB := db.GetDB()
	defer DB.Close()

	// Query to find the user ID from the sessions table using the UUID
	var userID int
	err = DB.QueryRow(`SELECT connected_user FROM sessions WHERE uuid = ?`, sessionUUID).Scan(&userID)
	if err != nil {
		// Session not found or other query error
		return nil
	}

	// Query to fetch user details from the users table
	var user db.User
	err = DB.QueryRow(
		`SELECT id, role, username, email, picture, password, created_at 
		 FROM users 
		 WHERE id = ?`, userID,
	).Scan(
		&user.ID, &user.Role, &user.Username, &user.Email, &user.Picture, &user.Password,
	)
	if err != nil {
		// User not found or query error
		return nil
	}

	// Return the populated user object
	return &user
}

func IsCookieValid(w http.ResponseWriter, r *http.Request) bool {
	if user := GetUserFromCookie(w, r); user != nil {
		return true
	}
	return false
}
