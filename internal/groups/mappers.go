package groups

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
)

func MapCreateGroupParams(create Create) database.CreateGroupParams {
	cgp := database.CreateGroupParams{
		AccountID:   create.AccountId,
		Description: api.NullString(create.RequestParams.Description),
		Name:        create.RequestParams.Name,
		ParentID:    api.NullUUID(create.RequestParams.ParentGroup),
	}
	return cgp
}

func MapListGroupsParams(list List) database.ListGroupsParams {
	lgp := database.ListGroupsParams{
		AccountID:   list.AccountId,
		Description: api.NullString(list.RequestParams.Description),
		Name:        api.NullString(list.RequestParams.Name),
		ParentID:    api.NullUUID(list.RequestParams.ParentGroup),
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &lgp.CreatedAtGt, &lgp.CreatedAtGte, &lgp.CreatedAtLt, &lgp.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &lgp.UpdatedAtGt, &lgp.UpdatedAtGte, &lgp.UpdatedAtLt, &lgp.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &lgp)
	return lgp
}

func MapUpdateGroupParams(update Update) database.UpdateGroupParams {
	return database.UpdateGroupParams{
		AccountID:   update.AccountId,
		Description: api.NullString(update.RequestParams.Description),
		ID:          update.GroupId,
		Name:        api.NullString(update.RequestParams.Name),
		ParentID:    api.NullUUID(update.RequestParams.ParentGroup),
	}
}
