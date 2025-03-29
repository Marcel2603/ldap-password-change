package change_password

import (
	"errors"
	"ldap-password-change/internal/service/ldap"
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

func Handler(ldapService ldap.Service, validator validation.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)
		userInfo := getUserInformation(r)
		validationError := validateUserInfo(validator, userInfo)
		if validationError != nil {
			toast := views.ErrorToastie("Some input was not valid " + validationError.Error())
			w.WriteHeader(http.StatusBadRequest)
			toast.Render(r.Context(), w)
			return
		}
		err := ldapService.ChangePassword(userInfo.username, userInfo.currentPassword, userInfo.newPassword)
		if err != nil {
			log.Println(err)
			toast := views.ErrorToastie("Failed to change password")
			w.WriteHeader(http.StatusInternalServerError)
			toast.Render(r.Context(), w)
			return
		}
		templ := views.SuccessfulPasswordChange()
		templ.Render(r.Context(), w)
	}
}

func getUserInformation(r *http.Request) userInformation {
	return userInformation{
		username:        r.FormValue("username"),
		currentPassword: r.FormValue("current-password"),
		newPassword:     r.FormValue("new-password"),
		confirmPassword: r.FormValue("confirm-password"),
	}
}

func validateUserInfo(validator validation.Validator, userInfo userInformation) error {
	validUsername := validator.ValidateUsername(userInfo.username)
	if !validUsername {
		return errors.New("invalid username")
	}
	validPassword := validator.ValidatePassword(userInfo.newPassword)
	if !validPassword {
		return errors.New("invalid password")
	}
	if userInfo.newPassword != userInfo.confirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}
