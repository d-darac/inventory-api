package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func createItems(
	nItems int32,
	group *uuid.UUID,
	inventory *uuid.UUID,
	account uuid.UUID,
	q *database.Queries,
) ([]database.CreateItemRow, error) {
	if nItems == 0 {
		return nil, fmt.Errorf("nItems must be greater than 0")
	}

	rows := []database.CreateItemRow{}
	for i := range nItems {
		n := (i * 100) + 99
		params := database.CreateItemParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Description: sql.NullString{
				String: fmt.Sprintf("Test item description %d", i),
				Valid:  true,
			},
			Name: fmt.Sprintf("Test Item %d", i),
			PriceAmount: sql.NullInt32{
				Int32: n,
				Valid: true,
			},
			PriceCurrency: database.NullCurrency{
				Currency: database.CurrencyEUR,
				Valid:    true,
			},
			Type:      database.ItemTypePRODUCT,
			AccountID: account,
		}

		if group != nil {
			params.GroupID = uuid.NullUUID{
				UUID:  *group,
				Valid: true,
			}
		}
		if inventory != nil {
			params.InventoryID = uuid.NullUUID{
				UUID:  *inventory,
				Valid: true,
			}
		}

		row, err := q.CreateItem(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("couldn't create test item: %v", err)
		}

		rows = append(rows, row)
		time.Sleep(time.Millisecond * 100)
	}
	return rows, nil
}
