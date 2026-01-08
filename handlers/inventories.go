package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/inventories"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/middleware"
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

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Create(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, inventory, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, inventory)
}

func (h *InventoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, err := api.GetIdFromPath(r)

	if err != nil {
		api.ResError(w, err)
		return
	}

	if err = h.Inventories.Delete(accountId, inventoryId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *InventoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := inventories.NewListInventoriesParams()

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventories, hasMore, err := h.Inventories.List(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(inventories) != 0 {
		listRes.Data = append(listRes.Data, inventories)
		listRes.HasMore = hasMore
	}

	if err := h.expandFieldsList(params.Expand, inventories, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *InventoriesHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &inventories.RetrieveInventoryParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Get(inventoryId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, inventory, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, inventory)
}

func (h *InventoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &inventories.UpdateInventoryParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Update(inventoryId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, inventory, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, inventory)
}

func (h *InventoriesHandler) expandFields(fields *[]string, inventory *inventories.Inventory, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "items") {
		iParams := items.NewListItemsParams()
		iParams.Inventory = &inventory.ID
		items, hasMore, err := h.Items.List(accountId, iParams)
		if err != nil {
			return err
		}
		inventory.Items.Data = make([]interface{}, 0)
		for _, i := range items {
			inventory.Items.Data = append(inventory.Items.Data, i)
		}
		inventory.Items.HasMore = hasMore
		inventory.Items.Url = "/v1/items"
	}
	return nil
}

func (h *InventoriesHandler) expandFieldsList(fields *[]string, inventories []*inventories.Inventory, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "items") {
		for _, i := range inventories {
			if err := h.expandFields(&[]string{"items"}, i, accountId); err != nil {
				return err
			}
		}
	}
	return nil
}
