package main

import (
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

	user, err := createUser(qtx)
	if err != nil {
		return err
	}

	account, err := createAccount(user.ID, qtx)
	if err != nil {
		return err
	}

	err = createAccountUserReference(account.ID, user.ID, qtx)
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
	fmt.Printf("User ID:             %s\n", user.ID)
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
		return fmt.Errorf("usage: cli groups <n_groups> <account_id> [...args]")
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
		return fmt.Errorf("usage: cli inventories <n_inventories> <account_id>")
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
		return fmt.Errorf("usage: cli items <n_items> <account_id> [...args]")
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
