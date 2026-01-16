package itemidentifiers

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateItemIdentifiersParams(create Create) database.CreateItemIdentifierParams {
	ciip := database.CreateItemIdentifierParams{
		AccountID: create.AccountId,
		Ean:       api.NullString(create.RequestParams.Ean),
		Gtin:      api.NullString(create.RequestParams.Gtin),
		Isbn:      api.NullString(create.RequestParams.Isbn),
		ItemID:    uuid.MustParse(create.RequestParams.Item),
		Jan:       api.NullString(create.RequestParams.Jan),
		Mpn:       api.NullString(create.RequestParams.Mpn),
		Nsn:       api.NullString(create.RequestParams.Nsn),
		Upc:       api.NullString(create.RequestParams.Upc),
		Qr:        api.NullString(create.RequestParams.Qr),
		Sku:       api.NullString(create.RequestParams.Sku),
	}
	return ciip
}

func MapListItemIdentifiersParams(list List) database.ListItemIdentifiersParams {
	liip := database.ListItemIdentifiersParams{
		AccountID: list.AccountId,
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &liip.CreatedAtGt, &liip.CreatedAtGte, &liip.CreatedAtLt, &liip.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &liip.UpdatedAtGt, &liip.UpdatedAtGte, &liip.UpdatedAtLt, &liip.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &liip)
	return liip
}

func MapUpdateItemIdentifiersParams(update Update) database.UpdateItemIdentifierParams {
	ciip := database.UpdateItemIdentifierParams{
		AccountID: update.AccountId,
		Ean:       api.NullString(update.RequestParams.Ean),
		Gtin:      api.NullString(update.RequestParams.Gtin),
		Isbn:      api.NullString(update.RequestParams.Isbn),
		ID:        update.ItemIdentifiersId,
		Jan:       api.NullString(update.RequestParams.Jan),
		Mpn:       api.NullString(update.RequestParams.Mpn),
		Nsn:       api.NullString(update.RequestParams.Nsn),
		Upc:       api.NullString(update.RequestParams.Upc),
		Qr:        api.NullString(update.RequestParams.Qr),
		Sku:       api.NullString(update.RequestParams.Sku),
	}
	return ciip
}
