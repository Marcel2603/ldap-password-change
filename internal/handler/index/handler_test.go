package index

import (
	"ldap-password-change/cmd/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	config.Configuration = config.Config{
		Validation: config.ValidationConfig{
			UsernamePattern: "^[a-zA-Z0-9]+$",
			PasswordPattern: "^.{8,}$",
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
