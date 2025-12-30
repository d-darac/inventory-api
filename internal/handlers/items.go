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

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	item, err := h.Items.Create(accountId, params)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "group") {
		_, err := api.ExpandField(&item.Group, item.Group.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Group.ID.UUID, "group"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "inventory") {
		_, err := api.ExpandField(&item.Inventory, item.Inventory.ID.UUID, accountId, &inventories.RetrieveInventoryParams{}, h.Inventories.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Inventory.ID.UUID, "inventory"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item identifiers") {
		_, err := api.ExpandField(&item.Identifiers, item.Identifiers.ID.UUID, accountId, &itemidentifiers.RetrieveItemIdentifiersParams{}, h.ItemIdentifiers.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Identifiers.ID.UUID, "item identifiers"))
		}
	}

	api.ResJSON(w, http.StatusCreated, item)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)

	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	err := h.Items.Delete(accountId, groupId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "item"))
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := items.NewListItemsParams()

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	items, hasMore, err := h.Items.List(accountId, params)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
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
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	params := &items.RetrieveItemParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	item, err := h.Items.Get(groupId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "item"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "group") {
		_, err := api.ExpandField(&item.Group, item.Group.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Group.ID.UUID, "group"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "inventory") {
		_, err := api.ExpandField(&item.Inventory, item.Inventory.ID.UUID, accountId, &inventories.RetrieveInventoryParams{}, h.Inventories.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Inventory.ID.UUID, "inventory"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item identifiers") {
		_, err := api.ExpandField(&item.Identifiers, item.Identifiers.ID.UUID, accountId, &itemidentifiers.RetrieveItemIdentifiersParams{}, h.ItemIdentifiers.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Identifiers.ID.UUID, "item identifiers"))
		}
	}

	api.ResJSON(w, http.StatusOK, item)
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	groupId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	params := &items.UpdateItemParams{}
	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	item, err := h.Items.Update(groupId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(groupId, "item"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "group") {
		_, err := api.ExpandField(&item.Group, item.Group.ID.UUID, accountId, &groups.RetrieveGroupParams{}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Group.ID.UUID, "group"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "inventory") {
		_, err := api.ExpandField(&item.Inventory, item.Inventory.ID.UUID, accountId, &inventories.RetrieveInventoryParams{}, h.Inventories.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Inventory.ID.UUID, "inventory"))
		}
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item identifiers") {
		_, err := api.ExpandField(&item.Identifiers, item.Identifiers.ID.UUID, accountId, &itemidentifiers.RetrieveItemIdentifiersParams{}, h.ItemIdentifiers.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(item.Identifiers.ID.UUID, "item identifiers"))
		}
	}

	api.ResJSON(w, http.StatusOK, item)
}
