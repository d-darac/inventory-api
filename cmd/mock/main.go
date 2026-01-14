package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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

	user, err := user(q)
	if err != nil {
		log.Fatalln(err)
	}

	account, err := account(*user, q)
	if err != nil {
		log.Fatalln(err)
	}

	key := env.MASTER_KEY
	iv := env.IV

	apiKey, err := apiKey(key, iv, *account, q)
	if err != nil {
		log.Fatalln(err)
	}

	parentGroups, err := groups(3, nil, *account, q)
	if err != nil {
		log.Fatalln(err)
	}

	childGroups := []database.CreateGroupRow{}
	for _, g := range parentGroups {
		childGroups, err = groups(5, &g, *account, q)
		if err != nil {
			log.Fatalln(err)
		}
	}

	inventories, err := inventories(15, *account, q)
	if err != nil {
		log.Fatalln(err)
	}

	items, err := items(15, &childGroups[0], &inventories[0], *account, q)
	if err != nil {
		log.Fatalln(err)
	}

	itemIdentifiers, err := itemIdentifiers(items, *account, q)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Api Key:             %s\n", *apiKey)
	fmt.Printf("User ID:             %s\n", user.ID)
	fmt.Printf("Account ID:          %s\n", account.ID)
	fmt.Printf("Group ID:            %s\n", childGroups[0].ID)
	fmt.Printf("Inventory ID:        %s\n", inventories[0].ID)
	fmt.Printf("Item ID:             %s\n", items[0].ID)
	fmt.Printf("Item Identifiers ID: %s\n", itemIdentifiers[0].ID)
	fmt.Println("----------------------------------------------------------------")
}

func user(q *database.Queries) (*database.CreateUserRow, error) {
	n := time.Now().UnixNano()
	password := fmt.Sprintf("super_strong_password_%d", n)
	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("couldn't hash password: %v", err)
	}
	user, err := q.CreateUser(context.Background(), database.CreateUserParams{
		Email:          fmt.Sprintf("email_%d@test.com", n),
		HashedPassword: hashedPass,
		Name: sql.NullString{
			String: fmt.Sprintf("Test User %d", n),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create test user: %v", err)
	}
	return &user, nil
}

func account(user database.CreateUserRow, q *database.Queries) (*database.CreateAccountRow, error) {
	n := time.Now().UnixNano()
	account, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: fmt.Sprintf("Test Account %d", n),
			Valid:  true,
		},
		OwnerID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create test account: %v", err)
	}
	return &account, nil
}

func apiKey(key, iv string, account database.CreateAccountRow, q *database.Queries) (*string, error) {
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
		AccountID:      account.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create test api key: %v", err)
	}
	return &apiKey, nil
}

func groups(nGroups int, parentGroup *database.CreateGroupRow, account database.CreateAccountRow, q *database.Queries) ([]database.CreateGroupRow, error) {
	if nGroups == 0 {
		return nil, fmt.Errorf("nGroups must be greater than 0")
	}

	rows := []database.CreateGroupRow{}
	for range nGroups {
		n := time.Now().UnixNano()
		params := database.CreateGroupParams{
			Name:      fmt.Sprintf("Test Group %d", n),
			AccountID: account.ID,
		}
		if parentGroup != nil {
			params.ParentID = uuid.NullUUID{
				UUID:  parentGroup.ID,
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

func inventories(nInventories int32, account database.CreateAccountRow, q *database.Queries) ([]database.CreateInventoryRow, error) {
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
			AccountID: account.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't create test inventory: %v", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func items(
	nItems int32,
	group *database.CreateGroupRow,
	inventory *database.CreateInventoryRow,
	account database.CreateAccountRow,
	q *database.Queries,
) ([]database.CreateItemRow, error) {
	if nItems == 0 {
		return nil, fmt.Errorf("nItems must be greater than 0")
	}
	rows := []database.CreateItemRow{}
	for i := range nItems {
		n := (i * 100) + 99
		params := database.CreateItemParams{
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
			AccountID: account.ID,
		}

		if group != nil {
			params.GroupID = uuid.NullUUID{
				UUID:  group.ID,
				Valid: true,
			}
		}
		if inventory != nil {
			params.InventoryID = uuid.NullUUID{
				UUID:  inventory.ID,
				Valid: true,
			}
		}

		row, err := q.CreateItem(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("couldn't create test item: %v", err)
		}

		rows = append(rows, row)
	}
	return rows, nil
}

func itemIdentifiers(items []database.CreateItemRow, account database.CreateAccountRow, q *database.Queries) ([]database.CreateItemIdentifierRow, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items slice must have at least one element")
	}

	rows := []database.CreateItemIdentifierRow{}
	for _, item := range items {
		n := time.Now().UnixNano()
		row, err := q.CreateItemIdentifier(context.Background(), database.CreateItemIdentifierParams{
			Sku: sql.NullString{
				String: fmt.Sprintf("%d", n),
				Valid:  true,
			},
			AccountID: account.ID,
			ItemID:    item.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't create test item identifiers: %v", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}
