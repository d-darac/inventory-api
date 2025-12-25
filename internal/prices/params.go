package prices

import (
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type CreateParams struct {
	Amount   int32             `json:"amount" validate:"required"`
	Currency database.Currency `json:"currency" validate:"required,currency"`
	Item     uuid.UUID         `json:"item" validate:"required,uuid"`
	Expand   *[]string         `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type ListParams struct {
	*database.PaginationParams
	CreatedAt *database.TimeRange `json:"created_at" validate:"omitnil"`
	UpdatedAt *database.TimeRange `json:"updated_at" validate:"omitnil"`
}

type RetrieveParams struct {
	Expand *[]string `json:"expand" validate:"omitnil,dive,oneof=item"`
}

type UpdateParams struct {
	Amount   *int32             `json:"amount" validate:"omitnil"`
	Currency *database.Currency `json:"currency" validate:"omitnil,currency"`
	Expand   *[]string          `json:"expand" validate:"omitnil,dive,oneof=item"`
}

func NewListParams() *ListParams {
	limit := int32(int(10))
	return &ListParams{
		PaginationParams: &database.PaginationParams{
			Limit: &limit,
		},
	}
}
