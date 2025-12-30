package handlers

import (
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	itemidentifiers "github.com/d-darac/inventory-api/internal/pkg/item_identifiers"
	"github.com/d-darac/inventory-api/internal/pkg/items"
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
	params := &itemidentifiers.CreateItemIdentifiersParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Create(accountId, params)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&itemIdentifiers.Item, itemIdentifiers.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiers.Item.ID.UUID, "item"))
		}
	}

	api.ResJSON(w, http.StatusCreated, itemIdentifiers)
}

func (h *ItemIdentifiersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifiersId, errMsg := api.GetIdFromPath(r)

	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	err := h.ItemIdentifiers.Delete(accountId, itemIdentifiersId)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiersId, "itemIdentifiers"))
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemIdentifiersHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := itemidentifiers.NewListItemIdentifiersParams()

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	itemIdentifiers, hasMore, err := h.ItemIdentifiers.List(accountId, params)
	if err != nil {
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
	itemIdentifiersId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	params := &itemidentifiers.RetrieveItemIdentifiersParams{}

	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Get(itemIdentifiersId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiersId, "itemIdentifiers"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&itemIdentifiers.Item, itemIdentifiers.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiers.Item.ID.UUID, "item"))
		}
	}

	api.ResJSON(w, http.StatusOK, itemIdentifiers)
}

func (h *ItemIdentifiersHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifiersId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	params := &itemidentifiers.UpdateItemIdentifiersParams{}
	if errRes := api.JsonDecode(r, params, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := h.validator.ValidateRequestParams(params); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Update(itemIdentifiersId, accountId, params)
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiersId, "itemIdentifiers"))
		return
	}

	if params.Expand != nil && slices.Contains(*params.Expand, "item") {
		_, err := api.ExpandField(&itemIdentifiers.Item, itemIdentifiers.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(itemIdentifiers.Item.ID.UUID, "item"))
		}
	}

	api.ResJSON(w, http.StatusOK, itemIdentifiers)
}
