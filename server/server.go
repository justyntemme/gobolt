package server

import (
	"fmt"
	"net/http"
)

// Start initializes and starts the HTTP server.
func Start(port string) error {
	// Register routes
	registerRoutes()

	// Start the HTTP server
	fmt.Printf("Server listening on http://localhost%s\n", port)
	return http.ListenAndServe(port, nil)
}
