package groups

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateGroupParams struct {
	Description *string    `json:"description" validate:"omitnil"`
	Name        string     `json:"name" validate:"required"`
	ParentGroup *uuid.UUID `json:"parent_group" validate:"omitnil,uuid"`
	Expand      *[]string  `json:"expand" validate:"omitnil,dive,oneof=parent_group"`
}

type ListGroupsParams struct {
	*database.PaginationParams
	CreatedAt   *database.TimeRange `json:"created_at" validate:"omitnil"`
	ParentGroup *uuid.UUID          `json:"parent_group" validate:"omitnil,uuid"`
	Description *string             `json:"description" validate:"omitnil"`
	Name        *string             `json:"name" validate:"omitnil"`
	UpdatedAt   *database.TimeRange `json:"updated_at" validate:"omitnil"`
}

type RetrieveGroupParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=parent_group"`
}

type UpdateGroupParams struct {
	Description *string    `json:"description" validate:"omitnil"`
	Name        *string    `json:"name" validate:"omitnil"`
	ParentGroup *uuid.UUID `json:"parent_group" validate:"omitnil,uuid"`
	Expand      *[]string  `json:"expand" validate:"omitnil,dive,oneof=parent_group"`
}

func NewListGroupParams() *ListGroupsParams {
	limit := int32(10)
	return &ListGroupsParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
