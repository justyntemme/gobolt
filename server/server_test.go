// server/server_test.go

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justyntemme/gobolt/dom"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer wraps the server with test utilities
type testServer struct {
	*Server
	tmpDir string
}

// setupTestServer creates a new server instance with a temporary content directory
func setupTestServer(t testing.TB) *testServer {
	// Create temporary directory for test content
	tmpDir, err := os.MkdirTemp("", "gobolt-test-*")
	require.NoError(t, err)

	// Create test DOM instance
	domInstance := dom.NewDOM()

	// Create test logger that discards output during tests
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	srv, err := NewServer(domInstance)
	require.NoError(t, err)

	// Override the base directory for testing
	srv.BaseDir = tmpDir

	return &testServer{
		Server: srv,
		tmpDir: tmpDir,
	}
}

// cleanup removes the temporary test directory
func (ts *testServer) cleanup(t testing.TB) {
	os.RemoveAll(ts.tmpDir)
}

// createTestContent creates a markdown file in the test directory
func (ts *testServer) createTestContent(t testing.TB, path, content string) {
	fullPath := filepath.Join(ts.tmpDir, path)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(fullPath, []byte(content), 0644)
	require.NoError(t, err)
}

func TestHandleContent(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.cleanup(t)

	// Create test content
	ts.createTestContent(t, "content/test.md", "# Test Page\nThis is a test.")

	// Load the markdown content
	err := ts.DOM.LoadMarkdown(ts.BaseDir)
	require.NoError(t, err)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "existing page",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "Test Page",
		},
		{
			name:           "non-existent page",
			path:           "/not-found",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "directory traversal attempt",
			path:           "../../../etc/passwd",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			ts.handleContent(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestSafeFilePath(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.cleanup(t)

	tests := []struct {
		name        string
		path        string
		shouldError bool
	}{
		{
			name:        "valid path",
			path:        "test.md",
			shouldError: false,
		},
		{
			name:        "directory traversal attempt",
			path:        "../../../etc/passwd",
			shouldError: true,
		},
		{
			name:        "absolute path attempt",
			path:        "/etc/passwd",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := ts.getSafeFilePath(tt.path)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, filepath.IsAbs(path))
				assert.True(t, strings.HasPrefix(path, ts.BaseDir))
			}
		})
	}
}

func TestGenerateNavigationHTML(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.cleanup(t)

	// Create test content with multiple pages
	ts.createTestContent(t, "content/page1.md", "# Page 1")
	ts.createTestContent(t, "content/page2.md", "# Page 2")
	ts.createTestContent(t, "content/subfolder/page3.md", "# Page 3")

	// Load the markdown content
	err := ts.DOM.LoadMarkdown(ts.BaseDir)
	fmt.Println(ts.BaseDir)
	require.NoError(t, err)
	for uri := range ts.DOM.Pages {
		// Consider top-level URIs only (e.g., "/about", not "/about/team")
		// if path.Dir(uri) == "/" || uri == "/"  // TODO Add check for only top level pages
		// by checking if the len after split by '/' is greater than 1
		NewUri := strings.TrimPrefix(uri, "/content")
		ts.DOM.Pages[NewUri] = ts.DOM.Pages[uri]
		delete(ts.DOM.Pages, uri)

	}

	// Generate navigation
	err = ts.GenerateNavigationHTML()
	require.NoError(t, err)

	// Verify navigation content
	assert.Contains(t, navigationHTML, `<li><a href="/page1">Page1</a></li>`)
	//assert.Contains(t, navigationHTML, "Page 2")
	// Subfolder page should not be in top-level navigation
	//assert.NotContains(t, navigationHTML, "Page 3")
}

// BenchmarkHandleContent benchmarks the content handler
