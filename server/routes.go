package server

import (
	"fmt"
	"net/http"
)

// registerRoutes sets up the routes and their handlers.
func registerRoutes() {
	http.HandleFunc("/", handleRoot)
	// http.HandleFunc("/search", handleSearch)
}

// handleRoot handles requests to the root endpoint.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<html><body><p>Welcome to the Blog!</p></body></html>`)
}
