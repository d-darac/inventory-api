package itemidentifiers

import (
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemIdentifiers struct {
	ID        *uuid.UUID     `json:"id,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	Ean       str.NullString `json:"ean"`
	Gtin      str.NullString `json:"gtin"`
	Isbn      str.NullString `json:"isbn"`
	Jan       str.NullString `json:"jan"`
	Mpn       str.NullString `json:"mpn"`
	Nsn       str.NullString `json:"nsn"`
	Upc       str.NullString `json:"upc"`
	Qr        str.NullString `json:"qr"`
	Sku       str.NullString `json:"sku"`
	Item      api.Expandable `json:"item"`
}
