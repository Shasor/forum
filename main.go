package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"login/src/database"
	"login/src/handlers"
)

func main() {
	// Initialiser la base de données
	db, err := database.InitDB("db/users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Créer la table des utilisateurs
	database.CreateTable(db)

	// Configurer les routes
	http.HandleFunc("/", handlers.LoginPage)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(db, w, r)
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignupPage(db, w, r)
	})

	http.HandleFunc("/dashboard", handlers.DashboardPage)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/delete-account", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlers.DeleteAccountHandler(db, w, r)
		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/account-deleted", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./web/template/account-deleted.html"))
		tmpl.Execute(w, nil)
	})

	// Servir les fichiers statiques (CSS)
	http.Handle("/web/static/css/", http.StripPrefix("/web/static/css/", http.FileServer(http.Dir("./web/static/css/"))))

	server := &http.Server{
		Addr:              ":8080",
		MaxHeaderBytes:    1 << 20,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	fmt.Println("Serveur démarré sur le port 8080")
	log.Fatal(server.ListenAndServe())
}
