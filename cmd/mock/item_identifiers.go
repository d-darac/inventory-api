package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func createItemIdentifiers(items []uuid.UUID, account uuid.UUID, q *database.Queries) ([]database.CreateItemIdentifierRow, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items slice must have at least one element")
	}

	rows := []database.CreateItemIdentifierRow{}
	for _, item := range items {
		n := time.Now().UnixNano()
		row, err := q.CreateItemIdentifier(context.Background(), database.CreateItemIdentifierParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Sku: sql.NullString{
				String: fmt.Sprintf("%d", n),
				Valid:  true,
			},
			AccountID: account,
			ItemID:    item,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't create test item identifiers: %v", err)
		}
		rows = append(rows, row)
		time.Sleep(time.Millisecond * 100)
	}
	return rows, nil
}
