package inventories

import (
	"time"

	"github.com/d-darac/inventory-assets/ints"
	"github.com/google/uuid"
)

type Inventory struct {
	ID        *uuid.UUID     `json:"id,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	InStock   int32          `json:"in_stock"`
	Orderable ints.NullInt32 `json:"orderable"`
	Reserved  ints.NullInt32 `json:"reserved"`
}
