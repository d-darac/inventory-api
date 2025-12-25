package inventories

import (
	"database/sql"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	InStock   int32          `json:"in_stock"`
	Orderable sql.NullInt32  `json:"orderable"`
	Item      api.Expandable `json:"item"`
}
