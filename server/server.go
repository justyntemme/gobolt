package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/justyntemme/gobolt/dom"
	"github.com/sirupsen/logrus"
)

// ServerConfig holds configuration values for the server.
type ServerConfig struct {
	BaseDir  string // The base directory to serve content from
	Logger   *logrus.Logger
	DOM      *dom.DOM
	Hostname string
}

type Server struct {
	ServerConfig
	httpServer *http.Server
	mux        *http.ServeMux
}

func NewServer(port string, dom *dom.DOM) *Server {
	mux := http.NewServeMux()
	// TODO add config package to read yaml files or params for ServerConfig Values
	return &Server{
		ServerConfig: ServerConfig{
			BaseDir:  "./content",
			Logger:   logrus.New(),
			DOM:      dom,
			Hostname: "localhost",
		},
		mux: mux,
		httpServer: &http.Server{
			Addr:    port,
			Handler: mux,
		},
	}
}

// registerRoutes sets up the routes and their handlers.
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.handleContent(w, r) // Call the method using the receiver
	})

	// Example: Add a custom route for serving CSS
	s.mux.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		s.handleCSS(w, r) // Serve CSS content
	})
}

func (s *Server) getSafeFilePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)

	absPath := filepath.Join(s.ServerConfig.BaseDir, cleanPath)

	absBaseDir, err := filepath.Abs(s.ServerConfig.BaseDir)
	if err != nil {
		return "", fmt.Errorf("Error resolving base directory")
	}
	absFilePath, err := filepath.Abs(absPath)
	if err != nil {
		return "", fmt.Errorf("Error resolving requested file path")
	}
	if !strings.HasPrefix(absFilePath, absBaseDir) {
		return "", fmt.Errorf("Forbidden: Access outside of the base directory is not allowed")
	}

	return absFilePath, nil
}

func (s *Server) handleCSS(w http.ResponseWriter, r *http.Request) {
	CSS := dom.GetThemeCSS()
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, CSS)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// writeCSSImport dynamically writes the CSS import statement to the provided writer.
func writeCSSImport(w io.Writer, hostname string) error {
	_, err := fmt.Fprintf(w, `<html><head><link rel="stylesheet" type="text/css" href="http://%s/css"></head><body>`, hostname)
	return err
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	s.Logger.Info("Recieved request at URI: ", r.URL)
	writeCSSImport(w, "localhost")
	path := strings.TrimPrefix(r.URL.Path, "/content/`")

	filePath, err := s.getSafeFilePath(path)
	if err != nil {
		s.Logger.Warn("Error with request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uri := "content" + r.URL.Path
	page, exists := s.ServerConfig.DOM.Pages[uri]
	if !exists {
		s.Logger.Warn("File path not found for request with uri: ", uri)
		for uri, page := range s.ServerConfig.DOM.Pages {
			s.Logger.Info("Found Page: ")
			s.Logger.Infof("URI: %s", uri)
			s.Logger.Infof("Markdown: %s", page.Markdown)
			s.Logger.Infof("HTML: %s", page.HTML)
		}
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, page.HTML)
	s.Logger.Debug(page.HTML)
	fmt.Fprintln(w, "</body></html>")
	duration := time.Since(startTime)
	s.ServerConfig.Logger.Infof(
		"Request processed in %s for path: %s with filepath %s",
		duration,
		r.URL.Path,
		filePath)
}

func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()

	fmt.Printf("Server listening on http://localhost%s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// Set a timeout for graceful shutdown
const shutdownTimeout = 5 * time.Second // seconds
