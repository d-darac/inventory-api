package groups

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

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

	if errListRes := validateCreateParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	cgp := mapCreateParams(cp)
	cgp.AccountID = accountId

	group, err := h.Db.CreateGroup(context.Background(), cgp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, groupResponse{
		ID:          group.ID,
		CreatedAt:   group.UpdatedAt,
		UpdatedAt:   group.UpdatedAt,
		Description: str.NullString(group.Description),
		Name:        group.Name,
		ParentGroup: group.ParentGroup,
	})
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
	listRes := api.ListResponse{
		Data: []interface{}{},
		Url:  r.URL.Path,
	}
	lp := &listParams{}

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validateListParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	groups, err := h.Db.ListGroups(context.Background(), database.ListGroupsParams{
		AccountID: accountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	for _, g := range groups {
		listRes.Data = append(listRes.Data, groupResponse{
			ID:          g.ID,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
			Description: str.NullString(g.Description),
			Name:        g.Name,
			ParentGroup: g.ParentGroup,
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

	group, err := h.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        groupdId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupdId, "group"))
		return
	}

	api.ResJSON(w, http.StatusOK, groupResponse{
		ID:          group.ID,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
		Description: str.NullString(group.Description),
		Name:        group.Name,
		ParentGroup: group.ParentGroup,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

}
