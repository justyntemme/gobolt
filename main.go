package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/justyntemme/gobolt/server"
)

func main() {
	// Start the server in a separate goroutine
	srv := server.NewServer(":80")
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until a signal is received

	fmt.Println("\nShutting down server...")
	if err := srv.Shutdown(); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}
}
