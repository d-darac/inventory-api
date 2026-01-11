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

type Create struct {
	AccountId     uuid.UUID
	RequestParams CreateGroupParams
}

type Delete struct {
	AccountId uuid.UUID
	GroupId   uuid.UUID
}

type Get struct {
	AccountId     uuid.UUID
	GroupId       uuid.UUID
	RequestParams RetrieveGroupParams
	OmitBase      bool
}

type List struct {
	AccountId     uuid.UUID
	RequestParams ListGroupsParams
}

type Update struct {
	AccountId     uuid.UUID
	GroupId       uuid.UUID
	RequestParams UpdateGroupParams
}

func NewGroupsService(db *database.Queries) *GroupsService {
	return &GroupsService{
		Db: db,
	}
}

func (s *GroupsService) Create(create Create) (*Group, error) {
	dbParams := MapCreateGroupParams(create)
	row, err := s.Db.CreateGroup(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	group := &Group{
		ID:          &row.ID,
		CreatedAt:   &row.CreatedAt,
		UpdatedAt:   &row.UpdatedAt,
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	return group, nil
}

func (s *GroupsService) Delete(delete Delete) error {
	_, err := s.Get(Get{
		AccountId:     delete.AccountId,
		GroupId:       delete.GroupId,
		RequestParams: RetrieveGroupParams{},
		OmitBase:      true,
	})
	if err != nil {
		return err
	}
	return s.Db.DeleteGroup(context.Background(), database.DeleteGroupParams{
		ID:        delete.GroupId,
		AccountID: delete.AccountId,
	})
}

func (s *GroupsService) Get(get Get) (*Group, error) {
	row, err := s.Db.GetGroup(context.Background(), database.GetGroupParams{
		ID:        get.GroupId,
		AccountID: get.AccountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NotFoundMessage(get.GroupId, "group")
		}
		return nil, err
	}

	group := &Group{
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	if !get.OmitBase {
		group.ID = &row.ID
		group.CreatedAt = &row.CreatedAt
		group.UpdatedAt = &row.UpdatedAt
	}

	return group, nil
}

func (s *GroupsService) List(list List) (groups []*Group, hasMore bool, err error) {
	if list.RequestParams.StartingAfter != nil {
		group, err := s.Get(Get{
			AccountId:     list.AccountId,
			GroupId:       *list.RequestParams.StartingAfter,
			RequestParams: RetrieveGroupParams{},
			OmitBase:      false,
		})
		if err != nil {
			return groups, hasMore, err
		}
		list.RequestParams.StartingAfterDate = group.CreatedAt
	}

	if list.RequestParams.EndingBefore != nil {
		group, err := s.Get(Get{
			AccountId:     list.AccountId,
			GroupId:       *list.RequestParams.EndingBefore,
			RequestParams: RetrieveGroupParams{},
			OmitBase:      false,
		})
		if err != nil {
			return groups, hasMore, err
		}
		list.RequestParams.EndingBeforeDate = group.CreatedAt
	}

	dbParams := MapListGroupsParams(list)

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
			ID:          &row.ID,
			CreatedAt:   &row.CreatedAt,
			UpdatedAt:   &row.UpdatedAt,
			Description: str.NullString(row.Description),
			Name:        row.Name,
			ParentGroup: api.Expandable{ID: row.ParentGroup},
		})
	}

	return groups, hasMore, err
}

func (s *GroupsService) Update(update Update) (*Group, error) {
	_, err := s.Get(Get{
		AccountId:     update.AccountId,
		GroupId:       update.GroupId,
		RequestParams: RetrieveGroupParams{},
		OmitBase:      true,
	})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateGroupParams(update)

	row, err := s.Db.UpdateGroup(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	group := &Group{
		ID:          &row.ID,
		CreatedAt:   &row.CreatedAt,
		UpdatedAt:   &row.UpdatedAt,
		Description: str.NullString(row.Description),
		Name:        row.Name,
		ParentGroup: api.Expandable{ID: row.ParentGroup},
	}

	return group, nil
}
