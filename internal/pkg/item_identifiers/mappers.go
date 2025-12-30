package itemidentifiers

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateItemIdentifiersParams(accountId uuid.UUID, cp *CreateItemIdentifiersParams) database.CreateItemIdentifierParams {
	ciip := database.CreateItemIdentifierParams{
		AccountID: accountId,
		Ean:       api.NullString(cp.Ean),
		Gtin:      api.NullString(cp.Gtin),
		Isbn:      api.NullString(cp.Isbn),
		ItemID:    cp.Item,
		Jan:       api.NullString(cp.Jan),
		Mpn:       api.NullString(cp.Mpn),
		Nsn:       api.NullString(cp.Nsn),
		Upc:       api.NullString(cp.Upc),
		Qr:        api.NullString(cp.Qr),
		Sku:       api.NullString(cp.Sku),
	}
	return ciip
}

func MapListItemIdentifiersParams(accountId uuid.UUID, lp *ListItemIdentifiersParams) database.ListItemIdentifiersParams {
	liip := database.ListItemIdentifiersParams{
		AccountID: accountId,
	}
	database.MapTimeRange(lp.CreatedAt, &liip.CreatedAtGt, &liip.CreatedAtGte, &liip.CreatedAtLt, &liip.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &liip.UpdatedAtGt, &liip.UpdatedAtGte, &liip.UpdatedAtLt, &liip.UpdatedAtLte)
	// database.MapPaginationParams(*lp.PaginationParams, &liip)
	return liip
}

func MapUpdateItemIdentifiersParams(id, accountId uuid.UUID, up *UpdateItemIdentifiersParams) database.UpdateItemIdentifierParams {
	ciip := database.UpdateItemIdentifierParams{
		AccountID: accountId,
		Ean:       api.NullString(up.Ean),
		Gtin:      api.NullString(up.Gtin),
		Isbn:      api.NullString(up.Isbn),
		ID:        id,
		Jan:       api.NullString(up.Jan),
		Mpn:       api.NullString(up.Mpn),
		Nsn:       api.NullString(up.Nsn),
		Upc:       api.NullString(up.Upc),
		Qr:        api.NullString(up.Qr),
		Sku:       api.NullString(up.Sku),
	}
	return ciip
}
