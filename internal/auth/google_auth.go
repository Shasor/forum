package auth

import (
	"encoding/json"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/utils"
	"net/http"
	"net/url"
	"os"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomState()

	params := url.Values{}
	params.Add("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	params.Add("redirect_uri", "https://localhost:8080/auth/google/callback")
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email")
	params.Add("state", state)
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: true})
	http.Redirect(w, r, GoogleAuthURL+"?"+params.Encode(), http.StatusFound)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	oauthCookie, err := r.Cookie("oauth_state")
	if err != nil || state != oauthCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := GoogleExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	userInfo, err := GoogleGetUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	user, err := GoogleCreateOrUpdateUser(userInfo)
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

func GoogleExchangeCodeForToken(code string) (string, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	values.Add("redirect_uri", "https://localhost:8080/auth/google/callback")
	values.Add("grant_type", "authorization_code")

	resp, err := http.PostForm(GoogleTokenURL, values)
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

func GoogleGetUserInfo(token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", GoogleUserInfoURL, nil)
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

func GoogleCreateOrUpdateUser(userInfo map[string]interface{}) (*db.User, error) {
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