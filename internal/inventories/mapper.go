package inventories

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(accountId uuid.UUID, cp *CreateParams) database.CreateInventoryParams {
	cip := database.CreateInventoryParams{
		AccountID: accountId,
		InStock:   cp.InStock,
		Orderable: api.NullInt32(cp.Orderable),
	}
	return cip
}

func mapListParams(accountId uuid.UUID, lp *ListParams) database.ListInventoriesParams {
	lip := database.ListInventoriesParams{
		AccountID: accountId,
	}
	database.MapTimeRange(lp.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lip)
	return lip
}

func mapUpdateParams(id, accountId uuid.UUID, up *UpdateParams) database.UpdateInventoryParams {
	return database.UpdateInventoryParams{
		AccountID: accountId,
		ID:        id,
		InStock:   api.NullInt32(up.InStock),
		Orderable: api.NullInt32(up.Orderable),
	}
}
