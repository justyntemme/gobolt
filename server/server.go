package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/justyntemme/gobolt/dom"
	"github.com/sirupsen/logrus"
)

// ServerConfig holds configuration values for the server.
type ServerConfig struct {
	BaseDir string // The base directory to serve content from
	Logger  *logrus.Logger
	DOM     *dom.DOM
}

type Server struct {
	ServerConfig
	httpServer *http.Server
}

func NewServer(port string, dom *dom.DOM) *Server {
	return &Server{
		ServerConfig: ServerConfig{
			BaseDir: "./content",
			Logger:  logrus.New(),
			DOM:     dom,
		},
		httpServer: &http.Server{
			Addr: port,
		},
	}
}

// registerRoutes sets up the routes and their handlers.
func (s *Server) registerRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.handleContent(w, r) // Call the method using the receiver
	})
	// http.HandleFunc("/search", handleSearch)
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

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	s.Logger.Info("Recieved request at URI: ", r.URL)
	path := strings.TrimPrefix(r.URL.Path, "/content/`")

	filePath, err := s.getSafeFilePath(path)
	if err != nil {
		s.Logger.Warn("Error with request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO implement the ability to correctly grab the right content markdown
	// INFO[0002] Recieved request at URI: /home
	// INFO[0002] map[content/about:0x140001404a0 content/home:0x140001404c0]
	// WARN[0002] File path not found for request: /Users/justyntemme/Documents/code/gobolt/content/home
	// uri := strings.TrimPrefix(r.URL.Path, "/content")
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
	duration := time.Since(startTime)
	s.ServerConfig.Logger.Info(
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
const shutdownTimeout = 5 // seconds
