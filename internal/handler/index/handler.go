package index

import (
	"ldap-password-change/cmd/config"
	"ldap-password-change/views"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := views.Index(config.Configuration)
	component.Render(r.Context(), w)
}
