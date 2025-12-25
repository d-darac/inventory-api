package handlers

import (
	"database/sql"
	"net/http"
	"slices"

	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/internal/mappers"
	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type ItemIdentifiersHandler struct {
	ItemIdentifiers itemidentifiers.ItemIdentifiersService
	Items           items.ItemsService
	validator       *api.Validator
}

func NewItemIdentifiersHandler(db *database.Queries) *ItemIdentifiersHandler {
	return &ItemIdentifiersHandler{
		ItemIdentifiers: *itemidentifiers.NewItemIdentifiersService(db),
		Items:           *items.NewItemsService(db),
		validator:       api.NewValidator(),
	}
}

func (h *ItemIdentifiersHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &itemidentifiers.CreateItemIdentifierParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	ciip := mappers.MapCreateItemIdentifierParams(accountId, cp)

	itemIdentifier, err := h.ItemIdentifiers.Create(ciip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "item") {
		id, err := api.ExpandField(&itemIdentifier.Item, database.GetItemParams{
			ID:        itemIdentifier.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}
	api.ResJSON(w, http.StatusCreated, itemIdentifier)
}

func (h *ItemIdentifiersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifierId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}
	_, err := h.ItemIdentifiers.Get(database.GetItemIdentifierParams{
		ID:        itemIdentifierId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifierId, "item identifier"))
		return
	}
	err = h.ItemIdentifiers.Delete(database.DeleteItemIdentifierParams{
		ID:        itemIdentifierId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}
	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemIdentifiersHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	lp := itemidentifiers.NewListItemIdentifierParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		itemIdentifier, err := h.ItemIdentifiers.Get(database.GetItemIdentifierParams{
			ID:        *lp.StartingAfter,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.StartingAfter, "item identifier"))
			return
		}
		lp.StartingAfterDate = &itemIdentifier.CreatedAt
	}
	if lp.EndingBefore != nil {
		itemIdentifier, err := h.ItemIdentifiers.Get(database.GetItemIdentifierParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "item identifier"))
			return
		}
		lp.EndingBeforeDate = &itemIdentifier.CreatedAt
	}

	liip := mappers.MapListItemIdentifiersParams(accountId, lp)

	itemIdentifiers, hasMore, err := h.ItemIdentifiers.List(liip)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if len(itemIdentifiers) != 0 {
		listRes.Data = append(listRes.Data, itemIdentifiers)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *ItemIdentifiersHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifierId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	rp := &itemidentifiers.RetrieveItemIdentifierParams{}
	if errRes := api.JsonDecode(r, rp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(rp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	itemIdentifier, err := h.ItemIdentifiers.Get(database.GetItemIdentifierParams{
		ID:        itemIdentifierId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifierId, "item identifier"))
		return
	}

	if rp.Expand != nil && slices.Contains(*rp.Expand, "item") {
		id, err := api.ExpandField(&itemIdentifier.Item, database.GetItemParams{
			ID:        itemIdentifier.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, itemIdentifier)
}

func (h *ItemIdentifiersHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifierId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.ItemIdentifiers.Get(database.GetItemIdentifierParams{
		ID:        itemIdentifierId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifierId, "item identifier"))
		return
	}

	up := &itemidentifiers.UpdateItemIdentifierParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	uiip := mappers.MapUpdateItemIdentifierParams(itemIdentifierId, accountId, up)

	itemIdentifier, err := h.ItemIdentifiers.Update(uiip)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "item") {
		id, err := api.ExpandField(&itemIdentifier.Item, database.GetItemParams{
			ID:        itemIdentifier.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, itemIdentifier)
}
