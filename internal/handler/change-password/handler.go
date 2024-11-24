package change_password

import (
	"fmt"
	"net/http"
)

type userInformation struct {
	username        string
	currentPassword string
	newPassword     string
	confirmPassword string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	userInfo := getUserInformation(r)
	fmt.Println(userInfo)
	w.Write([]byte("password changed"))
}

func getUserInformation(r *http.Request) userInformation {
	return userInformation{
		username:        r.FormValue("username"),
		currentPassword: r.FormValue("current-password"),
		newPassword:     r.FormValue("new-password"),
		confirmPassword: r.FormValue("confirm-password"),
	}
}
