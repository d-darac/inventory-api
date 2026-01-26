package inventories

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateInventoryParams struct {
	InStock   int32  `json:"in_stock" validate:"required"`
	Orderable *int32 `json:"orderable" validate:"omitnil"`
}

type ListInventoriesByIdsParams struct {
	Ids []uuid.UUID
}

type ListInventoriesParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
	InStock   *int32              `json:"in_stock" validate:"omitnil"`
	Orderable *int32              `json:"orderable" validate:"omitnil"`
	Reserved  *int32              `json:"reserved" validate:"omitnil"`
}

type RetrieveInventoryParams struct {
}

type UpdateInventoryParams struct {
	InStock   *int32 `json:"in_stock" validate:"omitnil"`
	Orderable *int32 `json:"orderable" validate:"omitnil"`
}

func NewListInventoriesParams() ListInventoriesParams {
	limit := int32(10)
	return ListInventoriesParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
