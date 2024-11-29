package json

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func ReadRequestBody(r *http.Request, result any) error {
	slog.Debug("request body", result)

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func WriteResponseBody(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
