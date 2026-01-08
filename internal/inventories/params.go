package inventories

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateInventoryParams struct {
	InStock   int32     `json:"in_stock" validate:"required"`
	Item      uuid.UUID `json:"item" validate:"required,uuid"`
	Orderable *int32    `json:"orderable" validate:"omitnil"`
	Expand    *[]string `json:"expand" validate:"omitnil,dive,oneof=items"`
}

type ListInventoriesParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
	InStock   *int32              `json:"in_stock" validate:"omitnil"`
	Orderable *int32              `json:"orderable" validate:"omitnil"`
	Reserved  *int32              `json:"reserved" validate:"omitnil"`
	Expand    *[]string           `json:"expand" validate:"omitnil,dive,oneof=items"`
}

type RetrieveInventoryParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=items"`
}

type UpdateInventoryParams struct {
	InStock   *int32    `json:"in_stock" validate:"omitnil"`
	Orderable *int32    `json:"orderable" validate:"omitnil"`
	Expand    *[]string `json:"expand" validate:"omitnil,dive,oneof=items"`
}

func NewListInventoriesParams() *ListInventoriesParams {
	limit := int32(int(10))
	return &ListInventoriesParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
