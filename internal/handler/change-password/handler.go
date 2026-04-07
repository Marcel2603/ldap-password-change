package changepassword

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Marcel2603/ldap-password-change/internal/service/ldap"
	"github.com/Marcel2603/ldap-password-change/internal/validation"
	"github.com/Marcel2603/ldap-password-change/views"

	"github.com/go-chi/chi/v5/middleware"
)

type userInformation struct {
	username        string
	currentPassword string
	newPassword     string
	confirmPassword string
}

type InvalidCredentialsError struct{}

func (m *InvalidCredentialsError) Error() string {
	return "Invalid Credentials!"
}
func Handler(ldapService ldap.Service, validator validation.Validator, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		l := logger.With(slog.String("req_id", reqID), slog.String("class", "handler_change_password"))

		userInfo := getUserInformation(r)
		validationError := validateUserInfo(validator, userInfo)
		if validationError != nil {
			renderErrorToastie(w, r, http.StatusBadRequest, "Some input was not valid", validationError, l)
			return
		}
		user, err := ldapService.SearchUser(userInfo.username)
		if err != nil {
			l.Error("Could not find user", slog.String("error", err.Error()))
			renderErrorToastie(w, r, http.StatusUnauthorized, "Invalid Credentials", &InvalidCredentialsError{}, l)
			return
		}

		changePasswordError := ldapService.ChangePassword(user.DN, userInfo.username, userInfo.currentPassword, userInfo.newPassword)
		if changePasswordError != nil {
			l.Error("Could not change password", slog.String("error", changePasswordError.Error()))
			renderErrorToastie(w, r, http.StatusInternalServerError, "Failed to change password", &InvalidCredentialsError{}, l)
			return
		}

		templ := views.SuccessfulPasswordChange()
		logRenderError(templ.Render(r.Context(), w), l)
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

func renderErrorToastie(w http.ResponseWriter, r *http.Request, statusCode int, errorTitle string, err error, l *slog.Logger) {
	toast := views.ErrorToastie(fmt.Sprintf("%s: %s", errorTitle, err.Error()))
	w.WriteHeader(statusCode)
	logRenderError(toast.Render(r.Context(), w), l)
}

func logRenderError(err error, l *slog.Logger) {
	if err != nil {
		l.Error("Could not render", slog.String("error", err.Error()))
	}
}
