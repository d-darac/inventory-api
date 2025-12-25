package itemidentifiers

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateItemIdentifierParams struct {
	Ean    *string   `json:"ean" validate:"omitnil,len=8|gte=12,lte=14"`
	Gtin   *string   `json:"gtin" validate:"omitnil,len=8|gte=12,lte=14"`
	Isbn   *string   `json:"isbn" validate:"omitnil,len=10|len=13"`
	Jan    *string   `json:"jan" validate:"omitnil,len=8|len=13"`
	Mpn    *string   `json:"mpn" validate:"omitnil"`
	Nsn    *string   `json:"nsn" validate:"omitnil,len=13"`
	Upc    *string   `json:"upc" validate:"omitnil,len=12"`
	Qr     *string   `json:"qr" validate:"omitnil"`
	Sku    *string   `json:"sku" validate:"omitnil"`
	Item   uuid.UUID `json:"item" validate:"required,uuid"`
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type ListItemIdentifiersParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
}

type RetrieveItemIdentifierParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type UpdateItemIdentifierParams struct {
	Ean    *string   `json:"ean" validate:"omitnil,len=8|gte=12,lte=14"`
	Gtin   *string   `json:"gtin" validate:"omitnil,len=8|gte=12,lte=14"`
	Isbn   *string   `json:"isbn" validate:"omitnil,len=10|len=13"`
	Jan    *string   `json:"jan" validate:"omitnil,len=8|len=13"`
	Mpn    *string   `json:"mpn" validate:"omitnil"`
	Nsn    *string   `json:"nsn" validate:"omitnil,len=13"`
	Upc    *string   `json:"upc" validate:"omitnil,len=12"`
	Qr     *string   `json:"qr" validate:"omitnil"`
	Sku    *string   `json:"sku" validate:"omitnil"`
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

func NewListItemIdentifierParams() *ListItemIdentifiersParams {
	limit := int32(int(10))
	return &ListItemIdentifiersParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
