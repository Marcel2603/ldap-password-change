package index

import (
	"ldap-password-change/views"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := views.Index()
	component.Render(r.Context(), w)
}
