package changepassword

import (
	"errors"
	"fmt"
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
		userInfo := getUserInformation(r)
		validationError := validateUserInfo(validator, userInfo)
		if validationError != nil {
			renderErrorToastie(w, r, http.StatusBadRequest, "Some input was not valid", validationError)
			return
		}

		changePasswordError := ldapService.ChangePassword(userInfo.username, userInfo.currentPassword, userInfo.newPassword)
		if changePasswordError != nil {
			log.Printf("Could not change password: %s\n", changePasswordError.Error())
			renderErrorToastie(w, r, http.StatusInternalServerError, "Failed to change password", changePasswordError)
			return
		}

		templ := views.SuccessfulPasswordChange()
		logRenderError(templ.Render(r.Context(), w))
	}
}

func getUserInformation(r *http.Request) *userInformation {
	return &userInformation{
		username:        r.FormValue("username"),
		currentPassword: r.FormValue("current-password"),
		newPassword:     r.FormValue("new-password"),
		confirmPassword: r.FormValue("confirm-password"),
	}
}

func validateUserInfo(validator validation.Validator, userInfo *userInformation) error {
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

func renderErrorToastie(w http.ResponseWriter, r *http.Request, statusCode int, errorTitle string, err error) {
	toast := views.ErrorToastie(fmt.Sprintf("%s: %s", errorTitle, err.Error()))
	w.WriteHeader(statusCode)
	logRenderError(toast.Render(r.Context(), w))
}

func logRenderError(err error) {
	if err != nil {
		log.Printf("Could not render: %s\n", err.Error())
	}
}
