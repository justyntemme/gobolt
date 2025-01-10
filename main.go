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
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{FullTimestamp: true},
		Level:     logrus.DebugLevel,
	}

	domInstance := dom.NewDOM()
	srv, err := server.NewServer(domInstance)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Load the Markdown content into the DOM (pass in the base directory)
	err = domInstance.LoadMarkdown(srv.BaseDir)
	if err != nil {
		log.Fatalf("Error loading markdown content: %v", err)
	}

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until a signal is received

	logger.Info("\nShutting down server...")
	if err := srv.Shutdown(); err != nil {
		logger.Infof("Error during shutdown: %v\n", err)
	}
}
