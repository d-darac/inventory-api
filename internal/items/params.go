package items

import (
	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateItemParams struct {
	Description     *string                                      `json:"description" validate:"omitnil"`
	Group           *uuid.UUID                                   `json:"group" validate:"omitnil,uuid"`
	GroupData       *groups.CreateGroupParams                    `json:"group_data" validate:"omitnil"`
	IdentifiersData *itemidentifiers.CreateItemIdentifiersParams `json:"identifiers_data" validate:"omitnil"`
	Inventory       *uuid.UUID                                   `json:"inventory" validate:"omitnil,uuid"`
	InventoryData   *inventories.CreateInventoryParams           `json:"inventory_data" valdiate:"omitnil"`
	Name            string                                       `json:"name" validate:"required"`
	PriceAmount     *int32                                       `json:"price_amount" validate:"omitnil"`
	PriceCurrency   *database.Currency                           `json:"price_currency" validate:"omitnil,currency"`
	Type            database.ItemType                            `json:"type" validate:"required,itemtype"`
	Expand          *[]string                                    `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

type ListItemsParams struct {
	*database.PaginationParams
	Active        *bool               `json:"active" validate:"omitnil"`
	CreatedAt     *database.TimeRange `json:"created_at" validate:"omitnil"`
	Description   *string             `json:"description" validate:"omitnil"`
	Group         *uuid.UUID          `json:"group" validate:"omitnil"`
	Inventory     *uuid.UUID          `json:"inventory" validate:"omitnil,uuid"`
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
	Group         *uuid.UUID         `json:"group" validate:"omitnil"`
	Inventory     *uuid.UUID         `json:"inventory" validate:"omitnil,uuid"`
	Name          *string            `json:"name" validate:"omitnil"`
	PriceAmount   *int32             `json:"price_amount" validate:"omitnil"`
	PriceCurrency *database.Currency `json:"price_currency" validate:"omitnil,currency"`
	Expand        *[]string          `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory"`
}

func NewListItemsParams() *ListItemsParams {
	limit := int32(10)
	return &ListItemsParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
