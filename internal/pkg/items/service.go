package items

import (
	"context"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
)

type ItemsService struct {
	Db *database.Queries
}

func NewItemsService(db *database.Queries) *ItemsService {
	return &ItemsService{
		Db: db,
	}
}

func (s *ItemsService) Create(cip database.CreateItemParams) (*Item, error) {
	it, err := s.Db.CreateItem(context.Background(), cip)
	if err != nil {
		return nil, err
	}
	return &Item{
		ID:            it.ID,
		CreatedAt:     it.CreatedAt,
		UpdatedAt:     it.UpdatedAt,
		Active:        it.Active,
		Description:   str.NullString(it.Description),
		Group:         api.Expandable{ID: it.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: it.Inventory},
		Name:          it.Name,
		PriceAmount:   it.PriceAmount,
		PriceCurrency: it.PriceCurrency,
		Variant:       it.Variant,
		Type:          it.Type,
	}, nil
}

func (s *ItemsService) Delete(dip database.DeleteItemParams) error {
	return s.Db.DeleteItem(context.Background(), dip)
}

func (s *ItemsService) Get(gip database.GetItemParams) (*Item, error) {
	it, err := s.Db.GetItem(context.Background(), gip)
	if err != nil {
		return nil, err
	}
	return &Item{
		ID:            it.ID,
		CreatedAt:     it.CreatedAt,
		UpdatedAt:     it.UpdatedAt,
		Active:        it.Active,
		Description:   str.NullString(it.Description),
		Group:         api.Expandable{ID: it.Group},
		Identifiers:   api.Expandable{},
		Inventory:     api.Expandable{ID: it.Inventory},
		Name:          it.Name,
		PriceAmount:   it.PriceAmount,
		PriceCurrency: it.PriceCurrency,
		Variant:       it.Variant,
		Type:          it.Type,
	}, nil
}

func (s *ItemsService) List(lip database.ListItemsParams) (items []*Item, hasMore bool, err error) {
	its, err := s.Db.ListItems(context.Background(), lip)
	if err != nil {
		return
	}
	if lip.Limit.Valid {
		hasMore = len(its) > int(lip.Limit.Int32)
	} else {
		hasMore = len(its) > 10
	}
	if hasMore {
		if lip.EndingBefore.Valid {
			its = its[1:]
		} else {
			its = its[:len(its)-1]
		}
	}
	for _, it := range its {
		items = append(items, &Item{
			ID:            it.ID,
			CreatedAt:     it.CreatedAt,
			UpdatedAt:     it.UpdatedAt,
			Active:        it.Active,
			Description:   str.NullString(it.Description),
			Group:         api.Expandable{ID: it.Group},
			Identifiers:   api.Expandable{ID: it.Identifiers},
			Inventory:     api.Expandable{ID: it.Inventory},
			Name:          it.Name,
			PriceAmount:   it.PriceAmount,
			PriceCurrency: it.PriceCurrency,
			Variant:       it.Variant,
			Type:          it.Type,
		})
	}
	return
}

func (s *ItemsService) Update(uip database.UpdateItemParams) (*Item, error) {
	it, err := s.Db.UpdateItem(context.Background(), uip)
	if err != nil {
		return nil, err
	}
	return &Item{
		ID:            it.ID,
		CreatedAt:     it.CreatedAt,
		UpdatedAt:     it.UpdatedAt,
		Active:        it.Active,
		Description:   str.NullString(it.Description),
		Group:         api.Expandable{ID: it.Group},
		Identifiers:   api.Expandable{ID: it.Identifiers},
		Inventory:     api.Expandable{ID: it.Inventory},
		Name:          it.Name,
		PriceAmount:   it.PriceAmount,
		PriceCurrency: it.PriceCurrency,
		Variant:       it.Variant,
		Type:          it.Type,
	}, nil
}
