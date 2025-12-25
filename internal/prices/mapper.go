package prices

import (
	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func mapCreateParams(accountId uuid.UUID, cp *CreateParams) database.CreatePriceParams {
	cpp := database.CreatePriceParams{
		AccountID: accountId,
		Amount:    cp.Amount,
		Currency:  cp.Currency,
		ItemID:    cp.Item,
	}
	return cpp
}

func mapListParams(accountId uuid.UUID, lp *ListParams) database.ListPricesParams {
	lpp := database.ListPricesParams{
		AccountID: accountId,
	}
	database.MapTimeRange(lp.CreatedAt, &lpp.CreatedAtGt, &lpp.CreatedAtGte, &lpp.CreatedAtLt, &lpp.CreatedAtLte)
	database.MapTimeRange(lp.UpdatedAt, &lpp.UpdatedAtGt, &lpp.UpdatedAtGte, &lpp.UpdatedAtLt, &lpp.UpdatedAtLte)
	database.MapPaginationParams(*lp.PaginationParams, &lpp)
	return lpp
}

func mapUpdateParams(id uuid.UUID, accountId uuid.UUID, up *UpdateParams) database.UpdatePriceParams {
	upp := database.UpdatePriceParams{
		AccountID: accountId,
		Amount:    api.NullInt32(up.Amount),
		Currency:  api.NullCurrency(up.Currency),
		ID:        id,
	}
	return upp
}
