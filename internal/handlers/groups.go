package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/pkg/groups"
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
	params := &groups.CreateGroupParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	group, err := h.Groups.Create(accountId, params)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "parent_group") {
		_, err := api.ExpandField(&group.ParentGroup, group.ParentGroup.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(group.ParentGroup.ID.UUID, "group"))
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

	err := h.Groups.Delete(accountId, groupId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *GroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := groups.NewListGroupParams()

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	groups, hasMore, err := h.Groups.List(accountId, params)
	if err != nil {
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

	params := &groups.RetrieveGroupParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	group, err := h.Groups.Get(groupId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "parent_group") {
		_, err := api.ExpandField(&group.ParentGroup, group.ParentGroup.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(group.ParentGroup.ID.UUID, "group"))
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

	params := &groups.UpdateGroupParams{}
	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	group, err := h.Groups.Update(groupId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "group"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "parent_group") {
		_, err := api.ExpandField(&group.ParentGroup, group.ParentGroup.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(group.ParentGroup.ID.UUID, "group"))
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}
