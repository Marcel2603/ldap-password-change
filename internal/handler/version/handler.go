package version

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func Handler(version string, commit string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{Version: version, Commit: commit})
	}
}
