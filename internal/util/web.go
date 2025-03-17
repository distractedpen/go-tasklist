package util

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, status int, data any) {
	response := make(map[string]any)
	w.Header().Set("Content-Type", "application/json")
	if status >= 400 && status <= 499 {
		w.WriteHeader(status)
		response["status"] = "error"
		response["message"] = "Error Parsing Request"
	} else if status >= 500 && status <= 599 {
		w.WriteHeader(status)
		response["status"] = "error"
		response["message"] = "Server Error"
	} else if status != 200 {
		w.WriteHeader(status)
		response["status"] = "success"
	} else {
		response["status"] = "success"
		response["data"] = data
	}
	json.NewEncoder(w).Encode(response)
}
