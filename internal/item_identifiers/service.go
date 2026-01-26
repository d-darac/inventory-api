package itemidentifiers

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemIdentifiersService struct {
	Db *database.Queries
}

type Create struct {
	AccountId     uuid.UUID
	RequestParams CreateItemIdentifiersParams
}

type Delete struct {
	AccountId         uuid.UUID
	ItemIdentifiersId uuid.UUID
}

type Get struct {
	AccountId         uuid.UUID
	ItemIdentifiersId uuid.UUID
	RequestParams     RetrieveItemIdentifiersParams
	OmitBase          bool
}

type ListByIds struct {
	AccountId     uuid.UUID
	RequestParams ListItemIdentifiersByIdsParams
}

type List struct {
	AccountId     uuid.UUID
	RequestParams ListItemIdentifiersParams
}

type Update struct {
	AccountId         uuid.UUID
	ItemIdentifiersId uuid.UUID
	RequestParams     UpdateItemIdentifiersParams
}

func NewItemIdentifiersService(db *database.Queries) *ItemIdentifiersService {
	return &ItemIdentifiersService{
		Db: db,
	}
}

func (s *ItemIdentifiersService) Create(create Create) (*ItemIdentifiers, error) {
	dbParams := MapCreateItemIdentifiersParams(create)
	row, err := s.Db.CreateItemIdentifier(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		ID:        &row.ID,
		CreatedAt: &row.CreatedAt,
		UpdatedAt: &row.UpdatedAt,
		Ean:       str.NullString(row.Ean),
		Gtin:      str.NullString(row.Gtin),
		Isbn:      str.NullString(row.Isbn),
		Jan:       str.NullString(row.Jan),
		Mpn:       str.NullString(row.Mpn),
		Nsn:       str.NullString(row.Nsn),
		Upc:       str.NullString(row.Upc),
		Qr:        str.NullString(row.Qr),
		Sku:       str.NullString(row.Sku),
	}

	return itemIdentifiers, nil
}

func (s *ItemIdentifiersService) Delete(delete Delete) error {
	_, err := s.Get(Get{
		AccountId:         delete.AccountId,
		ItemIdentifiersId: delete.ItemIdentifiersId,
		RequestParams:     RetrieveItemIdentifiersParams{},
		OmitBase:          true,
	})
	if err != nil {
		return err
	}
	return s.Db.DeleteItemIdentifier(context.Background(), database.DeleteItemIdentifierParams{
		ID:        delete.ItemIdentifiersId,
		AccountID: delete.AccountId,
	})
}

func (s *ItemIdentifiersService) Get(get Get) (*ItemIdentifiers, error) {
	row, err := s.Db.GetItemIdentifier(context.Background(), database.GetItemIdentifierParams{
		ID:        get.ItemIdentifiersId,
		AccountID: get.AccountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NotFoundMessage(get.ItemIdentifiersId, "item identifiers")
		}
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		Ean:  str.NullString(row.Ean),
		Gtin: str.NullString(row.Gtin),
		Isbn: str.NullString(row.Isbn),
		Jan:  str.NullString(row.Jan),
		Mpn:  str.NullString(row.Mpn),
		Nsn:  str.NullString(row.Nsn),
		Upc:  str.NullString(row.Upc),
		Qr:   str.NullString(row.Qr),
		Sku:  str.NullString(row.Sku),
	}

	if !get.OmitBase {
		itemIdentifiers.ID = &row.ID
		itemIdentifiers.CreatedAt = &row.CreatedAt
		itemIdentifiers.UpdatedAt = &row.UpdatedAt
	}
	return itemIdentifiers, nil
}

func (s *ItemIdentifiersService) ListByIds(list ListByIds) (itemIdentifiers []*ItemIdentifiers, err error) {
	params := database.ListItemIdentifiersByIdsParams{
		AccountID: list.AccountId,
		Ids:       list.RequestParams.Ids,
	}

	rows, err := s.Db.ListItemIdentifiersByIds(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return itemIdentifiers, nil
		}
		return
	}

	for _, row := range rows {
		itemIdentifiers = append(itemIdentifiers, &ItemIdentifiers{
			ID:        &row.ID,
			CreatedAt: &row.CreatedAt,
			UpdatedAt: &row.UpdatedAt,
			Ean:       str.NullString(row.Ean),
			Gtin:      str.NullString(row.Gtin),
			Isbn:      str.NullString(row.Isbn),
			Jan:       str.NullString(row.Jan),
			Mpn:       str.NullString(row.Mpn),
			Nsn:       str.NullString(row.Nsn),
			Upc:       str.NullString(row.Upc),
			Qr:        str.NullString(row.Qr),
			Sku:       str.NullString(row.Sku),
		})
	}

	return
}

func (s *ItemIdentifiersService) List(list List) (itemIdentifiers []*ItemIdentifiers, hasMore bool, err error) {
	if list.RequestParams.StartingAfter != nil {
		iis, err := s.Get(Get{
			AccountId:         list.AccountId,
			ItemIdentifiersId: *list.RequestParams.StartingAfter,
			RequestParams:     RetrieveItemIdentifiersParams{},
			OmitBase:          false,
		})
		if err != nil {
			return itemIdentifiers, hasMore, err
		}
		list.RequestParams.StartingAfterDate = iis.CreatedAt
	}

	if list.RequestParams.EndingBefore != nil {
		iis, err := s.Get(Get{
			AccountId:         list.AccountId,
			ItemIdentifiersId: *list.RequestParams.EndingBefore,
			RequestParams:     RetrieveItemIdentifiersParams{},
			OmitBase:          false,
		})
		if err != nil {
			return itemIdentifiers, hasMore, err
		}
		list.RequestParams.EndingBeforeDate = iis.CreatedAt
	}

	dbParams := MapListItemIdentifiersParams(list)

	rows, err := s.Db.ListItemIdentifiers(context.Background(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return itemIdentifiers, hasMore, nil
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
		itemIdentifiers = append(itemIdentifiers, &ItemIdentifiers{
			ID:        &row.ID,
			CreatedAt: &row.CreatedAt,
			UpdatedAt: &row.UpdatedAt,
			Ean:       str.NullString(row.Ean),
			Gtin:      str.NullString(row.Gtin),
			Isbn:      str.NullString(row.Isbn),
			Jan:       str.NullString(row.Jan),
			Mpn:       str.NullString(row.Mpn),
			Nsn:       str.NullString(row.Nsn),
			Upc:       str.NullString(row.Upc),
			Qr:        str.NullString(row.Qr),
			Sku:       str.NullString(row.Sku),
		})
	}

	return itemIdentifiers, hasMore, err
}

func (s *ItemIdentifiersService) Update(update Update) (*ItemIdentifiers, error) {
	_, err := s.Get(Get{
		AccountId:         update.AccountId,
		ItemIdentifiersId: update.ItemIdentifiersId,
		RequestParams:     RetrieveItemIdentifiersParams{},
		OmitBase:          true,
	})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateItemIdentifiersParams(update)

	row, err := s.Db.UpdateItemIdentifier(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		ID:        &row.ID,
		CreatedAt: &row.CreatedAt,
		UpdatedAt: &row.UpdatedAt,
		Ean:       str.NullString(row.Ean),
		Gtin:      str.NullString(row.Gtin),
		Isbn:      str.NullString(row.Isbn),
		Jan:       str.NullString(row.Jan),
		Mpn:       str.NullString(row.Mpn),
		Nsn:       str.NullString(row.Nsn),
		Upc:       str.NullString(row.Upc),
		Qr:        str.NullString(row.Qr),
		Sku:       str.NullString(row.Sku),
	}

	return itemIdentifiers, nil
}
