package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/d-darac/inventory-api/internal/str"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message"`
	Type    ErrorType `json:"type"`
}

type ListResponse struct {
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	Url     string        `json:"url"`
}

type GroupResponse struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Description str.NullString `json:"description"`
	Name        string         `json:"name"`
	ParentGroup uuid.NullUUID  `json:"parent_group"`
}

func ResError(
	w http.ResponseWriter,
	statCode int,
	msg string,
	errType ErrorType,
	errCode *ErrorCode,
	err error,
) {
	if err != nil {
		log.Println(err)
	}
	if statCode >= http.StatusInternalServerError {
		log.Printf("responding with status %d", statCode)
	}
	errRes := ErrorResponse{
		Message: msg,
		Type:    errType,
	}
	if errCode != nil {
		errRes.Code = *errCode
	}
	ResJSON(w, statCode, struct {
		Error interface{} `json:"error"`
	}{Error: errRes})
}

func ResJSON(w http.ResponseWriter, statCode int, payload interface{}) {
	setDefaultHeaders(w)
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statCode)
	w.Write([]byte(data))
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
