package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func checkAccountExists(cfg cfg, account uuid.UUID) error {
	dest := &[]uint8{}
	row := cfg.db.QueryRow("SELECT id FROM accounts WHERE id = $1", account)
	err := row.Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("account with id '%v' not found", account)
		} else {
			return err
		}
	}
	return nil
}

func createAccount(user *uuid.UUID, q *database.Queries) (*database.CreateAccountRow, error) {
	n := time.Now().UnixNano()
	params := database.CreateAccountParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Country:   database.CountryIE,
		Nickname: sql.NullString{
			String: fmt.Sprintf("Test Account %d", n),
			Valid:  true,
		},
	}
	if user != nil {
		params.OwnerID = uuid.NullUUID{
			UUID:  *user,
			Valid: true,
		}
	}
	account, err := q.CreateAccount(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("couldn't create test account: %v", err)
	}
	return &account, nil
}
