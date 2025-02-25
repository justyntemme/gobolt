package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/justyntemme/gobolt/dom"
	t "github.com/justyntemme/gobolt/template"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Global variable to hold the generated navigation HTML
var (
	navigationHTML string
	once           sync.Once
)

// NavData holds navigation link information
type NavData struct {
	Title string
	URI   string
}

func (s *Server) isTopLevel(uri string) bool {
	sp := strings.Split(uri, "/")
	// pop first element as its empty str
	if sp[0] == "" {
		sp = sp[1:]
	}
	s.Logger.Debugf("isTopLevel: uri is %s, and len is %d", sp, len(sp))
	return len(sp) == 1
}

// GenerateNavigationHTML dynamically generates the navigation bar based on site content.
func (s *Server) GenerateNavigationHTML() error {
	var loadErr error

	once.Do(func() {
		// Generate navigation links for all top-level pages
		navLinks := []NavData{}
		for uri := range s.DOM.Pages {
			// Consider top-level URIs only (e.g., "/about", not "/about/team")
			// if path.Dir(uri) == "/" || uri == "/"  // TODO Add check for only top level pages
			// by checking if the len after split by '/' is greater than 1
			uri = strings.TrimPrefix(uri, "content")

			topLevel := s.isTopLevel(uri)
			if topLevel {
				title := s.getPageTitle(uri)
				fmt.Print(uri)
				navLinks = append(navLinks, NavData{
					Title: title,
					URI:   uri,
				})

			} else {
				continue
			}

		}

		// Define a simple template for the navigation bar
		// baseTemplate := t.ReturnBaseTemplate()
		navTemplate := t.ReturnNavTemplate()

		// Parse the template
		tmpl, err := template.New("navigation").Parse(navTemplate)
		if err != nil {
			loadErr = fmt.Errorf("failed to parse navigation template: %w", err)
			return
		}

		// Execute the template with the navigation links
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, navLinks); err != nil {
			loadErr = fmt.Errorf("failed to execute navigation template: %w", err)
			return
		}

		// Store the generated HTML
		navigationHTML = buf.String()
		fmt.Println("Navigation HTML generated successfully.")
	})

	return loadErr
}

// getPageTitle derives a page title from the URI (optional utility function)
func (s *Server) getPageTitle(uri string) string {
	if uri == "/" {
		return "Home"
	}
	return cases.Title(language.Und).String(strings.Trim(path.Base(uri), "/"))
}

// ServerConfig holds configuration values for the server.
type ServerConfig struct {
	BaseDir  string // The base directory to serve content from
	DOM      *dom.DOM
	Hostname string
	Port     string
	Logger   *logrus.Logger
}

type Server struct {
	ServerConfig
	httpServer *http.Server
	mux        *http.ServeMux
}

func NewServer(dom *dom.DOM) (*Server, error) {
	mux := http.NewServeMux()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	// TODO add config package to read yaml files or params for ServerConfig Values
	c, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerConfig: ServerConfig{
			BaseDir:  c.BaseDir,
			DOM:      dom,
			Hostname: c.Hostname,
			Port:     c.Port,
			Logger:   logger,
		},
		mux: mux,
		httpServer: &http.Server{
			Addr:         c.Port,
			Handler:      mux,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			IdleTimeout:  time.Second * 5,
		},
	}
	return s, nil
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
	// TODO Implement symlink checking

	absPath := filepath.Join(s.BaseDir, cleanPath)

	absBaseDir, err := filepath.Abs(s.BaseDir)
	if err != nil {
		return "", fmt.Errorf("error resolving base directory")
	}
	absFilePath, err := filepath.Abs(absPath)
	if err != nil {
		return "", fmt.Errorf("error resolving requested file path")
	}
	if !strings.HasPrefix(absFilePath, absBaseDir) {
		return "", fmt.Errorf("forbidden: Access outside of the base directory is not allowed")
	}

	return absFilePath, nil
}

func (s *Server) handleCSS(w http.ResponseWriter, _ *http.Request) {
	CSS := dom.GetThemeCSS()
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, CSS)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	s.Logger.Debug("Recieved request at URI: ", r.URL)
	path := strings.TrimPrefix(r.URL.Path, "/"+s.BaseDir+"/`")

	filePath, err := s.getSafeFilePath(path)
	if err != nil {
		s.Logger.Warn("Error with request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uri := "content" + r.URL.Path
	page, exists := s.DOM.Pages[uri]
	if !exists {

		/* s.Logger.Warn("File path not found for request with uri: ", uri)
		for uri, page := range s.ServerConfig.DOM.Pages {
			s.Logger.Info("Found Page: ")
			s.Logger.Infof("URI: %s", uri)
			s.Logger.Infof("Markdown: %s", page.Markdown)
			s.Logger.Infof("HTML: %s", page.HTML)
		}
		*/
		http.NotFound(w, r)
		return
	}
	cssImport := t.ReturnCSSImportTemplate(s.Hostname)

	data := struct {
		CSSImport   template.HTML
		Navigation  template.HTML
		PageContent template.HTML
		Hostname    string
	}{
		CSSImport:   template.HTML(cssImport),
		Navigation:  template.HTML(navigationHTML),
		PageContent: template.HTML(page.HTML),
		Hostname:    s.Hostname,
	}

	tmpl, err := template.New("main").Parse(t.ReturnBaseTemplate())
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	// Set content type and render the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}

	duration := time.Since(startTime)
	s.Logger.Infof(
		"Request processed in %s for path: %s with filepath %s",
		duration,
		r.URL.Path,
		filePath)
}

func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()

	err := s.GenerateNavigationHTML()
	if err != nil {
		return err
	}

	err = dom.LoadCSS(s.BaseDir + "/styles.css")
	if err != nil {
		return err
	}

	fmt.Printf("Server listening on http://%s %s\n", s.Hostname, s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// Set a timeout for graceful shutdown
const shutdownTimeout = 5 * time.Second // seconds
