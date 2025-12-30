package itemidentifiers

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemIdentifiersService struct {
	Db *database.Queries
}

func NewItemIdentifiersService(db *database.Queries) *ItemIdentifiersService {
	return &ItemIdentifiersService{
		Db: db,
	}
}

func (s *ItemIdentifiersService) Create(accountId uuid.UUID, params *CreateItemIdentifiersParams) (*ItemIdentifiers, error) {
	dbParams := MapCreateItemIdentifiersParams(accountId, params)
	row, err := s.Db.CreateItemIdentifier(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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

func (s *ItemIdentifiersService) Delete(itemIdentifiersId, accountId uuid.UUID) error {
	_, err := s.Get(itemIdentifiersId, accountId, &RetrieveItemIdentifiersParams{})
	if err != nil {
		return err
	}
	return s.Db.DeleteItemIdentifier(context.Background(), database.DeleteItemIdentifierParams{
		ID:        itemIdentifiersId,
		AccountID: accountId,
	})
}

func (s *ItemIdentifiersService) Get(itemIdentifiersId, accountId uuid.UUID, params *RetrieveItemIdentifiersParams) (*ItemIdentifiers, error) {
	row, err := s.Db.GetItemIdentifier(context.Background(), database.GetItemIdentifierParams{
		ID:        itemIdentifiersId,
		AccountID: accountId,
	})
	if err != nil {
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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

func (s *ItemIdentifiersService) List(accountId uuid.UUID, params *ListItemIdentifiersParams) (itemIdentifiers []*ItemIdentifiers, hasMore bool, err error) {
	if params.StartingAfter != nil {
		iis, err := s.Get(*params.StartingAfter, accountId, &RetrieveItemIdentifiersParams{})
		if err != nil {
			return itemIdentifiers, hasMore, err
		}
		params.StartingAfterDate = &iis.CreatedAt
	}

	if params.EndingBefore != nil {
		iis, err := s.Get(*params.EndingBefore, accountId, &RetrieveItemIdentifiersParams{})
		if err != nil {
			return itemIdentifiers, hasMore, err
		}
		params.EndingBeforeDate = &iis.CreatedAt
	}

	dbParams := MapListItemIdentifiersParams(accountId, params)

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
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
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

func (s *ItemIdentifiersService) Update(itemIdentifiersId, accountId uuid.UUID, params *UpdateItemIdentifiersParams) (*ItemIdentifiers, error) {
	_, err := s.Get(itemIdentifiersId, accountId, &RetrieveItemIdentifiersParams{})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateItemIdentifiersParams(itemIdentifiersId, accountId, params)

	row, err := s.Db.UpdateItemIdentifier(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	itemIdentifiers := &ItemIdentifiers{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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
