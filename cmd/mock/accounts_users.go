package main

import (
	"context"
	"fmt"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func createAccountUserReference(account, user uuid.UUID, q *database.Queries) error {
	err := q.CreateAccountUserReference(context.Background(), database.CreateAccountUserReferenceParams{
		AccountID: account,
		UserID:    user,
	})
	if err != nil {
		return fmt.Errorf("couldn't create test account-user reference: %v", err)
	}
	return nil
}
