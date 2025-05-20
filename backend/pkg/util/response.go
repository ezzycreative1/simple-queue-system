package util

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	resp := map[string]interface{}{
		"status":  "success",
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	WriteJSONResponse(w, http.StatusOK, resp)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	resp := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	WriteJSONResponse(w, statusCode, resp)
}
