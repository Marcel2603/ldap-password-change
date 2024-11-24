package static_files

import (
	"fmt"
	"net/http"
	"os"
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
	w.Header().Set("cache-control", "public, max-age=3600")
	http.ServeFile(w, r, staticPath)
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
