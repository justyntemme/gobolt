package server

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// NewServer creates a new Server instance.
func NewServer(port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: port,
		},
	}
}

func (s *Server) Start() error {
	// Register routes
	registerRoutes()

	fmt.Printf("Server listening on http://localhost%s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// Set a timeout for graceful shutdown
const shutdownTimeout = 5 // seconds
