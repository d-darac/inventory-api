package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/middleware"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type InventoriesHandler struct {
	Groups          groups.GroupsService
	Inventories     inventories.InventoriesService
	Items           items.ItemsService
	ItemIdentifiers itemidentifiers.ItemIdentifiersService
	validator       *api.Validator
}

func NewInventoriesHandler(db *database.Queries) *InventoriesHandler {
	return &InventoriesHandler{
		Groups:          *groups.NewGroupsService(db),
		Inventories:     *inventories.NewInventoriesService(db),
		Items:           *items.NewItemsService(db),
		ItemIdentifiers: *itemidentifiers.NewItemIdentifiersService(db),
		validator:       api.NewValidator(),
	}
}

func (h *InventoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	params := inventories.CreateInventoryParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Create(inventories.Create{AccountId: accountId, RequestParams: params})
	if err != nil {
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

	if err = h.Inventories.Delete(inventories.Delete{AccountId: accountId, InventoryId: inventoryId}); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *InventoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := inventories.NewListInventoriesParams()

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventories, hasMore, err := h.Inventories.List(inventories.List{AccountId: accountId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(inventories) != 0 {
		listRes.Data = append(listRes.Data, inventories)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *InventoriesHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	inventoryId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}
	listRes := api.NewListResponse(r)
	params := inventories.NewListItemsParams()

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	items, hasMore, err := h.Items.List(items.List{
		AccountId: accountId,
		RequestParams: items.ListItemsParams{
			PaginationParams: params.PaginationParams,
			Active:           params.Active,
			CreatedAt:        params.CreatedAt,
			Description:      params.Description,
			Group:            params.Group,
			Inventory:        func() *string { s := inventoryId.String(); return &s }(),
			Name:             params.Name,
			PriceAmount:      params.PriceAmount,
			PriceCurrency:    params.PriceCurrency,
			Type:             params.Type,
			UpdatedAt:        params.UpdatedAt,
			Variant:          params.Variant,
		},
	})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(items) != 0 {
		listRes.Data = append(listRes.Data, items)
		listRes.HasMore = hasMore
	}

	if err := h.expandFieldsItemsList(params.Expand, items, accountId); err != nil {
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

	params := inventories.RetrieveInventoryParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Get(inventories.Get{
		AccountId:     accountId,
		InventoryId:   inventoryId,
		RequestParams: params,
		OmitBase:      false,
	})
	if err != nil {
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

	params := inventories.UpdateInventoryParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	inventory, err := h.Inventories.Update(inventories.Update{AccountId: accountId, InventoryId: inventoryId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, inventory)
}

func (h *InventoriesHandler) expandFieldsItems(fields *[]string, item *items.Item, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "group") {
		getParams := groups.Get{
			AccountId:     accountId,
			GroupId:       item.Group.ID.UUID,
			RequestParams: groups.RetrieveGroupParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&item.Group, h.Groups.Get, getParams); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "inventory") {
		getParams := inventories.Get{
			AccountId:     accountId,
			InventoryId:   item.Inventory.ID.UUID,
			RequestParams: inventories.RetrieveInventoryParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&item.Inventory, h.Inventories.Get, getParams); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "identifiers") {
		getParams := itemidentifiers.Get{
			AccountId:         accountId,
			ItemIdentifiersId: item.Identifiers.ID.UUID,
			RequestParams:     itemidentifiers.RetrieveItemIdentifiersParams{},
			OmitBase:          true,
		}
		if _, err := api.ExpandField(&item.Identifiers, h.ItemIdentifiers.Get, getParams); err != nil {
			return err
		}
	}
	return nil
}

func (h *InventoriesHandler) expandFieldsItemsList(fields *[]string, items []*items.Item, accountId uuid.UUID) error {
	for _, i := range items {
		if err := h.expandFieldsItems(fields, i, accountId); err != nil {
			return err
		}
	}
	return nil
}
