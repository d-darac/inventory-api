package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message"`
	Param   string    `json:"param,omitempty"`
	Type    ErrorType `json:"type"`
}

type ErrorListResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ListResponse struct {
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	Url     string        `json:"url"`
}

func (e ErrorResponse) ResError(w http.ResponseWriter, statCode int, err error) {
	if err != nil {
		log.Println(err)
	}
	if statCode >= http.StatusInternalServerError {
		log.Printf("responding with status %d", statCode)
	}
	ResJSON(w, statCode, struct {
		Error ErrorResponse `json:"error"`
	}{Error: e})
}

func (e ErrorListResponse) ResError(w http.ResponseWriter, statCode int, err error) {
	if err != nil {
		log.Println(err)
	}
	if statCode >= http.StatusInternalServerError {
		log.Printf("responding with status %d", statCode)
	}
	ResJSON(w, statCode, e)
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
		Error ErrorResponse `json:"error"`
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
	if statCode != http.StatusNoContent {
		w.Write([]byte(data))
	}
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
