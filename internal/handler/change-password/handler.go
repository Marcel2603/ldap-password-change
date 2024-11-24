package change_password

import (
  "fmt"
  "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
  r.ParseMultipartForm(0)
  fmt.Println(r.FormValue("username"))
  fmt.Println(r.MultipartForm)
  w.Write([]byte("password changed"))
}
