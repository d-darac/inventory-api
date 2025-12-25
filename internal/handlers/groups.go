package handlers

import (
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/mappers"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type GroupsHandler struct {
	Groups    groups.GroupsService
	validator *api.Validator
}

func NewGroupsHandler(db *database.Queries) *GroupsHandler {
	return &GroupsHandler{
		Groups:    *groups.NewGroupsService(db),
		validator: api.NewValidator(),
	}
}

func (h *GroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &groups.CreateGroupParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	cgp := mappers.MapCreateGroupParams(accountId, cp)

	group, err := h.Groups.Create(cgp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, database.GetGroupParams{
			ID:        group.ParentGroup.ID.UUID,
			AccountID: accountId,
		}, h.Groups.Get)
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
	_, err := h.Groups.Get(database.GetGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}
	err = h.Groups.Delete(database.DeleteGroupParams{
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
	lp := groups.NewListGroupParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		group, err := h.Groups.Get(database.GetGroupParams{
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
		group, err := h.Groups.Get(database.GetGroupParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "group"))
			return
		}
		lp.EndingBeforeDate = &group.CreatedAt
	}

	lgp := mappers.MapListGroupsParams(accountId, lp)

	groups, hasMore, err := h.Groups.List(lgp)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if len(groups) != 0 {
		listRes.Data = append(listRes.Data, groups)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *GroupsHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	rp := &groups.RetrieveGroupParams{}
	if errRes := api.JsonDecode(r, rp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(rp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	group, err := h.Groups.Get(database.GetGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	if rp.Expand != nil && slices.Contains(*rp.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, database.GetGroupParams{
			ID:        group.ParentGroup.ID.UUID,
			AccountID: accountId,
		}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Groups.Get(database.GetGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	up := &groups.UpdateGroupParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	ugp := mappers.MapUpdateGroupParams(groupId, accountId, up)

	group, err := h.Groups.Update(ugp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, database.GetGroupParams{
			ID:        group.ParentGroup.ID.UUID,
			AccountID: accountId,
		}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}
