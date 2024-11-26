package server

import (
	"fmt"
	"forum/internal/middlewares"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Middleware defines a function that takes an http.HandlerFunc and returns an http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route represents a route in the server.
type Route struct {
	Path    string
	Handler http.HandlerFunc
}

// Server represents our HTTP server.
type Server struct {
	port              string
	routes            []Route
	middlewares       []Middleware
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	maxHeaderBytes    int
}

// NewServer creates a new instance of Server with specified configurations.
func NewServer(port string, readTimeout, writeTimeout, idleTimeout, readHeaderTimeout time.Duration, maxHeaderBytes int) *Server {
	return &Server{
		port:              port,
		routes:            []Route{},
		middlewares:       []Middleware{},
		readTimeout:       readTimeout,
		writeTimeout:      writeTimeout,
		idleTimeout:       idleTimeout,
		readHeaderTimeout: readHeaderTimeout,
		maxHeaderBytes:    maxHeaderBytes,
	}
}

// Use adds a middleware to the server.
func (s *Server) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

// Handle adds a route to the server.
func (s *Server) Handle(path string, handler http.HandlerFunc) {
	s.routes = append(s.routes, Route{Path: path, Handler: handler})
	middlewares.Paths = append(middlewares.Paths, path)
}

// Start launches the server on the specified port with the defined settings.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	for _, route := range s.routes {
		handler := route.Handler
		// Apply all middlewares to the handler
		for _, mw := range s.middlewares {
			handler = mw(handler)
		}
		// Register the final handler with all middlewares applied
		mux.HandleFunc(route.Path, handler)
	}

	server := &http.Server{
		Addr:              s.port,
		ReadTimeout:       s.readTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		MaxHeaderBytes:    s.maxHeaderBytes,
		Handler:           mux,
		TLSConfig:         loadTLSConfig(),
	}
	mux.Handle("/static/", middlewares.StaticMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	fmt.Printf("Starting server on https://localhost%s\n", s.port)
	return server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
}

func checkAndCreateCert(certFile, keyFile string) error {
	// Check if the certificate file exists
	_, err := os.Stat(certFile)
	if os.IsNotExist(err) {
		// If the certificate does not exist, create it using the shell script
		fmt.Println("Certificate not found, generating new one...")
		return generateCertFromScript()
	}

	// Check if the key file exists
	_, err = os.Stat(keyFile)
	if os.IsNotExist(err) {
		// If the key file does not exist, create it using the shell script
		fmt.Println("Key file not found, generating new one...")
		return generateCertFromScript()
	}

	// If both the certificate and key exist, return nil (no error)
	return nil
}

func generateCertFromScript() error {
	// Execute the generate_cert.sh script
	cmd := exec.Command("./generate_cert.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the script and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute generate_cert.sh: %v", err)
	}

	return nil
}
