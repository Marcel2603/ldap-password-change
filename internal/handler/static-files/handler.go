package static_files

import (
  "net/http"
  "os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
  wd, _ := os.Getwd()
  http.ServeFile(w, r, wd+r.URL.Path)
}
