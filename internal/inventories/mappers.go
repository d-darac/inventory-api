package inventories

import (
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateInventoryParams(create Create) database.CreateInventoryParams {
	t := time.Now()
	cip := database.CreateInventoryParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
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
		UpdatedAt: time.Now(),
		AccountID: update.AccountId,
		ID:        update.InventoryId,
		InStock:   api.NullInt32(update.RequestParams.InStock),
		Orderable: api.NullInt32(update.RequestParams.Orderable),
	}
}
