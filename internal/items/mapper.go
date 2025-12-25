package items

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(accountId uuid.UUID, cp *CreateParams) database.CreateItemParams {
	cip := database.CreateItemParams{
		AccountID:   accountId,
		Description: api.NullString(cp.Description),
		GroupID:     api.NullUUID(cp.Group),
		Name:        cp.Name,
		Type:        cp.Type,
	}
	return cip
}

func mapListParams(accountId uuid.UUID, lp *ListParams) database.ListItemsParams {
	lip := database.ListItemsParams{
		AccountID:   accountId,
		Active:      api.NullBool(lp.Active),
		Description: api.NullString(lp.Description),
		GroupID:     api.NullUUID(lp.Group),
		Name:        api.NullString(lp.Name),
		Type:        api.NullItemType(lp.Type),
		Variant:     api.NullBool(lp.Variant),
	}
	database.MapTimeRange(lp.CreatedAt, &lip.CreatedAtGt, &lip.CreatedAtGte, &lip.CreatedAtLt, &lip.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lip.UpdatedAtGt, &lip.UpdatedAtGte, &lip.UpdatedAtLt, &lip.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lip)
	return lip
}

func mapUpdateParams(id uuid.UUID, accountId uuid.UUID, up *UpdateParams) database.UpdateItemParams {
	uip := database.UpdateItemParams{
		ID:          id,
		AccountID:   accountId,
		Active:      api.NullBool(up.Active),
		Description: api.NullString(up.Description),
		GroupID:     api.NullUUID(up.Group),
		Name:        api.NullString(up.Name),
	}
	return uip
}
