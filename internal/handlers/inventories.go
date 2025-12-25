package handlers

import (
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/inventories"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/internal/mappers"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type InventoriesHandler struct {
	Inventories inventories.InventoriesService
	Items       items.ItemsService
	validator   *api.Validator
}

func NewInventoriesHandler(db *database.Queries) *InventoriesHandler {
	return &InventoriesHandler{
		Inventories: *inventories.NewInventoriesService(db),
		Items:       *items.NewItemsService(db),
		validator:   api.NewValidator(),
	}
}

func (h *InventoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &inventories.CreateInventoryParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	cip := mappers.MapCreateInventoryParams(accountId, cp)

	inventory, err := h.Inventories.Create(cip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "item") {
		id, err := api.ExpandField(&inventory.Item, database.GetItemParams{
			ID:        inventory.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}
	api.ResJSON(w, http.StatusCreated, inventory)
}

func (h *InventoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}
	_, err := h.Inventories.Get(database.GetInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}
	err = h.Inventories.Delete(database.DeleteInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}
	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *InventoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	lp := inventories.NewListInventoriesParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		inventory, err := h.Inventories.Get(database.GetInventoryParams{
			ID:        *lp.StartingAfter,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.StartingAfter, "inventory"))
			return
		}
		lp.StartingAfterDate = &inventory.CreatedAt
	}
	if lp.EndingBefore != nil {
		inventory, err := h.Inventories.Get(database.GetInventoryParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "inventory"))
			return
		}
		lp.EndingBeforeDate = &inventory.CreatedAt
	}

	lip := mappers.MapListInventoryParams(accountId, lp)

	inventories, hasMore, err := h.Inventories.List(lip)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if len(inventories) != 0 {
		listRes.Data = append(listRes.Data, inventories)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *InventoriesHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	rp := &inventories.RetrieveInventoryParams{}
	if errRes := api.JsonDecode(r, rp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(rp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	inventory, err := h.Inventories.Get(database.GetInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}

	if rp.Expand != nil && slices.Contains(*rp.Expand, "item") {
		id, err := api.ExpandField(&inventory.Item, database.GetItemParams{
			ID:        inventory.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, inventory)
}

func (h *InventoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Inventories.Get(database.GetInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}

	up := &inventories.UpdateInventoryParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	uip := mappers.MapUpdateInventoryParams(inventoryId, accountId, up)

	inventory, err := h.Inventories.Update(uip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "item") {
		id, err := api.ExpandField(&inventory.Item, database.GetItemParams{
			ID:        inventory.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, inventory)
}
