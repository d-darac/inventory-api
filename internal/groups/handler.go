package groups

import (
	"context"
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	service "github.com/d-darac/inventory-api/internal/services"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

var validator = api.NewValidator()

type GroupsHandler struct {
	Db *database.Queries
}

func NewHandler(db *database.Queries) *GroupsHandler {
	return &GroupsHandler{
		Db: db,
	}
}

func (h *GroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupsService := service.NewGroupsService(h.Db)
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &createParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	cgp := mapCreateParams(accountId, cp)

	group, err := groupsService.Create(cgp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, groupsService.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}
	api.ResJSON(w, http.StatusCreated, group)
}

func (h *GroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	err = h.Db.DeleteGroup(context.Background(), database.DeleteGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *GroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	lp := newListParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		group, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
			ID:        *lp.StartingAfter,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.StartingAfter, "group"))
			return
		}
		lp.StartingAfterDate = &group.CreatedAt
	}
	if lp.EndingBefore != nil {
		group, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "group"))
			return
		}
		lp.EndingBeforeDate = &group.CreatedAt
	}

	lgp := mapListParams(accountId, lp)

	groups, err := h.Db.ListGroups(context.Background(), lgp)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	listRes.HasMore = len(groups) > int(*lp.Limit)
	if listRes.HasMore {
		if lp.EndingBefore != nil {
			groups = groups[1:]
		} else {
			groups = groups[:len(groups)-1]
		}
	}

	for _, g := range groups {
		listRes.Data = append(listRes.Data, service.Group{
			ID:          g.ID,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
			Description: str.NullString(g.Description),
			Name:        g.Name,
			ParentGroup: api.Expandable{ID: g.ParentGroup},
		})
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *GroupsHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	groupsService := service.NewGroupsService(h.Db)
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	ep := &expandParam{}
	if errRes := api.JsonDecode(r, ep, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(ep); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	group, err := groupsService.Get(groupId, accountId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	if ep.Expand != nil && slices.Contains(*ep.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, groupsService.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	groupsService := service.NewGroupsService(h.Db)
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := groupsService.Get(groupId, accountId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	up := &updateParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	ugp := mapUpdateParams(groupId, accountId, up)

	group, err := groupsService.Update(ugp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, groupsService.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}
