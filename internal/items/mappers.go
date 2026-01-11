package items

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
)

func MapCreateItemParams(create Create) database.CreateItemParams {
	cip := database.CreateItemParams{
		AccountID:     create.AccountId,
		Description:   api.NullString(create.RequestParams.Description),
		GroupID:       api.NullUUID(create.RequestParams.Group),
		InventoryID:   api.NullUUID(create.RequestParams.Inventory),
		Name:          create.RequestParams.Name,
		PriceAmount:   api.NullInt32(create.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(create.RequestParams.PriceCurrency),
		Type:          create.RequestParams.Type,
	}
	return cip
}

func MapListItemsParams(list List) database.ListItemsParams {
	lip := database.ListItemsParams{
		AccountID:     list.AccountId,
		Active:        api.NullBool(list.RequestParams.Active),
		Description:   api.NullString(list.RequestParams.Description),
		GroupID:       api.NullUUID(list.RequestParams.Group),
		InventoryID:   api.NullUUID(list.RequestParams.Inventory),
		Name:          api.NullString(list.RequestParams.Name),
		PriceAmount:   api.NullInt32(list.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(list.RequestParams.PriceCurrency),
		Type:          api.NullItemType(list.RequestParams.Type),
		Variant:       api.NullBool(list.RequestParams.Variant),
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &lip)
	return lip
}

func MapUpdateItemParams(update Update) database.UpdateItemParams {
	uip := database.UpdateItemParams{
		ID:            update.ItemId,
		AccountID:     update.AccountId,
		Active:        api.NullBool(update.RequestParams.Active),
		Description:   api.NullString(update.RequestParams.Description),
		GroupID:       api.NullUUID(update.RequestParams.Group),
		InventoryID:   api.NullUUID(update.RequestParams.Inventory),
		Name:          api.NullString(update.RequestParams.Name),
		PriceAmount:   api.NullInt32(update.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(update.RequestParams.PriceCurrency),
	}
	return uip
}
