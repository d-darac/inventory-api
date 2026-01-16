package inventories

import (
	"github.com/d-darac/inventory-assets/database"
)

type CreateInventoryParams struct {
	InStock   int32  `json:"in_stock" validate:"required"`
	Item      string `json:"item" validate:"required,uuid"`
	Orderable *int32 `json:"orderable" validate:"omitnil"`
}

type ListInventoriesParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
	InStock   *int32              `json:"in_stock" validate:"omitnil"`
	Orderable *int32              `json:"orderable" validate:"omitnil"`
	Reserved  *int32              `json:"reserved" validate:"omitnil"`
}

type ListItemsParams struct {
	*database.PaginationParams
	Active        *bool               `json:"active" validate:"omitnil"`
	CreatedAt     *database.TimeRange `json:"created_at" validate:"omitnil"`
	Description   *string             `json:"description" validate:"omitnil"`
	Group         *string             `json:"group" validate:"omitnil,uuid"`
	Name          *string             `json:"name" validate:"omitnil"`
	PriceAmount   *int32              `json:"price_amount" validate:"omitnil"`
	PriceCurrency *database.Currency  `json:"price_currency" validate:"omitnil,currency"`
	Type          *database.ItemType  `json:"type" validate:"omitnil,itemtype"`
	UpdatedAt     *database.TimeRange `json:"updated_at" validate:"omitnil"`
	Variant       *bool               `json:"variant" validate:"omitnil"`
	Expand        *[]string           `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
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

func NewListItemsParams() ListItemsParams {
	limit := int32(10)
	return ListItemsParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
