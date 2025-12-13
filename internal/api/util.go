package api

import (
	"net/http"

	"github.com/google/uuid"
)

func GetIdFromPath(r *http.Request) (uuid.UUID, string) {
	pathValue := r.PathValue("id")
	id, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.UUID{}, InvalidIdMessage(pathValue, "account")
	}
	return id, ""
}
