package groups

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type GroupsService struct {
	Db *database.Queries
}

func NewGroupsService(db *database.Queries) *GroupsService {
	return &GroupsService{
		Db: db,
	}
}

func (s *GroupsService) Create(accountId uuid.UUID, params *CreateGroupParams) (*Group, error) {
	dbParams := MapCreateGroupParams(accountId, params)
	row, err := s.Db.CreateGroup(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	group := &Group{
		ID:          row.ID,
		CreatedAt:   row.UpdatedAt,
		UpdatedAt:   row.UpdatedAt,
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	return group, nil
}

func (s *GroupsService) Delete(groupId, accountId uuid.UUID) error {
	_, err := s.Get(groupId, accountId, &RetrieveGroupParams{})
	if err != nil {
		return err
	}
	return s.Db.DeleteGroup(context.Background(), database.DeleteGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
}

func (s *GroupsService) Get(groupId, accountId uuid.UUID, params *RetrieveGroupParams) (*Group, error) {
	row, err := s.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        groupId,
		AccountID: accountId,
	})
	if err != nil {
		return nil, err
	}

	group := &Group{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	return group, nil
}

func (s *GroupsService) List(accountId uuid.UUID, params *ListGroupsParams) (groups []*Group, hasMore bool, err error) {
	if params.StartingAfter != nil {
		group, err := s.Get(*params.StartingAfter, accountId, &RetrieveGroupParams{})
		if err != nil {
			return groups, hasMore, err
		}
		params.StartingAfterDate = &group.CreatedAt
	}

	if params.EndingBefore != nil {
		group, err := s.Get(*params.EndingBefore, accountId, &RetrieveGroupParams{})
		if err != nil {
			return groups, hasMore, err
		}
		params.EndingBeforeDate = &group.CreatedAt
	}

	dbParams := MapListGroupsParams(accountId, params)

	rows, err := s.Db.ListGroups(context.Background(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return groups, hasMore, nil
		}
		return
	}

	if dbParams.Limit.Valid {
		hasMore = len(rows) > int(dbParams.Limit.Int32)
	} else {
		hasMore = len(rows) > 10
	}

	if hasMore {
		if dbParams.EndingBefore.Valid {
			rows = rows[1:]
		} else {
			rows = rows[:len(rows)-1]
		}
	}

	for _, row := range rows {
		groups = append(groups, &Group{
			ID:          row.ID,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			Description: str.NullString(row.Description),
			Name:        row.Name,
			ParentGroup: api.Expandable{ID: row.ParentGroup},
		})
	}

	return groups, hasMore, err
}

func (s *GroupsService) Update(groupId, accountId uuid.UUID, params *UpdateGroupParams) (*Group, error) {
	_, err := s.Get(groupId, accountId, &RetrieveGroupParams{})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateGroupParams(groupId, accountId, params)

	row, err := s.Db.UpdateGroup(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	group := &Group{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	return group, nil
}
