package static

import (
	"compress/gzip"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
)

var staticFS fs.FS

func NewHandler(root fs.FS) {
	sub, err := fs.Sub(root, "static")
	if err != nil {
		log.Fatalf("static: failed to sub embedded FS: %v", err)
	}
	staticFS = sub
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	if strings.HasPrefix(urlPath, "/static/") {
		serveEmbedded(w, r, urlPath)
		return
	}

	serveFilesystem(w, r, urlPath)
}

func serveEmbedded(w http.ResponseWriter, r *http.Request, urlPath string) {
	name := strings.TrimPrefix(urlPath, "/static/")

	data, err := fs.ReadFile(staticFS, name)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("cache-control", "public, max-age=36000")
	writeBytes(w, r, data)
}

func serveFilesystem(w http.ResponseWriter, r *http.Request, urlPath string) {
	wd, _ := os.Getwd()
	relPath := strings.TrimPrefix(filepath.Clean(urlPath), "/")
	fullPath := filepath.Join(wd, relPath)

	isDir, err := isDirectory(fullPath)
	if err != nil || isDir {
		log.Printf("stat %s: no such file or directory\n", fullPath)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("cache-control", "public, max-age=36000")
	writeBytes(w, r, data)
}

func writeBytes(w http.ResponseWriter, r *http.Request, data []byte) {
	encodings := r.Header.Get("accept-encoding")
	switch {
	case strings.Contains(encodings, "br"):
		w.Header().Set("content-encoding", "br")
		writer := brotli.NewWriter(w)
		defer writer.Close()
		_, _ = writer.Write(data)
	case strings.Contains(encodings, "gzip"):
		w.Header().Set("content-encoding", "gzip")
		writer := gzip.NewWriter(w)
		defer writer.Close()
		_, _ = writer.Write(data)
		writer.Flush()
	default:
		_, _ = w.Write(data)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func HandleFavicon(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("cache-control", "public, max-age=7776000")
	http.Redirect(writer, request, "/static/favicon.ico", http.StatusMovedPermanently)
}
