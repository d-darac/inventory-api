package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
)

func handleAll(cfg cfg) error {
	tx, err := cfg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := cfg.q.WithTx(tx)

	account, err := createAccount(nil, qtx)
	if err != nil {
		return err
	}

	key := cfg.env.MASTER_KEY
	iv := cfg.env.IV

	apiKey, err := createApiKey(key, iv, account.ID, qtx)
	if err != nil {
		return err
	}

	parentGroups, err := createGroups(3, nil, account.ID, qtx)
	if err != nil {
		return err
	}

	childGroups := []database.CreateGroupRow{}
	for _, g := range parentGroups {
		childGroups, err = createGroups(5, &g.ID, account.ID, qtx)
		if err != nil {
			return err
		}
	}

	inventories, err := createInventories(15, account.ID, qtx)
	if err != nil {
		return err
	}

	items, err := createItems(15, &childGroups[0].ID, &inventories[0].ID, account.ID, qtx)
	if err != nil {
		return err
	}

	itemIds := []uuid.UUID{}
	for _, i := range items {
		itemIds = append(itemIds, i.ID)
	}

	itemIdentifiers, err := createItemIdentifiers(itemIds, account.ID, qtx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Api Key:             %s\n", *apiKey)
	fmt.Printf("Account ID:          %s\n", account.ID)
	fmt.Printf("Group ID:            %s\n", childGroups[0].ID)
	fmt.Printf("Inventory ID:        %s\n", inventories[0].ID)
	fmt.Printf("Item ID:             %s\n", items[0].ID)
	fmt.Printf("Item Identifiers ID: %s\n", itemIdentifiers[0].ID)
	fmt.Println("----------------------------------------------------------------")

	return nil
}

func handleGroups(cfg cfg, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: cli groups COUNT ACCOUNT [ARG...]")
	}

	nGroups, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	if nGroups > 100 {
		return fmt.Errorf("can only create up to 100 groups")
	}

	account, err := uuid.Parse(args[1])
	if err != nil {
		return err
	}

	optional := args[2:]

	var parentGroup *uuid.UUID

	for _, arg := range optional {
		key, value, err := argToKeyVal(arg)
		if err != nil {
			return err
		}

		if key == "group" {
			id, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			parentGroup = &id
		} else {
			return fmt.Errorf("unknown argument: %s", key)
		}
	}

	if err = checkAccountExists(cfg, account); err != nil {
		return err
	}

	if parentGroup != nil {
		if err = checkGroupExists(cfg, *parentGroup); err != nil {
			return err
		}
	}

	groups, err := createGroups(nGroups, parentGroup, account, cfg.q)
	if err != nil {
		return err
	}

	ids := make([]uuid.UUID, 0, len(groups))
	for _, group := range groups {
		ids = append(ids, group.ID)
	}

	printIDs(ids)

	return nil
}

func handleInventories(cfg cfg, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: cli inventories COUNT ACCOUNT")
	}

	nInventories, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if nInventories > 100 {
		return fmt.Errorf("can only create up to 100 inventories")
	}

	account, err := uuid.Parse(args[1])
	if err != nil {
		return err
	}

	if err = checkAccountExists(cfg, account); err != nil {
		return err
	}

	inventories, err := createInventories(int32(nInventories), account, cfg.q)
	if err != nil {
		return err
	}

	ids := make([]uuid.UUID, 0, len(inventories))
	for _, inventory := range inventories {
		ids = append(ids, inventory.ID)
	}

	printIDs(ids)

	return nil
}

func handleItems(cfg cfg, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: cli items COUNT ACCOUNT [ARG...]")
	}

	nItems, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if nItems > 100 {
		return fmt.Errorf("can only create up to 100 items")
	}

	account, err := uuid.Parse(args[1])
	if err != nil {
		return err
	}

	optional := args[2:]

	var group *uuid.UUID
	var inventory *uuid.UUID

	for _, arg := range optional {
		key, value, err := argToKeyVal(arg)
		if err != nil {
			return err
		}
		switch key {
		case "group":
			id, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			group = &id
		case "inventory":
			id, err := uuid.Parse(value)
			if err != nil {
				return err
			}
			inventory = &id
		default:
			return fmt.Errorf("unknown argument: %s", key)
		}
	}

	if err = checkAccountExists(cfg, account); err != nil {
		return err
	}

	if group != nil {
		if err = checkGroupExists(cfg, *group); err != nil {
			return err
		}
	}

	if inventory != nil {
		if err = checkInventoryExists(cfg, *inventory); err != nil {
			return err
		}
	}

	items, err := createItems(int32(nItems), group, inventory, account, cfg.q)
	if err != nil {
		return err
	}

	ids := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	printIDs(ids)

	return nil
}

func handleWipeAll(cfg cfg) error {
	tx, err := cfg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("DELETE FROM accounts;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt2, err := tx.Prepare("DELETE FROM users;")
	if err != nil {
		return err
	}
	defer stmt2.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	if _, err := stmt2.Exec(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Database wiped!")
	return nil
}

func handleWipe(cfg cfg, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Usage: cli wipe ACCOUNT...")
	}
	tx, err := cfg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	accStmt := make(map[uuid.UUID]*sql.Stmt)
	for _, arg := range args {
		account, err := uuid.Parse(arg)
		if err != nil {
			return err
		}
		if err := checkAccountExists(cfg, account); err != nil {
			return err
		}
		stmt, err := tx.Prepare("DELETE FROM accounts WHERE id = $1;")
		if err != nil {
			return err
		}
		accStmt[account] = stmt
	}

	for account, stmt := range accStmt {
		if _, err := stmt.Exec(account); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Database wiped!")
	return nil
}

func handleHelp(args ...string) error {
	if len(args) > 0 {
		topic := args[0]
		switch topic {
		case "all":
			fmt.Println("Usage: cli all")
			fmt.Println()
			fmt.Println("Create all resources")
			return nil
		case "groups":
			fmt.Println("Usage: cli groups COUNT ACCOUNT [ARG...]")
			fmt.Println()
			fmt.Println("Create COUNT groups for a specified ACCOUNT")
			fmt.Println()
			fmt.Println("Optional Arguments:")
			fmt.Println("  group=GROUP_ID: use to assign groups to a specific parent group")
			return nil
		case "inventories":
			fmt.Println("Usage: cli inventories COUNT ACCOUNT")
			fmt.Println()
			fmt.Println("Create COUNT inventories for a specified ACCOUNT")
			return nil
		case "items":
			fmt.Println("Usage: cli items COUNT ACCOUNT [ARG...]")
			fmt.Println()
			fmt.Println("Create COUNT items for a specified ACCOUNT")
			fmt.Println()
			fmt.Println("Optional Arguments:")
			fmt.Println("  group=GROUP_ID:         use to assign items to a specific group")
			fmt.Println("  inventory=INVENTORY_ID: use to assign items to a specific inventory")
			return nil
		case "wipeall":
			fmt.Println("Usage: cli wipeall")
			fmt.Println()
			fmt.Println("Wipe the database of all records")
			return nil
		case "wipe":
			fmt.Println("Usage: cli wipe ACCOUNT...")
			fmt.Println()
			fmt.Println("Wipe the database of all records for a specified list of accounts")
			return nil
		case "help":
			fmt.Println("Usage: cli help [COMMAND]")
			fmt.Println()
			fmt.Println("Help about the command")
			return nil
		default:
			return fmt.Errorf("unknown help topic: %s", topic)
		}
	}

	fmt.Println("Usage: cli COMMAND [ARG...]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  all:         create all resources")
	fmt.Println("  groups:      create groups for a specific account")
	fmt.Println("  inventories: create inventories for a specific account")
	fmt.Println("  items:       create items for a specific account")
	fmt.Println("  wipeall:     delete all records in the database")
	fmt.Println("  wipe:        delete all records for a specific account")

	return nil
}
