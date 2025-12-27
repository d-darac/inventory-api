package items

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemsService struct {
	Db *database.Queries
}

func NewItemsService(db *database.Queries) *ItemsService {
	return &ItemsService{
		Db: db,
	}
}

func (s *ItemsService) Create(accountId uuid.UUID, params *CreateItemParams) (*Item, error) {
	dbParams := MapCreateItemParams(accountId, params)
	row, err := s.Db.CreateItem(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:            row.ID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   row.PriceAmount,
		PriceCurrency: row.PriceCurrency,
		Variant:       row.Variant,
		Type:          row.Type,
	}

	return item, nil
}

func (s *ItemsService) Delete(itemId, accountId uuid.UUID) error {
	_, err := s.Get(itemId, accountId, &RetrieveItemParams{})
	if err != nil {
		return err
	}
	return s.Db.DeleteItem(context.Background(), database.DeleteItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
}

func (s *ItemsService) Get(itemId, accountId uuid.UUID, params *RetrieveItemParams) (*Item, error) {
	row, err := s.Db.GetItem(context.Background(), database.GetItemParams{
		ID:        itemId,
		AccountID: accountId,
	})
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:            row.ID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   row.PriceAmount,
		PriceCurrency: row.PriceCurrency,
		Variant:       row.Variant,
		Type:          row.Type,
	}

	return item, nil
}

func (s *ItemsService) List(accountId uuid.UUID, params *ListItemsParams) (items []*Item, hasMore bool, err error) {
	if params.StartingAfter != nil {
		item, err := s.Get(*params.StartingAfter, accountId, &RetrieveItemParams{})
		if err != nil {
			return items, hasMore, err
		}
		params.StartingAfterDate = &item.CreatedAt
	}

	if params.EndingBefore != nil {
		item, err := s.Get(*params.EndingBefore, accountId, &RetrieveItemParams{})
		if err != nil {
			return items, hasMore, err
		}
		params.EndingBeforeDate = &item.CreatedAt
	}

	dbParams := MapListItemsParams(accountId, params)

	rows, err := s.Db.ListItems(context.Background(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, hasMore, nil
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
		items = append(items, &Item{
			ID:            row.ID,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
			Active:        row.Active,
			Description:   str.NullString(row.Description),
			Group:         api.Expandable{ID: row.Group},
			Identifiers:   api.Expandable{},
			Inventory:     api.Expandable{ID: row.Inventory},
			Name:          row.Name,
			PriceAmount:   row.PriceAmount,
			PriceCurrency: row.PriceCurrency,
			Variant:       row.Variant,
			Type:          row.Type,
		})
	}

	return items, hasMore, err
}

func (s *ItemsService) Update(itemId, accountId uuid.UUID, params *UpdateItemParams) (*Item, error) {
	_, err := s.Get(itemId, accountId, &RetrieveItemParams{})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateItemParams(itemId, accountId, params)

	row, err := s.Db.UpdateItem(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:            row.ID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   row.PriceAmount,
		PriceCurrency: row.PriceCurrency,
		Variant:       row.Variant,
		Type:          row.Type,
	}

	return item, nil
}
