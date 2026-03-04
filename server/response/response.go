package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Status:  status,
		Message: message,
	})
}

func InvalidJSON(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, "Invalid JSON")
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Internal server error")
}

func ValidationError(w http.ResponseWriter, details map[string]string) {
	JSON(w, http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Details: details,
	})
}

func IDNotFound(w http.ResponseWriter, id string) {
	errString := fmt.Sprintf("Can`t find item with id: %s", id)
	Error(w, http.StatusNotFound, errString)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, SuccessResponse{
		Status:  status,
		Message: "success",
		Data:    data,
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, data)
}

func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data)
}
