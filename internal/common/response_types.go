package common

import (
	"time"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type ErrResponse struct {
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message"`
	Type    ErrorType `json:"type"`
}

type ListResponse struct {
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	Url     string        `json:"url"`
}

type AccountResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Country   database.Country `json:"country"`
	Nickname  string           `json:"nickname"`
}

func NewErrorResponse(msg string, errType ErrorType, errCode *ErrorCode) *ErrResponse {
	errRes := ErrResponse{
		Message: msg,
		Type:    errType,
	}
	if errCode != nil {
		errRes.Code = *errCode
	}
	return &errRes
}
