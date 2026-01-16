package groups

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func MapCreateGroupParams(create Create) database.CreateGroupParams {
	cgp := database.CreateGroupParams{
		AccountID:   create.AccountId,
		Description: api.NullString(create.RequestParams.Description),
		Name:        create.RequestParams.Name,
		ParentID:    api.NullUUID(nil),
	}
	if create.RequestParams.ParentGroup != nil {
		parentGroupId := uuid.MustParse(*create.RequestParams.ParentGroup)
		cgp.ParentID = api.NullUUID(&parentGroupId)
	}
	return cgp
}

func MapListGroupsParams(list List) database.ListGroupsParams {
	lgp := database.ListGroupsParams{
		AccountID:   list.AccountId,
		Description: api.NullString(list.RequestParams.Description),
		Name:        api.NullString(list.RequestParams.Name),
		ParentID:    api.NullUUID(nil),
	}
	if list.RequestParams.ParentGroup != nil {
		parentGroupId := uuid.MustParse(*list.RequestParams.ParentGroup)
		lgp.ParentID = api.NullUUID(&parentGroupId)
	}
	database.MapTimeRange(list.RequestParams.CreatedAt, &lgp.CreatedAtGt, &lgp.CreatedAtGte, &lgp.CreatedAtLt, &lgp.CreatedAtLte)
	database.MapTimeRange(list.RequestParams.UpdatedAt, &lgp.UpdatedAtGt, &lgp.UpdatedAtGte, &lgp.UpdatedAtLt, &lgp.UpdatedAtLte)
	database.MapPaginationParams(*list.RequestParams.PaginationParams, &lgp)
	return lgp
}

func MapUpdateGroupParams(update Update) database.UpdateGroupParams {
	ugp := database.UpdateGroupParams{
		AccountID:   update.AccountId,
		Description: api.NullString(update.RequestParams.Description),
		ID:          update.GroupId,
		Name:        api.NullString(update.RequestParams.Name),
		ParentID:    api.NullUUID(nil),
	}
	if update.RequestParams.ParentGroup != nil {
		parentGroupId := uuid.MustParse(*update.RequestParams.ParentGroup)
		ugp.ParentID = api.NullUUID(&parentGroupId)
	}
	return ugp
}
