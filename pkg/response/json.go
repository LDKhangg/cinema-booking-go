package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	payload := map[string]interface{}{}
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
