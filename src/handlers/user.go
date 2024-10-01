package handlers

import (
	"database/sql"
	"log"
	"login/src/models" // Remplace par le chemin correct vers tes modèles
	"net/http"
)

func DeleteAccountHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode
	if r.Method != "POST" {
		log.Println("Mauvaise méthode :", r.Method)
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username, ok := GetSessionUsername(r)
	if !ok {
		log.Println("Aucun utilisateur connecté.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Log du nom d'utilisateur
	log.Println("Suppression de l'utilisateur :", username)

	// Suppression de l'utilisateur de la base de données
	err := models.DeleteUserByUsername(db, username)
	if err != nil {
		log.Println("Erreur lors de la suppression du compte :", err)
		http.Error(w, "Impossible de supprimer le compte", http.StatusInternalServerError)
		return
	}

	// Log de succès
	log.Println("Compte supprimé avec succès :", username)

	// Supprimer la session de l'utilisateur
	ClearSession(w, r)

	// Rediriger vers la page de confirmation
	http.Redirect(w, r, "/account-deleted", http.StatusSeeOther)
}
