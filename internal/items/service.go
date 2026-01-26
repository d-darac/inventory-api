package items

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/currency"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/ints"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemsService struct {
	Db *database.Queries
}

type Create struct {
	AccountId     uuid.UUID
	RequestParams CreateItemParams
}

type Delete struct {
	AccountId uuid.UUID
	ItemId    uuid.UUID
}

type Get struct {
	AccountId     uuid.UUID
	ItemId        uuid.UUID
	RequestParams RetrieveItemParams
	OmitBase      bool
}

type ListByIds struct {
	AccountId     uuid.UUID
	RequestParams ListItemsByIdsParams
}

type List struct {
	AccountId     uuid.UUID
	RequestParams ListItemsParams
}

type Update struct {
	AccountId     uuid.UUID
	ItemId        uuid.UUID
	RequestParams UpdateItemParams
}

func NewItemsService(db *database.Queries) *ItemsService {
	return &ItemsService{
		Db: db,
	}
}

func (s *ItemsService) Create(create Create) (*Item, error) {
	dbParams := MapCreateItemParams(create)
	row, err := s.Db.CreateItem(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:            &row.ID,
		CreatedAt:     &row.CreatedAt,
		UpdatedAt:     &row.UpdatedAt,
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   ints.NullInt32(row.PriceAmount),
		PriceCurrency: currency.NullCurrency(row.PriceCurrency),
		Variant:       row.Variant,
		Type:          row.Type,
	}

	return item, nil
}

func (s *ItemsService) Delete(delete Delete) error {
	_, err := s.Get(Get{
		AccountId:     delete.AccountId,
		ItemId:        delete.ItemId,
		RequestParams: RetrieveItemParams{},
		OmitBase:      true,
	})
	if err != nil {
		return err
	}
	return s.Db.DeleteItem(context.Background(), database.DeleteItemParams{
		ID:        delete.ItemId,
		AccountID: delete.AccountId,
	})
}

func (s *ItemsService) Get(get Get) (*Item, error) {
	row, err := s.Db.GetItem(context.Background(), database.GetItemParams{
		ID:        get.ItemId,
		AccountID: get.AccountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NotFoundMessage(get.ItemId, "item")
		}
		return nil, err
	}

	item := &Item{
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{ID: row.Identifiers},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   ints.NullInt32(row.PriceAmount),
		PriceCurrency: currency.NullCurrency(row.PriceCurrency),
		Variant:       row.Variant,
		Type:          row.Type,
	}

	if !get.OmitBase {
		item.ID = &row.ID
		item.CreatedAt = &row.CreatedAt
		item.UpdatedAt = &row.UpdatedAt
	}
	return item, nil
}

func (s *ItemsService) ListByIds(list ListByIds) (items []*Item, err error) {
	params := database.ListItemsByIdsParams{
		AccountID: list.AccountId,
		Ids:       list.RequestParams.Ids,
	}

	rows, err := s.Db.ListItemsByIds(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return
	}

	for _, row := range rows {
		items = append(items, &Item{
			ID:            &row.ID,
			CreatedAt:     &row.CreatedAt,
			UpdatedAt:     &row.UpdatedAt,
			Active:        row.Active,
			Description:   str.NullString(row.Description),
			Group:         api.Expandable{ID: row.Group},
			Identifiers:   api.Expandable{ID: row.Identifiers},
			Inventory:     api.Expandable{ID: row.Inventory},
			Name:          row.Name,
			PriceAmount:   ints.NullInt32(row.PriceAmount),
			PriceCurrency: currency.NullCurrency(row.PriceCurrency),
			Variant:       row.Variant,
			Type:          row.Type,
		})
	}

	return
}

func (s *ItemsService) List(list List) (items []*Item, hasMore bool, err error) {
	if list.RequestParams.StartingAfter != nil {
		item, err := s.Get(Get{
			AccountId:     list.AccountId,
			ItemId:        *list.RequestParams.StartingAfter,
			RequestParams: RetrieveItemParams{},
			OmitBase:      false,
		})
		if err != nil {
			return items, hasMore, err
		}
		list.RequestParams.StartingAfterDate = item.CreatedAt
	}

	if list.RequestParams.EndingBefore != nil {
		item, err := s.Get(Get{
			AccountId:     list.AccountId,
			ItemId:        *list.RequestParams.EndingBefore,
			RequestParams: RetrieveItemParams{},
			OmitBase:      false,
		})
		if err != nil {
			return items, hasMore, err
		}
		list.RequestParams.EndingBeforeDate = item.CreatedAt
	}

	dbParams := MapListItemsParams(list)

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
			ID:            &row.ID,
			CreatedAt:     &row.CreatedAt,
			UpdatedAt:     &row.UpdatedAt,
			Active:        row.Active,
			Description:   str.NullString(row.Description),
			Group:         api.Expandable{ID: row.Group},
			Identifiers:   api.Expandable{ID: row.Identifiers},
			Inventory:     api.Expandable{ID: row.Inventory},
			Name:          row.Name,
			PriceAmount:   ints.NullInt32(row.PriceAmount),
			PriceCurrency: currency.NullCurrency(row.PriceCurrency),
			Variant:       row.Variant,
			Type:          row.Type,
		})
	}

	return items, hasMore, err
}

func (s *ItemsService) Update(update Update) (*Item, error) {
	_, err := s.Get(Get{
		AccountId:     update.AccountId,
		ItemId:        update.ItemId,
		RequestParams: RetrieveItemParams{},
		OmitBase:      true,
	})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateItemParams(update)

	row, err := s.Db.UpdateItem(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:            &row.ID,
		CreatedAt:     &row.CreatedAt,
		UpdatedAt:     &row.UpdatedAt,
		Active:        row.Active,
		Description:   str.NullString(row.Description),
		Group:         api.Expandable{ID: row.Group},
		Identifiers:   api.Expandable{ID: row.Identifiers},
		Inventory:     api.Expandable{ID: row.Inventory},
		Name:          row.Name,
		PriceAmount:   ints.NullInt32(row.PriceAmount),
		PriceCurrency: currency.NullCurrency(row.PriceCurrency),
		Variant:       row.Variant,
		Type:          row.Type,
	}

	return item, nil
}
