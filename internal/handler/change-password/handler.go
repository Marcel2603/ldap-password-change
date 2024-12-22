package change_password

import (
	"errors"
	"ldap-password-change/internal/validation"
	"ldap-password-change/views"
	"log"
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
	log.Println(r.Body)
	userInfo := getUserInformation(r)
	log.Println(userInfo)
	validationError := validateUserInfo(userInfo)
	if validationError != nil {
		http.Error(w, validationError.Error(), http.StatusBadRequest)
		return
	}
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

func validateUserInfo(userInfo userInformation) error {
	validUsername := validation.ValidateUsername(userInfo.username)
	if !validUsername {
		return errors.New("invalid username")
	}
	validPassword := validation.ValidatePassword(userInfo.newPassword)
	if !validPassword {
		return errors.New("invalid password")
	}
	if userInfo.newPassword != userInfo.confirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}
