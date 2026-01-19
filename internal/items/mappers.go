package items

import (
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateItemParams(create Create) database.CreateItemParams {
	t := time.Now()
	cip := database.CreateItemParams{
		ID:            uuid.New(),
		CreatedAt:     t,
		UpdatedAt:     t,
		AccountID:     create.AccountId,
		Description:   api.NullString(create.RequestParams.Description),
		GroupID:       api.NullUUID(nil),
		InventoryID:   api.NullUUID(nil),
		Name:          create.RequestParams.Name,
		PriceAmount:   api.NullInt32(create.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(create.RequestParams.PriceCurrency),
		Type:          create.RequestParams.Type,
	}
	if create.RequestParams.Group != nil {
		groupId := uuid.MustParse(*create.RequestParams.Group)
		cip.GroupID = api.NullUUID(&groupId)
	}
	if create.RequestParams.Inventory != nil {
		inventoryId := uuid.MustParse(*create.RequestParams.Inventory)
		cip.InventoryID = api.NullUUID(&inventoryId)
	}
	return cip
}

func MapListItemsParams(list List) database.ListItemsParams {
	lip := database.ListItemsParams{
		AccountID:     list.AccountId,
		Active:        api.NullBool(list.RequestParams.Active),
		Description:   api.NullString(list.RequestParams.Description),
		GroupID:       api.NullUUID(nil),
		InventoryID:   api.NullUUID(nil),
		Name:          api.NullString(list.RequestParams.Name),
		PriceAmount:   api.NullInt32(list.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(list.RequestParams.PriceCurrency),
		Type:          api.NullItemType(list.RequestParams.Type),
		Variant:       api.NullBool(list.RequestParams.Variant),
	}
	if list.RequestParams.Group != nil {
		groupId := uuid.MustParse(*list.RequestParams.Group)
		lip.GroupID = api.NullUUID(&groupId)
	}
	if list.RequestParams.Inventory != nil {
		inventoryId := uuid.MustParse(*list.RequestParams.Inventory)
		lip.InventoryID = api.NullUUID(&inventoryId)
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &lip)
	return lip
}

func MapUpdateItemParams(update Update) database.UpdateItemParams {
	t := time.Now()
	uip := database.UpdateItemParams{
		UpdatedAt:     t,
		ID:            update.ItemId,
		AccountID:     update.AccountId,
		Active:        api.NullBool(update.RequestParams.Active),
		Description:   api.NullString(update.RequestParams.Description),
		GroupID:       api.NullUUID(nil),
		InventoryID:   api.NullUUID(nil),
		Name:          api.NullString(update.RequestParams.Name),
		PriceAmount:   api.NullInt32(update.RequestParams.PriceAmount),
		PriceCurrency: api.NullCurrency(update.RequestParams.PriceCurrency),
	}
	if update.RequestParams.Group != nil {
		groupId := uuid.MustParse(*update.RequestParams.Group)
		uip.GroupID = api.NullUUID(&groupId)
	}
	if update.RequestParams.Inventory != nil {
		inventoryId := uuid.MustParse(*update.RequestParams.Inventory)
		uip.InventoryID = api.NullUUID(&inventoryId)
	}
	return uip
}
