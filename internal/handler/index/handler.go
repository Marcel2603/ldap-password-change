package index

import (
	"net/http"

	"github.com/Marcel2603/ldap-password-change/cmd/config"
	"github.com/Marcel2603/ldap-password-change/views"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := views.Index(config.Configuration)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
