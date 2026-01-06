package inventories

import (
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/ints"
	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	InStock   int32            `json:"in_stock"`
	Orderable ints.NullInt32   `json:"orderable"`
	Items     api.ListResponse `json:"items"`
}
