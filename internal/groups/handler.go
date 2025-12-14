package groups

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/d-darac/inventory-api/internal/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type Handler struct {
	Db *database.Queries
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	accountId := uuid.MustParse("bfa3562d-7c2a-4880-8634-3decaf00872e")
	type parameters struct {
		Description *string    `json:"description,omitempty"`
		Name        *string    `json:"name,omitempty"`
		ParentGroup *uuid.UUID `json:"parent_group,omitempty"`
	}

	params := &parameters{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	createGroupParams := database.CreateGroupParams{
		AccountID: accountId,
	}
	if params.Description != nil {
		createGroupParams.Description = sql.NullString{
			String: *params.Description,
			Valid:  true,
		}
	}
	if params.ParentGroup != nil {
		createGroupParams.ParentID = uuid.NullUUID{
			UUID:  *params.ParentGroup,
			Valid: true,
		}
	}

	group, err := h.Db.CreateGroup(context.Background(), createGroupParams)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, api.GroupResponse{
		ID:          group.ID,
		CreatedAt:   group.UpdatedAt,
		UpdatedAt:   group.UpdatedAt,
		Description: str.NullString(group.Description),
		Name:        group.Name,
		ParentGroup: group.ParentGroup,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

}
