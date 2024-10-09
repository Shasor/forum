package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"login/src/models" // Remplace par le chemin correct vers tes modèles
	"net/http"
)

func DeleteAccountHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if !IsCookieExist(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Vérification de la méthode
	if r.Method != "POST" {
		log.Println("Mauvaise méthode :", r.Method)
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username, err := GetSessionUsername(r)
	if err != nil {
		log.Println("Aucun utilisateur connecté.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Log du nom d'utilisateur
	log.Println("Suppression de l'utilisateur :", username)

	// Suppression de l'utilisateur de la base de données
	err = models.DeleteUserByUsername(db, username)
	if err != nil {
		log.Println("Erreur lors de la suppression du compte :", err)
		http.Error(w, "Impossible de supprimer le compte", http.StatusInternalServerError)
		return
	}

	// Log de succès
	log.Println("Compte supprimé avec succès :", username)

	// Supprimer la session de l'utilisateur
	ClearSession(w, r)

	// Afficher la page de confirmation
	tmpl := template.Must(template.ParseFiles("./web/template/account-deleted.html"))
	tmpl.Execute(w, nil)
}

// UserProfilePage displays the profile of a specific user by their username
func UserProfilePage(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if !IsCookieExist(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extract the username from the URL
	path := r.URL.Path
	username := path[len("/profile/"):] // Get everything after "/profile/"

	// Fetch the user details from the database
	user, err := models.SelectUser(db, username)
	if err != nil {
		if err.Error() == "utilisateur non trouvé" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching user profile:", err)
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	// Render the profile page
	tmpl := template.Must(template.ParseFiles("./web/template/profile.html"))
	tmpl.Execute(w, map[string]interface{}{
		"User": user,
	})
}
