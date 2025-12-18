package groups

import (
	"database/sql"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(accountId uuid.UUID, cp *createParams) database.CreateGroupParams {
	cgp := database.CreateGroupParams{
		Name:      cp.Name,
		AccountID: accountId,
	}
	if cp.Description != nil {
		cgp.Description = sql.NullString{
			String: *cp.Description,
			Valid:  true,
		}
	}
	if cp.ParentGroup != nil {
		cgp.ParentID = uuid.NullUUID{
			UUID:  *cp.ParentGroup,
			Valid: true,
		}
	}
	return cgp
}

func mapListParams(accountId uuid.UUID, lp *listParams) database.ListGroupsParams {
	lgp := database.ListGroupsParams{
		AccountID: accountId,
	}

	database.MapTimeRange(lp.CreatedAt, &lgp.CreatedAtGt, &lgp.CreatedAtGte, &lgp.CreatedAtLt, &lgp.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lgp.UpdatedAtGt, &lgp.UpdatedAtGte, &lgp.UpdatedAtLt, &lgp.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lgp)

	if lp.Description != nil {
		lgp.Description = sql.NullString{
			String: *lp.Description,
			Valid:  true,
		}
	}
	if lp.Name != nil {
		lgp.Name = sql.NullString{
			String: *lp.Name,
			Valid:  true,
		}
	}
	if lp.ParentGroup != nil {
		lgp.ParentID = uuid.NullUUID{
			UUID:  *lp.ParentGroup,
			Valid: true,
		}
	}
	return lgp
}

func mapUpdateParams(id uuid.UUID, accountId uuid.UUID, up *updateParams) database.UpdateGroupParams {
	ugp := database.UpdateGroupParams{
		ID:        id,
		AccountID: accountId,
	}

	if up.Description != nil {
		ugp.Description = sql.NullString{
			String: *up.Description,
			Valid:  true,
		}
	}
	if up.Name != nil {
		ugp.Name = sql.NullString{
			String: *up.Name,
			Valid:  true,
		}
	}
	if up.ParentGroup != nil {
		ugp.ParentID = uuid.NullUUID{
			UUID:  *up.ParentGroup,
			Valid: true,
		}
	}
	return ugp
}
