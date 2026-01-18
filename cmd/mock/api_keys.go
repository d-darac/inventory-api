package main

import (
	"context"
	"fmt"
	"time"

	"github.com/d-darac/inventory-assets/auth"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"
)

func createApiKey(key, iv string, account uuid.UUID, q *database.Queries) (*string, error) {
	apiKey := auth.GenApiKey(32)
	encryptedApiKey, err := auth.EncryptApiKeySecret(apiKey, key, iv)
	if err != nil {
		return nil, fmt.Errorf("couldn't encrypt the api key: %v", err)
	}
	n := time.Now().UnixNano()
	_, err = q.CreateApiKey(context.Background(), database.CreateApiKeyParams{
		Name:           fmt.Sprintf("Test Api Key %d", n),
		Secret:         encryptedApiKey,
		RedactedSecret: str.RedactString(apiKey, 4),
		AccountID:      account,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create test api key: %v", err)
	}
	return &apiKey, nil
}
