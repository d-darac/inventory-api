package handlers

import (
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/mappers"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/pkg/groups"
	"github.com/d-darac/inventory-api/internal/pkg/inventories"
	"github.com/d-darac/inventory-api/internal/pkg/items"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type ItemsHandler struct {
	Groups      *groups.GroupsService
	Inventories *inventories.InventoriesService
	Items       *items.ItemsService
	validator   *api.Validator
}

func NewItemsHandler(db *database.Queries) *ItemsHandler {
	return &ItemsHandler{
		Groups:      groups.NewGroupsService(db),
		Inventories: inventories.NewInventoriesService(db),
		Items:       items.NewItemsService(db),
		validator:   api.NewValidator(),
	}
}

func (h *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &items.CreateItemParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if cp.GroupData != nil {
		cgp := mappers.MapCreateGroupParams(accountId, cp.GroupData)
		group, err := h.Groups.Create(cgp)
		if err != nil {
			api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
			return
		}
		cp.Group = &group.ID
	}

	if cp.InventoryData != nil {
		cip := mappers.MapCreateInventoryParams(accountId, cp.InventoryData)
		inventory, err := h.Inventories.Create(cip)
		if err != nil {
			api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
			return
		}
		cp.Inventory = &inventory.ID
	}

	cip := mappers.MapCreateItemParams(accountId, cp)

	item, err := h.Items.Create(cip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, item)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Items.Get(database.GetItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemId, "item"))
		return
	}

	err = h.Items.Delete(database.DeleteItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	lp := items.NewListItemsParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		item, err := h.Items.Get(database.GetItemParams{
			ID:        *lp.StartingAfter,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.StartingAfter, "item"))
			return
		}
		lp.StartingAfterDate = &item.CreatedAt
	}
	if lp.EndingBefore != nil {
		item, err := h.Items.Get(database.GetItemParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "item"))
			return
		}
		lp.EndingBeforeDate = &item.CreatedAt
	}

	lip := mappers.MapListItemsParams(accountId, lp)

	items, hasMore, err := h.Items.List(lip)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	listRes.Data = append(listRes.Data, items)
	listRes.HasMore = hasMore

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *ItemsHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	rp := &items.RetrieveItemParams{}
	if errRes := api.JsonDecode(r, rp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(rp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	item, err := h.Items.Get(database.GetItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemId, "item"))
		return
	}

	if rp.Expand != nil && slices.Contains(*rp.Expand, "group") {
		id, err := api.ExpandField(&item.Group, database.GetGroupParams{
			ID:        item.Group.ID.UUID,
			AccountID: accountId,
		}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, item)
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Items.Get(database.GetItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemId, "item"))
		return
	}

	up := &items.UpdateItemParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	uip := mappers.MapUpdateItemParams(itemId, accountId, up)

	item, err := h.Items.Update(uip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "group") {
		id, err := api.ExpandField(&item.Group, database.GetGroupParams{
			ID:        item.Group.ID.UUID,
			AccountID: accountId,
		}, h.Groups.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "group"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, item)
}
