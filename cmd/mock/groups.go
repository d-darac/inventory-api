package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func checkGroupExists(cfg cfg, group uuid.UUID) error {
	dest := &[]uint8{}
	row := cfg.db.QueryRow("SELECT id FROM groups WHERE id = $1", group)
	err := row.Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("group with id '%v' not found", group)
		} else {
			return err
		}
	}
	return nil
}

func createGroups(nGroups int, parentGroup *uuid.UUID, account uuid.UUID, q *database.Queries) ([]database.CreateGroupRow, error) {
	if nGroups == 0 {
		return nil, fmt.Errorf("nGroups must be greater than 0")
	}

	rows := []database.CreateGroupRow{}
	for range nGroups {
		n := time.Now().UnixNano()
		params := database.CreateGroupParams{
			Name:      fmt.Sprintf("Test Group %d", n),
			AccountID: account,
		}
		if parentGroup != nil {
			params.ParentID = uuid.NullUUID{
				UUID:  *parentGroup,
				Valid: true,
			}
		}
		row, err := q.CreateGroup(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("couldn't create test group: %v", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
