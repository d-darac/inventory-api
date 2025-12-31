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

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Create(accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if ok := h.expandFields(params.Expand, itemIdentifiers, accountId, w); !ok {
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

	err = h.ItemIdentifiers.Delete(accountId, itemIdentifiersId)
	if err != nil {
		api.ResError(w, err)
		return
	}

	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *ItemIdentifiersHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	params := itemidentifiers.NewListItemIdentifiersParams()

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, hasMore, err := h.ItemIdentifiers.List(accountId, params)
	if err != nil {
		api.ResError(w, err)
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
	itemIdentifiersId, err := api.GetIdFromPath(r)
	if err != nil {
		api.ResError(w, err)
		return
	}

	params := &itemidentifiers.RetrieveItemIdentifiersParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Get(itemIdentifiersId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if ok := h.expandFields(params.Expand, itemIdentifiers, accountId, w); !ok {
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

	params := &itemidentifiers.UpdateItemIdentifiersParams{}

	if err := api.JsonDecode(r, params, w); err != nil {
		api.ResError(w, err)
		return
	}

	if errs := h.validator.ValidateRequestParams(params); errs != nil {
		api.ResErrorList(w, errs)
		return
	}

	itemIdentifiers, err := h.ItemIdentifiers.Update(itemIdentifiersId, accountId, params)
	if err != nil {
		api.ResError(w, err)
		return
	}

	if ok := h.expandFields(params.Expand, itemIdentifiers, accountId, w); !ok {
		return
	}

	api.ResJSON(w, http.StatusOK, itemIdentifiers)
}

func (h *ItemIdentifiersHandler) expandFields(fields *[]string, itemIdentifiers *itemidentifiers.ItemIdentifiers, accountId uuid.UUID, w http.ResponseWriter) bool {
	if fields != nil && slices.Contains(*fields, "item") {
		_, err := api.ExpandField(&itemIdentifiers.Item, itemIdentifiers.Item.ID.UUID, accountId, &items.RetrieveItemParams{}, h.Items.Get)
		if err != nil {
			api.ResError(w, err)
			return false
		}
	}
	return true
}
