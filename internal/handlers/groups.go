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

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	group, err := h.Groups.Create(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, group, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, group)
}

func (h *GroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, err := api.GetIdFromPath(r)

	if err != nil {
		api.ResError(w, err)
		return
	}

	if err = h.Groups.Delete(accountId, groupId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *GroupsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := groups.NewListGroupsParams()

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	groups, hasMore, err := h.Groups.List(accountId, params)
	if err != nil {
		api.ResError(w, err)
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
	groupId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &groups.RetrieveGroupParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	group, err := h.Groups.Get(groupId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, group, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &groups.UpdateGroupParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	group, err := h.Groups.Update(groupId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, group, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *GroupsHandler) expandFields(fields *[]string, group *groups.Group, accountId uuid.UUID, w http.ResponseWriter) error {
	if fields != nil && slices.Contains(*fields, "parent_group") {
		if _, err := api.ExpandField(&group.ParentGroup, group.ParentGroup.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get); err != nil {
			return err
		}
	}
	return nil
}
