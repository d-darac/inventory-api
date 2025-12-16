package groups

import (
	"time"

	"github.com/google/uuid"
)

type createParams struct {
	Description *string    `json:"description" validate:"lte=1024"`
	Name        string     `json:"name" validate:"required,lte=64"`
	ParentGroup *uuid.UUID `json:"parent_group"`
}

type listParams struct {
	CreatedAt struct {
		Gt  *time.Time `json:"gt"`
		Lt  *time.Time `json:"lt"`
		Gte *time.Time `json:"gte"`
		Lte *time.Time `json:"lte"`
	} `json:"created_at"`
	UpdatedAt struct {
		Gt  *time.Time `json:"gt"`
		Lt  *time.Time `json:"lt"`
		Gte *time.Time `json:"gte"`
		Lte *time.Time `json:"lte"`
	} `json:"updated_at"`
	StartingAfter *uuid.UUID `json:"starting_after" validate:"excluded_with=EndingBefore"`
	EndingBefore  *uuid.UUID `json:"ending_before" validate:"excluded_with=StartingAfter"`
	ParentGroup   *uuid.UUID `json:"parent_group"`
	Description   *string    `json:"description" validate:"lte=1024"`
	Name          *string    `json:"name" validate:"lte=64"`
	Limit         *int       `json:"limit" validate:"gt=0,lte=100"`
}
