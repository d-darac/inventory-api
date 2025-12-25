package services

import (
	"context"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

type ItemIdentifier struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Ean       str.NullString `json:"ean"`
	Gtin      str.NullString `json:"gtin"`
	Isbn      str.NullString `json:"isbn"`
	Jan       str.NullString `json:"jan"`
	Mpn       str.NullString `json:"mpn"`
	Nsn       str.NullString `json:"nsn"`
	Upc       str.NullString `json:"upc"`
	Qr        str.NullString `json:"qr"`
	Sku       str.NullString `json:"sku"`
	Item      api.Expandable `json:"item"`
}

type ItemIdentifiersService struct {
	Db *database.Queries
}

func NewItemIdentifiersService(db *database.Queries) *ItemIdentifiersService {
	return &ItemIdentifiersService{
		Db: db,
	}
}

func (s *ItemIdentifiersService) Create(ciip database.CreateItemIdentifierParams) (*ItemIdentifier, error) {
	ii, err := s.Db.CreateItemIdentifier(context.Background(), ciip)
	if err != nil {
		return nil, err
	}
	return &ItemIdentifier{
		ID:        ii.ID,
		CreatedAt: ii.CreatedAt,
		UpdatedAt: ii.UpdatedAt,
		Ean:       str.NullString(ii.Ean),
		Gtin:      str.NullString(ii.Gtin),
		Isbn:      str.NullString(ii.Isbn),
		Jan:       str.NullString(ii.Jan),
		Mpn:       str.NullString(ii.Mpn),
		Nsn:       str.NullString(ii.Nsn),
		Upc:       str.NullString(ii.Upc),
		Qr:        str.NullString(ii.Qr),
		Sku:       str.NullString(ii.Sku),
	}, nil
}

func (s *ItemIdentifiersService) Delete(diip database.DeleteItemIdentifierParams) error {
	return s.Db.DeleteItemIdentifier(context.Background(), diip)
}

func (s *ItemIdentifiersService) Get(giip database.GetItemIdentifierParams) (*ItemIdentifier, error) {
	ii, err := s.Db.GetItemIdentifier(context.Background(), giip)
	if err != nil {
		return nil, err
	}
	return &ItemIdentifier{
		ID:        ii.ID,
		CreatedAt: ii.CreatedAt,
		UpdatedAt: ii.UpdatedAt,
		Ean:       str.NullString(ii.Ean),
		Gtin:      str.NullString(ii.Gtin),
		Isbn:      str.NullString(ii.Isbn),
		Jan:       str.NullString(ii.Jan),
		Mpn:       str.NullString(ii.Mpn),
		Nsn:       str.NullString(ii.Nsn),
		Upc:       str.NullString(ii.Upc),
		Qr:        str.NullString(ii.Qr),
		Sku:       str.NullString(ii.Sku),
	}, nil
}

func (s *ItemIdentifiersService) List(liip database.ListItemIdentifiersParams) (itemIdentifiers []*ItemIdentifier, hasMore bool, err error) {
	iis, err := s.Db.ListItemIdentifiers(context.Background(), liip)
	if err != nil {
		return
	}
	if liip.Limit.Valid {
		hasMore = len(iis) > int(liip.Limit.Int32)
	} else {
		hasMore = len(iis) > 10
	}
	if hasMore {
		if liip.EndingBefore.Valid {
			iis = iis[1:]
		} else {
			iis = iis[:len(iis)-1]
		}
	}
	for _, ii := range iis {
		itemIdentifiers = append(itemIdentifiers, &ItemIdentifier{
			ID:        ii.ID,
			CreatedAt: ii.CreatedAt,
			UpdatedAt: ii.UpdatedAt,
			Ean:       str.NullString(ii.Ean),
			Gtin:      str.NullString(ii.Gtin),
			Isbn:      str.NullString(ii.Isbn),
			Jan:       str.NullString(ii.Jan),
			Mpn:       str.NullString(ii.Mpn),
			Nsn:       str.NullString(ii.Nsn),
			Upc:       str.NullString(ii.Upc),
			Qr:        str.NullString(ii.Qr),
			Sku:       str.NullString(ii.Sku),
		})
	}
	return
}

func (s *ItemIdentifiersService) Update(uiip database.UpdateItemIdentifierParams) (*ItemIdentifier, error) {
	ii, err := s.Db.UpdateItemIdentifier(context.Background(), uiip)
	if err != nil {
		return nil, err
	}
	return &ItemIdentifier{
		ID:        ii.ID,
		CreatedAt: ii.CreatedAt,
		UpdatedAt: ii.UpdatedAt,
		Ean:       str.NullString(ii.Ean),
		Gtin:      str.NullString(ii.Gtin),
		Isbn:      str.NullString(ii.Isbn),
		Jan:       str.NullString(ii.Jan),
		Mpn:       str.NullString(ii.Mpn),
		Nsn:       str.NullString(ii.Nsn),
		Upc:       str.NullString(ii.Upc),
		Qr:        str.NullString(ii.Qr),
		Sku:       str.NullString(ii.Sku),
	}, nil
}
