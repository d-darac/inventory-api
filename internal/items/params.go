package items

import (
	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateInventoryParams struct {
	InStock   int32  `json:"in_stock" validate:"required"`
	Orderable *int32 `json:"orderable" validate:"omitnil"`
}

type CreateItemIdentifiersParams struct {
	Ean  *string `json:"ean" validate:"omitnil"`
	Gtin *string `json:"gtin" validate:"omitnil"`
	Isbn *string `json:"isbn" validate:"omitnil"`
	Jan  *string `json:"jan" validate:"omitnil"`
	Mpn  *string `json:"mpn" validate:"omitnil"`
	Nsn  *string `json:"nsn" validate:"omitnil"`
	Upc  *string `json:"upc" validate:"omitnil"`
	Qr   *string `json:"qr" validate:"omitnil"`
	Sku  *string `json:"sku" validate:"omitnil"`
}

type CreateItemParams struct {
	Description     *string                      `json:"description" validate:"omitnil"`
	Group           *string                      `json:"group" validate:"omitnil,uuid,excluded_with=GroupData"`
	GroupData       *groups.CreateGroupParams    `json:"group_data" validate:"omitnil"`
	IdentifiersData *CreateItemIdentifiersParams `json:"identifiers_data" validate:"omitnil"`
	Inventory       *string                      `json:"inventory" validate:"omitnil,uuid,excluded_with=InventoryData"`
	InventoryData   *CreateInventoryParams       `json:"inventory_data" valdiate:"omitnil"`
	Name            string                       `json:"name" validate:"required"`
	PriceAmount     *int32                       `json:"price_amount" validate:"omitnil"`
	PriceCurrency   *database.Currency           `json:"price_currency" validate:"omitnil,currency"`
	Type            database.ItemType            `json:"type" validate:"required,itemtype"`
	Expand          *[]string                    `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

type ListItemsByIdsParams struct {
	Ids []uuid.UUID
}

type ListItemsParams struct {
	*database.PaginationParams
	Active        *bool               `json:"active" validate:"omitnil"`
	CreatedAt     *database.TimeRange `json:"created_at" validate:"omitnil"`
	Description   *string             `json:"description" validate:"omitnil"`
	Group         *string             `json:"group" validate:"omitnil"`
	Inventory     *string             `json:"inventory" validate:"omitnil,uuid"`
	Name          *string             `json:"name" validate:"omitnil"`
	PriceAmount   *int32              `json:"price_amount" validate:"omitnil"`
	PriceCurrency *database.Currency  `json:"price_currency" validate:"omitnil,currency"`
	Type          *database.ItemType  `json:"type" validate:"omitnil,itemtype"`
	UpdatedAt     *database.TimeRange `json:"updated_at" validate:"omitnil"`
	Variant       *bool               `json:"variant" validate:"omitnil"`
	Expand        *[]string           `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

type RetrieveItemParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

type UpdateItemParams struct {
	Active        *bool              `json:"active" validate:"omitnil"`
	Description   *string            `json:"description" validate:"omitnil"`
	Group         *string            `json:"group" validate:"omitnil"`
	Inventory     *string            `json:"inventory" validate:"omitnil,uuid"`
	Name          *string            `json:"name" validate:"omitnil"`
	PriceAmount   *int32             `json:"price_amount" validate:"omitnil"`
	PriceCurrency *database.Currency `json:"price_currency" validate:"omitnil,currency"`
	Expand        *[]string          `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

func NewListItemsParams() ListItemsParams {
	limit := int32(10)
	return ListItemsParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
