package server

import (
	"fmt"
	"forum/internal/handlers"
	"net/http"
	"time"
)

func InitServer() {
	var port = ":8080"
	server := NewServer(port, 10*time.Second, 10*time.Second, 30*time.Second, 2*time.Second, 1<<20)

	http.HandleFunc("/", handlers.IndexHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
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
