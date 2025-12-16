package groups

import (
	"database/sql"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(cp *createParams) database.CreateGroupParams {
	cgp := database.CreateGroupParams{
		Name: cp.Name,
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
