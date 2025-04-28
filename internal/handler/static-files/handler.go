package static

import (
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"net/http"
	"os"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()

	staticPath := wd + r.URL.Path
	isDir, err := isDirectory(staticPath)
	if err != nil || isDir {
		fmt.Println(err)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("cache-control", "public, max-age=36000")
	writeFile(w, r, staticPath)
	writer := brotli.NewWriter(w)
	defer writer.Close()
	w.Header().Set("content-encoding", "br")
	file, _ := os.ReadFile(staticPath)
	writer.Write(file)
	writer.Flush()
}

func writeFile(w http.ResponseWriter, r *http.Request, path string) {
	encodings := r.Header.Get("accept-encoding")
	switch {
	case strings.Contains(encodings, "br"):
		w.Header().Set("content-encoding", "br")
		writer := brotli.NewWriter(w)
		defer writer.Close()
		file, _ := os.ReadFile(path)
		writer.Write(file)
		writer.Flush()
	case strings.Contains(encodings, "gzip"):
		w.Header().Set("content-encoding", "gzip")
		writer := gzip.NewWriter(w)
		defer writer.Close()
		file, _ := os.ReadFile(path)
		writer.Write(file)
		writer.Flush()
	default:
		file, _ := os.ReadFile(path)
		w.Write(file)
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
