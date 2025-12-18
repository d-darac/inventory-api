package groups

import (
	"time"

	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type group struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Description str.NullString `json:"description"`
	Name        string         `json:"name"`
	ParentGroup uuid.NullUUID  `json:"parent_group"`
}
