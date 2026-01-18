package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func checkInventoryExists(cfg cfg, inventory uuid.UUID) error {
	dest := &[]uint8{}
	row := cfg.db.QueryRow("SELECT id FROM inventories WHERE id = $1", inventory)
	err := row.Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("inventory with id '%v' not found\n", inventory)
		} else {
			return err
		}
	}
	return nil
}

func createInventories(nInventories int32, account uuid.UUID, q *database.Queries) ([]database.CreateInventoryRow, error) {
	if nInventories == 0 {
		return nil, fmt.Errorf("nInventories must be greater than 0")
	}
	rows := []database.CreateInventoryRow{}
	for i := range nInventories {
		n := (i * (i + 1)) + i
		row, err := q.CreateInventory(context.Background(), database.CreateInventoryParams{
			InStock: n,
			Orderable: sql.NullInt32{
				Int32: n,
				Valid: true,
			},
			AccountID: account,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't create test inventory: %v", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
