package inventories

import (
	"context"
	"database/sql"

	"github.com/d-darac/inventory-assets/api"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/ints"
	"github.com/google/uuid"
)

type InventoriesService struct {
	Db *database.Queries
}

func NewInventoriesService(db *database.Queries) *InventoriesService {
	return &InventoriesService{
		Db: db,
	}
}

func (s *InventoriesService) Create(accountId uuid.UUID, params *CreateInventoryParams) (*Inventory, error) {
	dbParams := MapCreateInventoryParams(accountId, params)
	row, err := s.Db.CreateInventory(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	Inventory := &Inventory{
		ID:        row.ID,
		CreatedAt: row.UpdatedAt,
		UpdatedAt: row.UpdatedAt,
		InStock:   row.InStock,
		Orderable: ints.NullInt32{
			Int32: row.Orderable.Int32,
			Valid: true,
		},
	}

	return Inventory, nil
}

func (s *InventoriesService) Delete(inventoryId, accountId uuid.UUID) error {
	_, err := s.Get(inventoryId, accountId, &RetrieveInventoryParams{})
	if err != nil {
		return err
	}
	return s.Db.DeleteInventory(context.Background(), database.DeleteInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
}

func (s *InventoriesService) Get(inventoryId, accountId uuid.UUID, params *RetrieveInventoryParams) (*Inventory, error) {
	row, err := s.Db.GetInventory(context.Background(), database.GetInventoryParams{
		ID:        inventoryId,
		AccountID: accountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NotFoundMessage(inventoryId, "inventory")
		}
		return nil, err
	}

	Inventory := &Inventory{
		ID:        row.ID,
		CreatedAt: row.UpdatedAt,
		UpdatedAt: row.UpdatedAt,
		InStock:   row.InStock,
		Orderable: ints.NullInt32{
			Int32: row.Orderable.Int32,
			Valid: true,
		},
	}

	return Inventory, nil
}

func (s *InventoriesService) List(accountId uuid.UUID, params *ListInventoriesParams) (Inventories []*Inventory, hasMore bool, err error) {
	if params.StartingAfter != nil {
		Inventory, err := s.Get(*params.StartingAfter, accountId, &RetrieveInventoryParams{})
		if err != nil {
			return Inventories, hasMore, err
		}
		params.StartingAfterDate = &Inventory.CreatedAt
	}

	if params.EndingBefore != nil {
		Inventory, err := s.Get(*params.EndingBefore, accountId, &RetrieveInventoryParams{})
		if err != nil {
			return Inventories, hasMore, err
		}
		params.EndingBeforeDate = &Inventory.CreatedAt
	}

	dbParams := MapListInventoriesParams(accountId, params)

	rows, err := s.Db.ListInventories(context.Background(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return Inventories, hasMore, nil
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
		Inventories = append(Inventories, &Inventory{
			ID:        row.ID,
			CreatedAt: row.UpdatedAt,
			UpdatedAt: row.UpdatedAt,
			InStock:   row.InStock,
			Orderable: ints.NullInt32{
				Int32: row.Orderable.Int32,
				Valid: true,
			},
		})
	}

	return Inventories, hasMore, err
}

func (s *InventoriesService) Update(inventoryId, accountId uuid.UUID, params *UpdateInventoryParams) (*Inventory, error) {
	_, err := s.Get(inventoryId, accountId, &RetrieveInventoryParams{})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateInventoryParams(inventoryId, accountId, params)

	row, err := s.Db.UpdateInventory(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	Inventory := &Inventory{
		ID:        row.ID,
		CreatedAt: row.UpdatedAt,
		UpdatedAt: row.UpdatedAt,
		InStock:   row.InStock,
		Orderable: ints.NullInt32{
			Int32: row.Orderable.Int32,
			Valid: true,
		},
	}

	return Inventory, nil
}
