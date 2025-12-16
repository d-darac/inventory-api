package groups

import "github.com/google/uuid"

type createParams struct {
	Description *string    `json:"description" validate:"lte=1024"`
	Name        string     `json:"name" validate:"required,lte=64"`
	ParentGroup *uuid.UUID `json:"parent_group"`
}
