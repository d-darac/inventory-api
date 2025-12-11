package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResError(w http.ResponseWriter, statCode int, errRes interface{}, err error) {
	if err != nil {
		log.Println(err)
	}

	if statCode >= http.StatusInternalServerError {
		log.Printf("responding with status %d", statCode)
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
