package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/d-darac/inventory-api/env"
	"github.com/d-darac/inventory-assets/auth"
	"github.com/d-darac/inventory-assets/database"
	"github.com/d-darac/inventory-assets/str"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

func main() {
	env := env.GetEnv()
	dbUrl := env.DB_URL
	platform := env.PLATFORM

	if strings.ToLower(platform) != "dev" {
		log.Fatalf("this should only be run in dev environment")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	password := "super_strong_password"
	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("couldn't hash password: %v", err)
	}

	usr, err := q.CreateUser(context.Background(), database.CreateUserParams{
		Email:          "email@test.com",
		HashedPassword: hashedPass,
		Name: sql.NullString{
			String: "Test User",
			Valid:  true,
		},
	})
	if err != nil {
		log.Fatalf("couldn't create test user: %v", err)
	}

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
		OwnerID: uuid.NullUUID{
			UUID:  usr.ID,
			Valid: true,
		},
	})
	if err != nil {
		log.Fatalf("couldn't create test account: %v", err)
	}

	apiKey := auth.GenApiKey(32)
	key := env.MASTER_KEY
	iv := env.IV

	encryptedApiKey, err := auth.EncryptApiKeySecret(apiKey, key, iv)
	if err != nil {
		log.Fatalf("couldn't encrypt the api key: %v", err)
	}

	_, err = q.CreateApiKey(context.Background(), database.CreateApiKeyParams{
		Name:           "Test Api Key",
		Secret:         encryptedApiKey,
		RedactedSecret: str.RedactString(apiKey, 4),
		AccountID:      acc.ID,
	})
	if err != nil {
		log.Fatalf("couldn't create test api key: %v", err)
	}

	grp, err := q.CreateGroup(context.Background(), database.CreateGroupParams{
		Name:      "Test Group",
		AccountID: acc.ID,
	})
	if err != nil {
		log.Fatalf("couldn't create test group: %v", err)
	}

	inv, err := q.CreateInventory(context.Background(), database.CreateInventoryParams{
		InStock: 10,
		Orderable: sql.NullInt32{
			Int32: 10,
			Valid: true,
		},
		AccountID: acc.ID,
	})
	if err != nil {
		log.Fatalf("couldn't create test inventory: %v", err)
	}

	itm, err := q.CreateItem(context.Background(), database.CreateItemParams{
		Name: "Test Item",
		PriceAmount: sql.NullInt32{
			Int32: int32(5),
			Valid: true,
		},
		PriceCurrency: database.NullCurrency{
			Currency: database.CurrencyEUR,
			Valid:    true,
		},
		Type:      database.ItemTypePRODUCT,
		AccountID: acc.ID,
		GroupID: uuid.NullUUID{
			UUID:  grp.ID,
			Valid: true,
		},
		InventoryID: uuid.NullUUID{
			UUID:  inv.ID,
			Valid: true,
		},
	})
	if err != nil {
		log.Fatalf("couldn't create test item: %v", err)
	}

	itmids, err := q.CreateItemIdentifier(context.Background(), database.CreateItemIdentifierParams{
		Ean: sql.NullString{
			String: "00000123",
			Valid:  true,
		},
		AccountID: acc.ID,
		ItemID:    itm.ID,
	})
	if err != nil {
		log.Fatalf("couldn't create test item identifiers: %v", err)
	}

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Api Key:             %s\n", apiKey)
	fmt.Printf("User ID:             %s\n", usr.ID)
	fmt.Printf("Account ID:          %s\n", acc.ID)
	fmt.Printf("Group ID:            %s\n", grp.ID)
	fmt.Printf("Inventory ID:        %s\n", inv.ID)
	fmt.Printf("Item ID:             %s\n", itm.ID)
	fmt.Printf("Item Identifiers ID: %s\n", itmids.ID)
	fmt.Println("----------------------------------------------------------------")

}
