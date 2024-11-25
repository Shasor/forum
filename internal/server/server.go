package server

import (
	"bufio"
	"fmt"
	"forum/internal/handlers"
	"net/http"
	"os"
	"strings"
	"time"
	"forum/internal/api"
)

func InitServer() {
	if err := loadEnv(".env"); err != nil {
		fmt.Printf("Error load .env: %v\n", err)
	}
	var port = ":8080"
	server := NewServer(port, 10*time.Second, 10*time.Second, 30*time.Second, 2*time.Second, 1<<20)
	http.HandleFunc("/", recoverMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.ErrorsHandler(w, r, http.StatusNotFound, "The page you're looking for doesn't exist")
		} else {
			handlers.IndexHandler(w, r)
		}
	}))
	http.HandleFunc("/signup", recoverMiddleware(handlers.SignupHandler))
	http.HandleFunc("/login", recoverMiddleware(handlers.LoginHandler))
	http.HandleFunc("/create-post", recoverMiddleware(handlers.CreatePostHandler))
	http.HandleFunc("/react", recoverMiddleware(handlers.ReactToPost))
	http.HandleFunc("/follow", recoverMiddleware(handlers.FollowHandler))
	http.HandleFunc("/logout", recoverMiddleware(handlers.LogoutHandler))
	http.HandleFunc("/delete", recoverMiddleware(handlers.DeleteHandler))
	http.HandleFunc("/edit", recoverMiddleware(handlers.EditProfileHandler))
	http.HandleFunc("/users/", recoverMiddleware(api.GetUserInfo))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", staticMiddleware(http.StripPrefix("/static/", fs)))
	fmt.Printf("Starting server on http://localhost%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func NewServer(port string, readTimeout, writeTimeout, idleTimeout, readHeaderTimeout time.Duration, maxHeaderBytes int) *http.Server {
	return &http.Server{
		Addr:              port,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	return scanner.Err()
}
