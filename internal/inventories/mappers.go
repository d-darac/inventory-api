package inventories

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
)

func MapCreateInventoryParams(create Create) database.CreateInventoryParams {
	cip := database.CreateInventoryParams{
		AccountID: create.AccountId,
		InStock:   create.RequestParams.InStock,
		Orderable: api.NullInt32(create.RequestParams.Orderable),
	}
	return cip
}

func MapListInventoriesParams(list List) database.ListInventoriesParams {
	lip := database.ListInventoriesParams{
		AccountID: list.AccountId,
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &lip)
	return lip
}

func MapUpdateInventoryParams(update Update) database.UpdateInventoryParams {
	return database.UpdateInventoryParams{
		AccountID: update.AccountId,
		ID:        update.InventoryId,
		InStock:   api.NullInt32(update.RequestParams.InStock),
		Orderable: api.NullInt32(update.RequestParams.Orderable),
	}
}
