package items

import (
	"database/sql"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type Item struct {
	ID            uuid.UUID             `json:"id"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Active        bool                  `json:"active"`
	Description   str.NullString        `json:"description"`
	Group         api.Expandable        `json:"group"`
	Identifiers   api.Expandable        `json:"identifiers"`
	Inventory     api.Expandable        `json:"inventory"`
	Name          string                `json:"name"`
	PriceAmount   sql.NullInt32         `json:"price_amount"`
	PriceCurrency database.NullCurrency `json:"price_currency"`
	Variant       bool                  `json:"variant"`
	Type          database.ItemType     `json:"type"`
}
