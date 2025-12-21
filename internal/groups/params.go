package groups

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type createParams struct {
	*expandParam
	Description *string    `json:"description" validate:"omitnil"`
	Name        string     `json:"name" validate:"required"`
	ParentGroup *uuid.UUID `json:"parent_group"`
}

type listParams struct {
	*database.PaginationParams
	CreatedAt   *database.TimeRange `json:"created_at" validate:"omitnil"`
	ParentGroup *uuid.UUID          `json:"parent_group"`
	Description *string             `json:"description" validate:"omitnil"`
	Name        *string             `json:"name" validate:"omitnil"`
	UpdatedAt   *database.TimeRange `json:"updated_at" validate:"omitnil"`
}

type expandParam struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=parent_group"`
}

type updateParams struct {
	*expandParam
	Description *string    `json:"description" validate:"omitnil"`
	Name        *string    `json:"name" validate:"omitnil"`
	ParentGroup *uuid.UUID `json:"parent_group" validate:"omitnil"`
}

func newListParams() *listParams {
	limit := int32(int(10))
	return &listParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
