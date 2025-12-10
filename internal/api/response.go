package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func resError(w http.ResponseWriter, statCode int, msg string, errType ErrorType, errCode *ErrorCode, err error) {
	if err != nil {
		log.Println(err)
	}

	if statCode >= http.StatusInternalServerError {
		log.Printf("responding with %d error: %s", statCode, msg)
	}

	errRes := ErrResponse{
		Message: msg,
		Type:    errType,
	}
	if errCode != nil {
		errRes.Code = *errCode
	}

	resJSON(w, statCode, struct {
		Error ErrResponse `json:"error"`
	}{Error: errRes})
}

func resJSON(w http.ResponseWriter, statCode int, payload interface{}) {
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
