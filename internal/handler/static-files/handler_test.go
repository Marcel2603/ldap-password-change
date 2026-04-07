package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func init() {
	NewHandler(fstest.MapFS{
		"static/test.txt": &fstest.MapFile{Data: []byte("embedded content")},
	})
}

func TestHandler(t *testing.T) {
	t.Parallel()
	t.Run("Embedded static file — no encoding", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/test.txt", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
		if rr.Body.String() != "embedded content" {
			t.Errorf("Expected 'embedded content', got %q", rr.Body.String())
		}
	})

	t.Run("Embedded static file — br encoding", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/test.txt", nil)
		req.Header.Add("Accept-Encoding", "br")
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
		if rr.Header().Get("content-encoding") != "br" {
			t.Errorf("Expected br encoding")
		}
	})

	t.Run("Embedded static file — gzip encoding", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/test.txt", nil)
		req.Header.Add("Accept-Encoding", "gzip")
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
		if rr.Header().Get("content-encoding") != "gzip" {
			t.Errorf("Expected gzip encoding")
		}
	})

	t.Run("Embedded static file — not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/nonexistent.xyz", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %v", rr.Code)
		}
	})

	t.Run("Custom file from filesystem", func(t *testing.T) {
		tmpDir := t.TempDir()
		_ = os.Chdir(tmpDir)
		_ = os.MkdirAll("custom", 0755)
		_ = os.WriteFile(filepath.Join("custom", "test.txt"), []byte("custom content"), 0644)

		req := httptest.NewRequest("GET", "/custom/test.txt", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
		if rr.Body.String() != "custom content" {
			t.Errorf("Expected 'custom content', got %v", rr.Body.String())
		}
	})

	t.Run("Custom file — not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/custom/nonexistent.txt", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %v", rr.Code)
		}
	})

	t.Run("Directory path — not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/custom", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %v", rr.Code)
		}
	})
}

func TestHandleFavicon(t *testing.T) {
	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	rr := httptest.NewRecorder()
	HandleFavicon(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("Expected 301, got %v", rr.Code)
	}
	if rr.Header().Get("Location") != "/static/favicon.ico" {
		t.Errorf("Expected location /static/favicon.ico, got %v", rr.Header().Get("Location"))
	}
}
