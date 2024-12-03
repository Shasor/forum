package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/utils"
	"net/http"
	"net/url"
	"os"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()

	params := url.Values{}
	params.Add("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	params.Add("redirect_uri", "https://localhost:8080/auth/google/callback")
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email")
	params.Add("state", state)
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, googleAuthURL+"?"+params.Encode(), http.StatusFound)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	oauthCookie, err := r.Cookie("oauth_state")
	if err != nil || state != oauthCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	userInfo, err := getUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	user, err := createOrUpdateUser(userInfo)
	if err != nil {
		http.Error(w, "Failed to create or update user", http.StatusInternalServerError)
		return
	}

	// Créer une session pour l'utilisateur
	handlers.ClearSession(w, r, "oauth_state")
	handlers.SetSession(w, user.Username)
	// ...

	http.Redirect(w, r, "/", http.StatusFound)
}

func exchangeCodeForToken(code string) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	values.Add("redirect_uri", "https://localhost:8080/auth/google/callback")
	values.Add("grant_type", "authorization_code")

	resp, err := http.PostForm(googleTokenURL, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func getUserInfo(token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", googleUserInfoURL, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func createOrUpdateUser(userInfo map[string]interface{}) (*db.User, error) {
	// Implémentez la logique pour créer ou mettre à jour l'utilisateur dans votre base de données
	// Utilisez les informations de userInfo pour remplir les champs de l'utilisateur
	// Retournez l'utilisateur créé ou mis à jour
	fmt.Println(userInfo)
	fmt.Println()

	var user *db.User
	if !db.UserExistByEmail(userInfo["email"].(string)) {
		picture, err := utils.GetFileFromURL(userInfo["picture"].(string))
		if err != nil {
			panic(err)
		}
		user, err = db.CreateUser("user", userInfo["given_name"].(string), userInfo["email"].(string), picture, "")
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		user, err = db.SelectUserByEmail(userInfo["email"].(string))
		if err != nil {
			panic(err)
		}
	}
	return user, nil
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
