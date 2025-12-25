package groups

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(accountId uuid.UUID, cp *CreateParams) database.CreateGroupParams {
	cgp := database.CreateGroupParams{
		AccountID:   accountId,
		Description: api.NullString(cp.Description),
		Name:        cp.Name,
		ParentID:    api.NullUUID(cp.ParentGroup),
	}
	return cgp
}

func mapListParams(accountId uuid.UUID, lp *ListParams) database.ListGroupsParams {
	lgp := database.ListGroupsParams{
		AccountID:   accountId,
		Description: api.NullString(lp.Description),
		Name:        api.NullString(lp.Name),
		ParentID:    api.NullUUID(lp.ParentGroup),
	}
	database.MapTimeRange(lp.CreatedAt, &lgp.CreatedAtGt, &lgp.CreatedAtGte, &lgp.CreatedAtLt, &lgp.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lgp.UpdatedAtGt, &lgp.UpdatedAtGte, &lgp.UpdatedAtLt, &lgp.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lgp)
	return lgp
}

func mapUpdateParams(id uuid.UUID, accountId uuid.UUID, up *UpdateParams) database.UpdateGroupParams {
	return database.UpdateGroupParams{
		AccountID:   accountId,
		Description: api.NullString(up.Description),
		ID:          id,
		Name:        api.NullString(up.Name),
		ParentID:    api.NullUUID(up.ParentGroup),
	}
}
