package items

import (
	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/prices"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateParams struct {
	Description     *string                       `json:"description" validate:"omitnil"`
	Group           *uuid.UUID                    `json:"group" validate:"omitnil,uuid"`
	GroupData       *groups.CreateParams          `json:"group_data" validate:"omitnil"`
	IdentifiersData *itemidentifiers.CreateParams `json:"identifiers_data" validate:"omitnil"`
	InventoryData   *inventories.CreateParams     `json:"inventory_data" valdiate:"omitnil"`
	Name            string                        `json:"name" validate:"required"`
	PriceData       *prices.CreateParams          `json:"price_data" validate:"omitnil"`
	Type            database.ItemType             `json:"type" validate:"required,itemtype"`
	Expand          *[]string                     `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory price"`
}

type ListParams struct {
	*database.PaginationParams
	Active      *bool               `json:"active" validate:"omitnil"`
	CreatedAt   *database.TimeRange `json:"created_at" validate:"omitnil"`
	Description *string             `json:"description" validate:"omitnil"`
	Group       *uuid.UUID          `json:"group" validate:"omitnil"`
	Name        *string             `json:"name" validate:"omitnil"`
	Type        *database.ItemType  `json:"type" validate:"omitnil,itemtype"`
	UpdatedAt   *database.TimeRange `json:"updated_at" validate:"omitnil"`
	Variant     *bool               `json:"variant" validate:"omitnil"`
}

type RetrieveParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory price"`
}

type UpdateParams struct {
	Active      *bool      `json:"active" validate:"omitnil"`
	Description *string    `json:"description" validate:"omitnil"`
	Group       *uuid.UUID `json:"group" validate:"omitnil"`
	Name        *string    `json:"name" validate:"omitnil"`
	Expand      *[]string  `json:"expand" validate:"omitnil,dive,oneof=group identifiers inventory price"`
}

func NewListParams() *ListParams {
	limit := int32(int(10))
	return &ListParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
