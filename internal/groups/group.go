package groups

import (
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type Group struct {
	ID          *uuid.UUID     `json:"id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	Description str.NullString `json:"description"`
	Name        string         `json:"name"`
	ParentGroup api.Expandable `json:"parent_group"`
}
