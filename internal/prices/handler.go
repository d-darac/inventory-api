package prices

import (
	"database/sql"
	"net/http"
	"slices"

	"github.com/d-darac/inventory-api/internal/middleware"
	"github.com/d-darac/inventory-api/internal/services"
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

var validator = api.NewValidator()

type PricesHandler struct {
	Prices services.PricesService
	Items  services.ItemsService
}

func NewHandler(db *database.Queries) *PricesHandler {
	return &PricesHandler{
		Prices: *services.NewPricesService(db),
		Items:  *services.NewItemsService(db),
	}
}

func (h *PricesHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	cp := &CreateParams{}

	if errRes := api.JsonDecode(r, cp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(cp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	cpp := mapCreateParams(accountId, cp)

	price, err := h.Prices.Create(cpp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if cp.Expand != nil && slices.Contains(*cp.Expand, "item") {
		id, err := api.ExpandField(&price.Item, database.GetItemParams{
			ID:        price.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}
	api.ResJSON(w, http.StatusCreated, price)
}

func (h *PricesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	priceId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}
	_, err := h.Prices.Get(database.GetPriceParams{
		ID:        priceId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(priceId, "price"))
		return
	}
	err = h.Prices.Delete(database.DeletePriceParams{
		ID:        priceId,
		AccountID: accountId,
	})
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}
	api.ResJSON(w, http.StatusNoContent, nil)
}

func (h *PricesHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	listRes := api.NewListResponse(r)
	lp := NewListParams()

	if errRes := api.JsonDecode(r, lp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(lp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if lp.StartingAfter != nil {
		price, err := h.Prices.Get(database.GetPriceParams{
			ID:        *lp.StartingAfter,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.StartingAfter, "price"))
			return
		}
		lp.StartingAfterDate = &price.CreatedAt
	}
	if lp.EndingBefore != nil {
		price, err := h.Prices.Get(database.GetPriceParams{
			ID:        *lp.EndingBefore,
			AccountID: accountId,
		})
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(*lp.EndingBefore, "price"))
			return
		}
		lp.EndingBeforeDate = &price.CreatedAt
	}

	lpp := mapListParams(accountId, lp)

	prices, hasMore, err := h.Prices.List(lpp)
	if err != nil {
		if err == sql.ErrNoRows {
			api.ResJSON(w, http.StatusOK, listRes)
			return
		}
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if len(prices) != 0 {
		listRes.Data = append(listRes.Data, prices)
		listRes.HasMore = hasMore
	}

	api.ResJSON(w, http.StatusOK, listRes)
}

func (h *PricesHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	priceId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	rp := &RetrieveParams{}
	if errRes := api.JsonDecode(r, rp, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(rp); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	price, err := h.Prices.Get(database.GetPriceParams{
		ID:        priceId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(priceId, "price"))
		return
	}

	if rp.Expand != nil && slices.Contains(*rp.Expand, "item") {
		id, err := api.ExpandField(&price.Item, database.GetItemParams{
			ID:        price.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, price)
}

func (h *PricesHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId := r.Context().Value(middleware.AuthAccountID).(uuid.UUID)
	priceId, errMsg := api.GetIdFromPath(r)
	if len(errMsg) > 0 {
		api.ResError(w, http.StatusBadRequest, errMsg, api.InvalidRequestError, nil, nil)
		return
	}

	_, err := h.Prices.Get(database.GetPriceParams{
		ID:        priceId,
		AccountID: accountId,
	})
	if err != nil {
		api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(priceId, "price"))
		return
	}

	up := &UpdateParams{}
	if errRes := api.JsonDecode(r, up, w); errRes != nil {
		errRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	if errListRes := validator.ValidateRequestParams(up); errListRes != nil {
		errListRes.ResError(w, http.StatusBadRequest, nil)
		return
	}

	upp := mapUpdateParams(priceId, accountId, up)

	price, err := h.Prices.Update(upp)
	if err != nil {
		api.ResError(w, http.StatusInternalServerError, api.ApiErrorMessage(), api.ApiError, nil, err)
		return
	}

	if up.Expand != nil && slices.Contains(*up.Expand, "item") {
		id, err := api.ExpandField(&price.Item, database.GetItemParams{
			ID:        price.Item.ID.UUID,
			AccountID: accountId,
		}, h.Items.Get)
		if err != nil {
			api.HandleSqlErrNoRows(err, w, api.NotFoundMessage(id, "item"))
			return
		}
	}

	api.ResJSON(w, http.StatusOK, price)
}
