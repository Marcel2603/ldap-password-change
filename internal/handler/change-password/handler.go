package change_password

import (
	"fmt"
	"ldap-password-change/internal/handler/ldap"
	"ldap-password-change/views"
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
	ldap.Handler()
	templ := views.SuccessfulPasswordChange()
	templ.Render(r.Context(), w)
}

func getUserInformation(r *http.Request) userInformation {
	return userInformation{
		username:        r.FormValue("username"),
		currentPassword: r.FormValue("current-password"),
		newPassword:     r.FormValue("new-password"),
		confirmPassword: r.FormValue("confirm-password"),
	}
}
