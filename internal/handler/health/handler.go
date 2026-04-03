package health

import (
	"encoding/json"
	"ldap-password-change/internal/service/ldap"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func LivenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{Status: "ok"})
}

func ReadinessHandler(svc ldap.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := svc.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(Response{Status: "unavailable", Message: "ldap unreachable"})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{Status: "ok"})
	}
}
