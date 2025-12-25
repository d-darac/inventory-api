package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	InStock   int32          `json:"in_stock"`
	Orderable sql.NullInt32  `json:"orderable"`
	Item      api.Expandable `json:"item"`
}

type InventoriesService struct {
	Db *database.Queries
}

func NewInventoriesService(db *database.Queries) *InventoriesService {
	return &InventoriesService{
		Db: db,
	}
}

func (s *InventoriesService) Create(cip database.CreateInventoryParams) (*Inventory, error) {
	in, err := s.Db.CreateInventory(context.Background(), cip)
	if err != nil {
		return nil, err
	}
	return &Inventory{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		InStock:   in.InStock,
		Orderable: in.Orderable,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: in.Item, Valid: true}},
	}, nil
}

func (s *InventoriesService) Delete(dip database.DeleteInventoryParams) error {
	return s.Db.DeleteInventory(context.Background(), dip)
}

func (s *InventoriesService) Get(gip database.GetInventoryParams) (*Inventory, error) {
	in, err := s.Db.GetInventory(context.Background(), gip)
	if err != nil {
		return nil, err
	}
	return &Inventory{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		InStock:   in.InStock,
		Orderable: in.Orderable,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: in.Item, Valid: true}},
	}, nil
}

func (s *InventoriesService) List(lip database.ListInventoriesParams) (inventories []*Inventory, hasMore bool, err error) {
	ins, err := s.Db.ListInventories(context.Background(), lip)
	if err != nil {
		return
	}
	if lip.Limit.Valid {
		hasMore = len(ins) > int(lip.Limit.Int32)
	} else {
		hasMore = len(ins) > 10
	}
	if hasMore {
		if lip.EndingBefore.Valid {
			ins = ins[1:]
		} else {
			ins = ins[:len(ins)-1]
		}
	}
	for _, i := range ins {
		inventories = append(inventories, &Inventory{
			ID:        i.ID,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			InStock:   i.InStock,
			Orderable: i.Orderable,
			Item:      api.Expandable{ID: uuid.NullUUID{UUID: i.Item, Valid: true}},
		})
	}
	return
}

func (s *InventoriesService) Update(uip database.UpdateInventoryParams) (*Inventory, error) {
	in, err := s.Db.UpdateInventory(context.Background(), uip)
	if err != nil {
		return nil, err
	}
	return &Inventory{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		InStock:   in.InStock,
		Orderable: in.Orderable,
		Item:      api.Expandable{ID: uuid.NullUUID{UUID: in.Item, Valid: true}},
	}, nil
}
