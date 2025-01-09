package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/justyntemme/gobolt/dom"
	"github.com/justyntemme/gobolt/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.InfoLevel)
	baseDir := "./content" // TODO make this an argument to the binary

	logger.Info("Starting Loader")
	domInstance := dom.NewDOM()

	// Load the Markdown content into the DOM (pass in the base directory)
	err := domInstance.LoadMarkdown(baseDir)
	if err != nil {
		log.Fatalf("Error loading markdown content: %v", err)
	}

	logger.Info("DOM loaded successfully.")
	logger.Info("Creating new Server")
	// Start the server in a separate goroutine
	srv := server.NewServer(":80", domInstance)
	srv.Hostname = "localhost"
	BaseDir := srv.BaseDir
	dom.LoadCSS(BaseDir + "/styles.css")
	logger.Info("Loaded CSS at " + BaseDir + "/styles.css")
	css := dom.GetThemeCSS()
	logger.Debug(css)
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
