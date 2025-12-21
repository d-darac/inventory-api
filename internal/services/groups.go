package services

import (
	"context"

	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Description str.NullString `json:"description"`
	Name        string         `json:"name"`
	ParentGroup api.Expandable `json:"parent_group"`
}

type GroupsService struct {
	Db *database.Queries
}

func NewGroupsService(db *database.Queries) *GroupsService {
	return &GroupsService{
		Db: db,
	}
}

func (s *GroupsService) Create(cgp database.CreateGroupParams) (*Group, error) {
	gr, err := s.Db.CreateGroup(context.Background(), cgp)
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          gr.ID,
		CreatedAt:   gr.UpdatedAt,
		UpdatedAt:   gr.UpdatedAt,
		Description: str.NullString(gr.Description),
		Name:        gr.Name,
		ParentGroup: api.Expandable{ID: gr.ParentGroup},
	}, nil
}

func (s *GroupsService) Get(id, accountId uuid.UUID) (*Group, error) {
	gr, err := s.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        id,
		AccountID: accountId,
	})
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          gr.ID,
		CreatedAt:   gr.CreatedAt,
		UpdatedAt:   gr.UpdatedAt,
		Description: str.NullString(gr.Description),
		Name:        gr.Name,
		ParentGroup: api.Expandable{ID: gr.ParentGroup},
	}, nil
}

func (s *GroupsService) Update(ugp database.UpdateGroupParams) (*Group, error) {
	ugr, err := s.Db.UpdateGroup(context.Background(), ugp)
	if err != nil {
		return nil, err
	}
	return &Group{
		ID:          ugr.ID,
		CreatedAt:   ugr.CreatedAt,
		UpdatedAt:   ugr.UpdatedAt,
		Description: str.NullString(ugr.Description),
		Name:        ugr.Name,
		ParentGroup: api.Expandable{ID: ugr.ParentGroup},
	}, nil
}
