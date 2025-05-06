package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

const (
	internalServerError = "internal server error"
)

func RespondWithJSON(w http.ResponseWriter, log *slog.Logger, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Error(fmt.Errorf("response marshalling error: %w", err).Error())

		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		w.Write([]byte(internalServerError))

		return
	}

	w.WriteHeader(code)
	//nolint:errcheck
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, log *slog.Logger, code int, message string) {
	RespondWithJSON(w, log, code, ErrorResponse{Errors: message})
}
