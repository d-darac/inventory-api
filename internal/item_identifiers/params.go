package itemidentifiers

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateItemIdentifiersParams struct {
	Ean    *string  `json:"ean" validate:"omitnil"`
	Gtin   *string  `json:"gtin" validate:"omitnil"`
	Isbn   *string  `json:"isbn" validate:"omitnil"`
	Jan    *string  `json:"jan" validate:"omitnil"`
	Mpn    *string  `json:"mpn" validate:"omitnil"`
	Nsn    *string  `json:"nsn" validate:"omitnil"`
	Upc    *string  `json:"upc" validate:"omitnil"`
	Qr     *string  `json:"qr" validate:"omitnil"`
	Sku    *string  `json:"sku" validate:"omitnil"`
	Item   string   `json:"item" validate:"required,uuid"`
	Expand []string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type ListItemIdentifiersByIdsParams struct {
	Ids []uuid.UUID
}

type ListItemIdentifiersParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
	Expand    []string            `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type RetrieveItemIdentifiersParams struct {
	Expand []string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type UpdateItemIdentifiersParams struct {
	Ean    *string  `json:"ean" validate:"omitnil"`
	Gtin   *string  `json:"gtin" validate:"omitnil"`
	Isbn   *string  `json:"isbn" validate:"omitnil"`
	Jan    *string  `json:"jan" validate:"omitnil"`
	Mpn    *string  `json:"mpn" validate:"omitnil"`
	Nsn    *string  `json:"nsn" validate:"omitnil"`
	Upc    *string  `json:"upc" validate:"omitnil"`
	Qr     *string  `json:"qr" validate:"omitnil"`
	Sku    *string  `json:"sku" validate:"omitnil"`
	Expand []string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

func NewListItemIdentifiersParams() ListItemIdentifiersParams {
	limit := int32(10)
	return ListItemIdentifiersParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
