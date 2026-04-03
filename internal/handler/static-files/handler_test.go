package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandler(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	_ = os.MkdirAll("static", 0755)
	_ = os.WriteFile(filepath.Join("static", "test.txt"), []byte("test content"), 0644)

	t.Run("File exists no encoding", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/test.txt", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %v", rr.Code)
		}
	})

	t.Run("File exists br encoding", func(t *testing.T) {
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

	t.Run("File exists gzip encoding", func(t *testing.T) {
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

	t.Run("Directory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static", nil)
		rr := httptest.NewRecorder()
		Handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %v", rr.Code)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/notfound.txt", nil)
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
