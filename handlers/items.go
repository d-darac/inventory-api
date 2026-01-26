package handlers

import (
	"net/http"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/middleware"
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
	params := items.CreateItemParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	if params.Group != nil {
		_, err := h.Groups.Get(groups.Get{AccountId: accountId, GroupId: uuid.MustParse(*params.Group)})
		if err != nil {
			api.ResError(w, err)
			return
		}
	}

	if params.Inventory != nil {
		_, err := h.Inventories.Get(inventories.Get{AccountId: accountId, InventoryId: uuid.MustParse(*params.Inventory)})
		if err != nil {
			api.ResError(w, err)
			return
		}
	}

	if params.GroupData != nil {
		group, err := h.Groups.Create(groups.Create{
			AccountId: accountId,
			RequestParams: groups.CreateGroupParams{
				Description: params.GroupData.Description,
				Name:        params.GroupData.Name,
				ParentGroup: params.GroupData.ParentGroup,
			},
		})
		if err != nil {
			api.ResError(w, err)
			return
		}
		groupId := group.ID.String()
		params.Group = &groupId
	}

	if params.InventoryData != nil {
		inventory, err := h.Inventories.Create(inventories.Create{
			AccountId: accountId,
			RequestParams: inventories.CreateInventoryParams{
				InStock:   params.InventoryData.InStock,
				Orderable: params.InventoryData.Orderable,
			},
		})
		if err != nil {
			api.ResError(w, err)
			return
		}
		inventoryId := inventory.ID.String()
		params.Inventory = &inventoryId
	}

	item, err := h.Items.Create(items.Create{AccountId: accountId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if params.IdentifiersData != nil {
		itemIdentifier, err := h.ItemIdentifiers.Create(itemidentifiers.Create{
			AccountId: accountId,
			RequestParams: itemidentifiers.CreateItemIdentifiersParams{
				Ean:  params.IdentifiersData.Ean,
				Gtin: params.IdentifiersData.Gtin,
				Isbn: params.IdentifiersData.Isbn,
				Jan:  params.IdentifiersData.Jan,
				Mpn:  params.IdentifiersData.Mpn,
				Nsn:  params.IdentifiersData.Nsn,
				Upc:  params.IdentifiersData.Upc,
				Qr:   params.IdentifiersData.Qr,
				Sku:  params.IdentifiersData.Sku,
				Item: item.ID.String(),
			},
		})
		if err != nil {
			api.ResError(w, err)
			return
		}
		item.Identifiers = api.Expandable{
			ID: api.NullUUID(itemIdentifier.ID),
		}
	}

	if err := h.ExpandFields(params.Expand, item, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, item)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, err := api.GetIdFromPath(r)

	if err != nil {
		api.ResError(w, err)
		return
	}

	if err = h.Items.Delete(items.Delete{AccountId: accountId, ItemId: itemId}); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := items.NewListItemsParams()

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	items, hasMore, err := h.Items.List(items.List{AccountId: accountId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(items) != 0 {
		listRes.Data = append(listRes.Data, items)
		listRes.HasMore = hasMore
	}

	if err := h.ExpandFieldsList(params.Expand, items, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *ItemsHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := items.RetrieveItemParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	item, err := h.Items.Get(items.Get{
		AccountId:     accountId,
		ItemId:        itemId,
		RequestParams: params,
		OmitBase:      false,
	})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.ExpandFields(params.Expand, item, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, item)
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := items.UpdateItemParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	if params.Group != nil {
		_, err := h.Groups.Get(groups.Get{AccountId: accountId, GroupId: uuid.MustParse(*params.Group)})
		if err != nil {
			api.ResError(w, err)
			return
		}
	}

	if params.Inventory != nil {
		_, err := h.Inventories.Get(inventories.Get{AccountId: accountId, InventoryId: uuid.MustParse(*params.Inventory)})
		if err != nil {
			api.ResError(w, err)
			return
		}
	}

	item, err := h.Items.Update(items.Update{
		AccountId:     accountId,
		ItemId:        itemId,
		RequestParams: params,
	})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.ExpandFields(params.Expand, item, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, item)
}
