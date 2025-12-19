package groups

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type createParams struct {
	Description *string    `json:"description" validate:"omitnil,lte=1024"`
	Name        string     `json:"name" validate:"required,lte=64"`
	ParentGroup *uuid.UUID `json:"parent_group"`
	expandParam
}

type listParams struct {
	*database.PaginationParams
	CreatedAt   *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt   *database.TimeRange `json:"updated_at" validate:"omitnil"`
	ParentGroup *uuid.UUID          `json:"parent_group"`
	Description *string             `json:"description" validate:"omitnil,lte=1024"`
	Name        *string             `json:"name" validate:"omitnil,lte=64"`
}

type expandParam struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=parent_group"`
}

type updateParams struct {
	Description *string    `json:"description" validate:"omitnil,lte=1024"`
	Name        *string    `json:"name" validate:"omitnil,lte=64"`
	ParentGroup *uuid.UUID `json:"parent_group" validate:"omitnil"`
	expandParam
}

func newListParams() *listParams {
	limit := int32(int(10))
	return &listParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
