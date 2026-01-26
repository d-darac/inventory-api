package handlers

import (
	"net/http"

	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-api/middleware"
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
	params := itemidentifiers.CreateItemIdentifiersParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Create(itemidentifiers.Create{AccountId: accountId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.ExpandFields(params.Expand, itemIdentifiers, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusCreated, itemIdentifiers)
}

func (h *ItemIdentifiersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifiersId, err := api.GetIdFromPath(r)

	if err != nil {
		api.ResError(w, err)
		return
	}

	if err = h.ItemIdentifiers.Delete(itemidentifiers.Delete{AccountId: accountId, ItemIdentifiersId: itemIdentifiersId}); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemIdentifiersHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := itemidentifiers.NewListItemIdentifiersParams()

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, hasMore, err := h.ItemIdentifiers.List(itemidentifiers.List{AccountId: accountId, RequestParams: params})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if len(itemIdentifiers) != 0 {
		listRes.Data = append(listRes.Data, itemIdentifiers)
		listRes.HasMore = hasMore
	}

	if err := h.ExpandFieldsList(params.Expand, itemIdentifiers, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *ItemIdentifiersHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifiersId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := itemidentifiers.RetrieveItemIdentifiersParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Get(itemidentifiers.Get{
		AccountId:         accountId,
		ItemIdentifiersId: itemIdentifiersId,
		RequestParams:     params,
		OmitBase:          false,
	})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.ExpandFields(params.Expand, itemIdentifiers, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, itemIdentifiers)
}

func (h *ItemIdentifiersHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	itemIdentifiersId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := itemidentifiers.UpdateItemIdentifiersParams{}

	if err := api.JsonDecode(r, &params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Update(itemidentifiers.Update{
		AccountId:         accountId,
		ItemIdentifiersId: itemIdentifiersId,
		RequestParams:     params,
	})
	if err != nil {
		api.ResError(w, err)
		return
	}

	if err := h.ExpandFields(params.Expand, itemIdentifiers, accountId); err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusOK, itemIdentifiers)
}
