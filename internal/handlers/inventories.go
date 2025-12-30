package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/pkg/inventories"
	"github.com/d-darac/inventory-api/internal/pkg/items"
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
	params := &inventories.CreateInventoryParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	inventory, err := h.Inventories.Create(accountId, params)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&inventory.Item, inventory.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventory.Item.ID.UUID, "item"))
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

	err := h.Inventories.Delete(accountId, inventoryId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *InventoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := inventories.NewListInventoriesParams()

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	inventories, hasMore, err := h.Inventories.List(accountId, params)
	if err != nil {
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

	params := &inventories.RetrieveInventoryParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	inventory, err := h.Inventories.Get(inventoryId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&inventory.Item, inventory.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventory.Item.ID.UUID, "item"))
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

	params := &inventories.UpdateInventoryParams{}
	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	inventory, err := h.Inventories.Update(inventoryId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventoryId, "inventory"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&inventory.Item, inventory.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(inventory.Item.ID.UUID, "item"))
		}
	}

	api.ResJSON(w, http.StatusOK, inventory)
}
