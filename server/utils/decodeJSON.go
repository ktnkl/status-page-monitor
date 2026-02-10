package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"status-page-monitor/server/response"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		response.InvalidJSON(w)
		log.Println("json decoding error:", err)
		return false
	}

	return true
}
