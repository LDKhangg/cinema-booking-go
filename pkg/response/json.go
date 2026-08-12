package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func FromError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.Status, map[string]string{
			"error": appErr.Message,
			"code":  appErr.Code,
		})
		return
	}

	JSON(w, http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
		"code":  "internal_error",
	})
}

func Success(w http.ResponseWriter, status int, message string, data any) {
	payload := map[string]any{}
	if message != "" {
		payload["message"] = message
	}
	if data != nil {
		payload["data"] = data
	} else if message == "" {
		payload = nil
	}
	JSON(w, status, payload)
}
