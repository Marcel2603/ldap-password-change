package validate

import (
	"ldap-password-change/views"
	"net/http"
	"regexp"
)

type validationError struct {
	ErrorMsg string
}

func (v validationError) Error() string {
	return v.ErrorMsg
}

var (
	ALLOWED_PASSWORD_CHARS = "[a-zA-Z0-9!@#$%^&*]+$"
	MIN_PASSWORD_LENGTH    = 8
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	switch r.URL.Path {
	case "/validate/username":
		validateUsername(w, r)
	case "/validate/password":
		validatePassword(w, r)
	case "/validate/confirm-password":
		validateConfirmPassword(w, r)
	default:
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func validateConfirmPassword(w http.ResponseWriter, r *http.Request) {
	newPassword := r.FormValue("new-password")
	confirmPassword := r.FormValue("confirm-password")
	if newPassword == confirmPassword {
		component := views.ConfirmPassword(confirmPassword, nil)
		component.Render(r.Context(), w)
	} else {
		errorMsg := "Passwords do not match"
		component := views.ConfirmPassword(confirmPassword, validationError{errorMsg})
		component.Render(r.Context(), w)
	}
}

func validatePassword(w http.ResponseWriter, r *http.Request) {
	newPassword := r.FormValue("new-password")
	valid := regexp.MustCompile(ALLOWED_PASSWORD_CHARS).MatchString(newPassword)
	hasMinLenght := len(newPassword) >= MIN_PASSWORD_LENGTH
	if hasMinLenght && valid {
		component := views.NewPassword(newPassword, nil)
		component.Render(r.Context(), w)
	} else {
		errorMsg := "Password must be at least 8 characters long and contain only letters, numbers, and special characters"
		component := views.NewPassword(newPassword, validationError{errorMsg})
		component.Render(r.Context(), w)
	}
}

func validateUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	valid := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(username)
	if len(username) >= MIN_PASSWORD_LENGTH && valid {
		component := views.UsernameForm(username, nil)
		component.Render(r.Context(), w)
	} else {
		component := views.UsernameForm(username, validationError{"Username must be at least 8 characters long and contain only letters"})
		component.Render(r.Context(), w)
	}
}
