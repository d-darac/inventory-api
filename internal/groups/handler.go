package groups

import (
	"context"
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

var validator = api.NewValidator()

type Handler struct {
	Db *database.Queries
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	group, err := h.createGroup(accountId, cgp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, h.getGroup)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}
	api.ResJSON(w, http.StatusCreated, group)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupdId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        groupdId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupdId, "group"))
		return
	}

	err = h.Db.DeleteGroup(context.Background(), database.DeleteGroupParams{
		ID:        groupdId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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
		listRes.Data = append(listRes.Data, Group{
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

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupdId, errMsg := api.GetIdFromPath(r)
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

	group, err := h.getGroup(groupdId, accountId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupdId, "group"))
		return
	}

	if ep.Expand != nil && slices.Contains(*ep.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, h.getGroup)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupdId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.getGroup(groupdId, accountId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupdId, "group"))
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

	ugp := mapUpdateParams(groupdId, accountId, up)

	group, err := h.updateGroup(groupdId, accountId, ugp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "parent_group") {
		id, err := api.ExpandField(&group.ParentGroup, accountId, h.getGroup)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, group)
}

func (h *Handler) createGroup(accountId uuid.UUID, cgp database.CreateGroupParams) (*Group, error) {
	gr, err := h.Db.CreateGroup(context.Background(), cgp)
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          gr.ID,
		CreatedAt:   gr.UpdatedAt,
		UpdatedAt:   gr.UpdatedAt,
		Description: str.NullString(gr.Description),
		Name:        gr.Name,
		ParentGroup: api.Expandable{ID: gr.ParentGroup},
	}, nil
}

func (h *Handler) getGroup(id, accountId uuid.UUID) (*Group, error) {
	gr, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        id,
		AccountID: accountId,
	})
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          gr.ID,
		CreatedAt:   gr.CreatedAt,
		UpdatedAt:   gr.UpdatedAt,
		Description: str.NullString(gr.Description),
		Name:        gr.Name,
		ParentGroup: api.Expandable{ID: gr.ParentGroup},
	}, nil
}

func (h *Handler) updateGroup(id, accountId uuid.UUID, ugp database.UpdateGroupParams) (*Group, error) {
	ugr, err := h.Db.UpdateGroup(context.Background(), ugp)
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          ugr.ID,
		CreatedAt:   ugr.CreatedAt,
		UpdatedAt:   ugr.UpdatedAt,
		Description: str.NullString(ugr.Description),
		Name:        ugr.Name,
		ParentGroup: api.Expandable{ID: ugr.ParentGroup},
	}, nil
}
