package mappers

import (
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateItemParams(accountId uuid.UUID, cp *items.CreateItemParams) database.CreateItemParams {
	cip := database.CreateItemParams{
		AccountID:     accountId,
		Description:   api.NullString(cp.Description),
		GroupID:       api.NullUUID(cp.Group),
		InventoryID:   api.NullUUID(cp.Inventory),
		Name:          cp.Name,
		PriceAmount:   api.NullInt32(cp.PriceAmount),
		PriceCurrency: api.NullCurrency(cp.PriceCurrency),
		Type:          cp.Type,
	}
	return cip
}

func MapListItemsParams(accountId uuid.UUID, lp *items.ListItemsParams) database.ListItemsParams {
	lip := database.ListItemsParams{
		AccountID:     accountId,
		Active:        api.NullBool(lp.Active),
		Description:   api.NullString(lp.Description),
		GroupID:       api.NullUUID(lp.Group),
		InventoryID:   api.NullUUID(lp.Inventory),
		Name:          api.NullString(lp.Name),
		PriceAmount:   api.NullInt32(lp.PriceAmount),
		PriceCurrency: api.NullCurrency(lp.PriceCurrency),
		Type:          api.NullItemType(lp.Type),
		Variant:       api.NullBool(lp.Variant),
	}
	database.MapTimeRange(lp.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lip)
	return lip
}

func MapUpdateItemParams(id uuid.UUID, accountId uuid.UUID, up *items.UpdateItemParams) database.UpdateItemParams {
	uip := database.UpdateItemParams{
		ID:            id,
		AccountID:     accountId,
		Active:        api.NullBool(up.Active),
		Description:   api.NullString(up.Description),
		GroupID:       api.NullUUID(up.Group),
		InventoryID:   api.NullUUID(up.Inventory),
		Name:          api.NullString(up.Name),
		PriceAmount:   api.NullInt32(up.PriceAmount),
		PriceCurrency: api.NullCurrency(up.PriceCurrency),
	}
	return uip
}
