package mappers

import (
	"github.com/d-darac/inventory-api/internal/pkg/inventories"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateInventoryParams(accountId uuid.UUID, cp *inventories.CreateInventoryParams) database.CreateInventoryParams {
	cip := database.CreateInventoryParams{
		AccountID: accountId,
		InStock:   cp.InStock,
		Orderable: api.NullInt32(cp.Orderable),
	}
	return cip
}

func MapListInventoryParams(accountId uuid.UUID, lp *inventories.ListInventoriesParams) database.ListInventoriesParams {
	lip := database.ListInventoriesParams{
		AccountID: accountId,
	}
	database.MapTimeRange(lp.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lip)
	return lip
}

func MapUpdateInventoryParams(id, accountId uuid.UUID, up *inventories.UpdateInventoryParams) database.UpdateInventoryParams {
	return database.UpdateInventoryParams{
		AccountID: accountId,
		ID:        id,
		InStock:   api.NullInt32(up.InStock),
		Orderable: api.NullInt32(up.Orderable),
	}
}
