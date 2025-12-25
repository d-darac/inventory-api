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

func (s *GroupsService) Delete(dgp database.DeleteGroupParams) error {
	return s.Db.DeleteGroup(context.Background(), dgp)
}

func (s *GroupsService) Get(ggp database.GetGroupParams) (*Group, error) {
	gr, err := s.Db.GetGroup(context.Background(), ggp)
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

func (s *GroupsService) List(lgp database.ListGroupsParams) (groups []*Group, hasMore bool, err error) {
	grs, err := s.Db.ListGroups(context.Background(), lgp)
	if err != nil {
		return
	}
	if lgp.Limit.Valid {
		hasMore = len(grs) > int(lgp.Limit.Int32)
	} else {
		hasMore = len(grs) > 10
	}
	if hasMore {
		if lgp.EndingBefore.Valid {
			grs = grs[1:]
		} else {
			grs = grs[:len(grs)-1]
		}
	}
	for _, g := range grs {
		groups = append(groups, &Group{
			ID:          g.ID,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
			Description: str.NullString(g.Description),
			Name:        g.Name,
			ParentGroup: api.Expandable{ID: g.ParentGroup},
		})
	}
	return groups, hasMore, err
}

func (s *GroupsService) Update(ugp database.UpdateGroupParams) (*Group, error) {
	gr, err := s.Db.UpdateGroup(context.Background(), ugp)
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
