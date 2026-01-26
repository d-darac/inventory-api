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

type Create struct {
	AccountId     uuid.UUID
	RequestParams CreateInventoryParams
}

type Delete struct {
	AccountId   uuid.UUID
	InventoryId uuid.UUID
}

type Get struct {
	AccountId     uuid.UUID
	InventoryId   uuid.UUID
	RequestParams RetrieveInventoryParams
	OmitBase      bool
}

type List struct {
	AccountId     uuid.UUID
	RequestParams ListInventoriesParams
}

type Update struct {
	AccountId     uuid.UUID
	InventoryId   uuid.UUID
	RequestParams UpdateInventoryParams
}

func NewInventoriesService(db *database.Queries) *InventoriesService {
	return &InventoriesService{
		Db: db,
	}
}

func (s *InventoriesService) Create(create Create) (*Inventory, error) {
	dbParams := MapCreateInventoryParams(create)
	row, err := s.Db.CreateInventory(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	inventory := &Inventory{
		ID:        &row.ID,
		CreatedAt: &row.CreatedAt,
		UpdatedAt: &row.UpdatedAt,
		InStock:   row.InStock,
		Orderable: ints.NullInt32(row.Orderable),
	}

	return inventory, nil
}

func (s *InventoriesService) Delete(delete Delete) error {
	_, err := s.Get(Get{
		AccountId:     delete.AccountId,
		InventoryId:   delete.InventoryId,
		RequestParams: RetrieveInventoryParams{},
		OmitBase:      true,
	})
	if err != nil {
		return err
	}
	return s.Db.DeleteInventory(context.Background(), database.DeleteInventoryParams{
		ID:        delete.InventoryId,
		AccountID: delete.AccountId,
	})
}

func (s *InventoriesService) Get(get Get) (*Inventory, error) {
	row, err := s.Db.GetInventory(context.Background(), database.GetInventoryParams{
		ID:        get.InventoryId,
		AccountID: get.AccountId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NotFoundMessage(get.InventoryId, "inventory")
		}
		return nil, err
	}

	inventory := &Inventory{
		InStock:   row.InStock,
		Orderable: ints.NullInt32(row.Orderable),
	}

	if !get.OmitBase {
		inventory.ID = &row.ID
		inventory.CreatedAt = &row.CreatedAt
		inventory.UpdatedAt = &row.UpdatedAt
	}

	return inventory, nil
}

func (s *InventoriesService) List(list List) (inventories []*Inventory, hasMore bool, err error) {
	if list.RequestParams.StartingAfter != nil {
		Inventory, err := s.Get(Get{
			AccountId:     list.AccountId,
			InventoryId:   *list.RequestParams.StartingAfter,
			RequestParams: RetrieveInventoryParams{},
			OmitBase:      false,
		})
		if err != nil {
			return inventories, hasMore, err
		}
		list.RequestParams.StartingAfterDate = Inventory.CreatedAt
	}

	if list.RequestParams.EndingBefore != nil {
		Inventory, err := s.Get(Get{
			AccountId:     list.AccountId,
			InventoryId:   *list.RequestParams.EndingBefore,
			RequestParams: RetrieveInventoryParams{},
			OmitBase:      false,
		})
		if err != nil {
			return inventories, hasMore, err
		}
		list.RequestParams.EndingBeforeDate = Inventory.CreatedAt
	}

	dbParams := MapListInventoriesParams(list)

	rows, err := s.Db.ListInventories(context.Background(), dbParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return inventories, hasMore, nil
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
		inventories = append(inventories, &Inventory{
			ID:        &row.ID,
			CreatedAt: &row.CreatedAt,
			UpdatedAt: &row.UpdatedAt,
			InStock:   row.InStock,
			Orderable: ints.NullInt32(row.Orderable),
			Reserved:  ints.NullInt32(row.Reserved),
		})
	}

	return inventories, hasMore, err
}

func (s *InventoriesService) Update(update Update) (*Inventory, error) {
	_, err := s.Get(Get{
		AccountId:     update.AccountId,
		InventoryId:   update.InventoryId,
		RequestParams: RetrieveInventoryParams{},
		OmitBase:      true,
	})
	if err != nil {
		return nil, err
	}

	dbParams := MapUpdateInventoryParams(update)

	row, err := s.Db.UpdateInventory(context.Background(), dbParams)
	if err != nil {
		return nil, err
	}

	inventory := &Inventory{
		ID:        &row.ID,
		CreatedAt: &row.CreatedAt,
		UpdatedAt: &row.UpdatedAt,
		InStock:   row.InStock,
		Orderable: ints.NullInt32(row.Orderable),
	}

	return inventory, nil
}
