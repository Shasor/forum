package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"login/src/database"
	"login/src/handlers"
)

func main() {
	// Initialiser la base de données
	db, err := database.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Créer la table des utilisateurs
	database.InitAllDB(db)

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
		handlers.DeleteAccountHandler(db, w, r)
	})

	// http.HandleFunc("/account-deleted", func(w http.ResponseWriter, r *http.Request) {
	// 	if !handlers.IsCookieExist(r) || r.Method != http.MethodPost {
	// 		handlers.LoginHandler(db, w, r)
	// 	}
	// })

	// New route to handle post creation
	http.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlers.CreatePostHandler(db, w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Servir les fichiers statiques (CSS)
	http.Handle("/web/static/css/", http.StripPrefix("/web/static/css/", http.FileServer(http.Dir("./web/static/css/"))))
	http.Handle("/web/img/", http.StripPrefix("/web/img/", http.FileServer(http.Dir("./web/img/"))))

	server := &http.Server{
		Addr:              ":8080",
		MaxHeaderBytes:    1 << 20,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	fmt.Println("Serveur démarré : http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
