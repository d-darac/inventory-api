package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/pkg/groups"
	"github.com/d-darac/inventory-api/internal/pkg/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/pkg/item_identifiers"
	"github.com/d-darac/inventory-api/internal/pkg/items"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type ItemsHandler struct {
	Groups          *groups.GroupsService
	Inventories     *inventories.InventoriesService
	ItemIdentifiers *itemidentifiers.ItemIdentifiersService
	Items           *items.ItemsService
	validator       *api.Validator
}

func NewItemsHandler(db *database.Queries) *ItemsHandler {
	return &ItemsHandler{
		Groups:          groups.NewGroupsService(db),
		Inventories:     inventories.NewInventoriesService(db),
		ItemIdentifiers: itemidentifiers.NewItemIdentifiersService(db),
		Items:           items.NewItemsService(db),
		validator:       api.NewValidator(),
	}
}

func (h *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	params := &items.CreateItemParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	item, err := h.Items.Create(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, item, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, item)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, err := api.GetIdFromPath(r)

	if err != nil {
		api.ResError(w, err)
		return
	}

	if err = h.Items.Delete(accountId, groupId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := items.NewListItemsParams()

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	items, hasMore, err := h.Items.List(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(items) != 0 {
		listRes.Data = append(listRes.Data, items)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *ItemsHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &items.RetrieveItemParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	item, err := h.Items.Get(groupId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, item, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, item)
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &items.UpdateItemParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	item, err := h.Items.Update(groupId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.expandFields(params.Expand, item, accountId, w); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, item)
}

func (h *ItemsHandler) expandFields(fields *[]string, item *items.Item, accountId uuid.UUID, w http.ResponseWriter) error {
	if fields != nil && slices.Contains(*fields, "group") {
		if _, err := api.ExpandField(&item.Group, item.Group.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "inventory") {
		if _, err := api.ExpandField(&item.Inventory, item.Inventory.ID.UUID, accountId, &inventories.RetrieveInventoryParams{}, h.Inventories.Get); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "item identifiers") {
		if _, err := api.ExpandField(&item.Identifiers, item.Identifiers.ID.UUID, accountId, &itemidentifiers.RetrieveItemIdentifiersParams{}, h.ItemIdentifiers.Get); err != nil {
			return err
		}

	}
	return nil
}
