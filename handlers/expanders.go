package handlers

import (
	"slices"

	"github.com/d-darac/inventory-api/internal/groups"
	"github.com/d-darac/inventory-api/internal/inventories"
	itemidentifiers "github.com/d-darac/inventory-api/internal/item_identifiers"
	"github.com/d-darac/inventory-api/internal/items"
	"github.com/d-darac/inventory-assets/api"
	"github.com/google/uuid"
)

func (h *GroupsHandler) ExpandFieldsList(fields *[]string, groups []*groups.Group, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "parent_group") {
		err := h.expandGroups(groups, accountId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *GroupsHandler) ExpandFields(fields *[]string, group *groups.Group, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "parent_group") {
		getParams := groups.Get{
			AccountId:     accountId,
			GroupId:       group.ParentGroup.ID.UUID,
			RequestParams: groups.RetrieveGroupParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&group.ParentGroup, h.Groups.Get, getParams); err != nil {
			return err
		}
	}
	return nil
}

func (h *GroupsHandler) expandGroups(grps []*groups.Group, accountId uuid.UUID) error {
	ids := make([]uuid.UUID, 0, len(grps))

	withNonNillParentGroup := make([]*groups.Group, 0)
	for _, group := range grps {
		if group.ParentGroup.ID.Valid {
			ids = append(ids, group.ParentGroup.ID.UUID)
			withNonNillParentGroup = append(withNonNillParentGroup, group)
		}
	}

	parentGroups, err := h.Groups.ListByIds(groups.ListByIds{
		AccountId: accountId,
		RequestParams: groups.ListGroupsByIdsParams{
			Ids: ids,
		},
	})
	if err != nil {
		return err
	}

	idParentGroupMap := make(map[uuid.UUID]*groups.Group, 0)
	for _, parentGroup := range parentGroups {
		idParentGroupMap[*parentGroup.ID] = parentGroup
	}

	for _, group := range withNonNillParentGroup {
		if parentGroup, ok := idParentGroupMap[group.ParentGroup.ID.UUID]; ok {
			group.ParentGroup.Resource = parentGroup
		}
	}

	return nil
}

func (h *ItemIdentifiersHandler) ExpandFieldsList(fields *[]string, itemIdentifiers []*itemidentifiers.ItemIdentifiers, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "item") {
		err := h.expandItems(itemIdentifiers, accountId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *ItemIdentifiersHandler) ExpandFields(fields *[]string, itemIdentifiers *itemidentifiers.ItemIdentifiers, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "item") {
		getParams := items.Get{
			AccountId:     accountId,
			ItemId:        itemIdentifiers.Item.ID.UUID,
			RequestParams: items.RetrieveItemParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&itemIdentifiers.Item, h.Items.Get, getParams); err != nil {
			return err
		}
	}
	return nil
}

func (h *ItemIdentifiersHandler) expandItems(idtfs []*itemidentifiers.ItemIdentifiers, accountId uuid.UUID) error {
	itemsIds := make([]uuid.UUID, 0, len(idtfs))

	withNonNillItem := make([]*itemidentifiers.ItemIdentifiers, 0)
	for _, itemIdentifiers := range idtfs {
		if itemIdentifiers.Item.ID.Valid {
			itemsIds = append(itemsIds, itemIdentifiers.Item.ID.UUID)
			withNonNillItem = append(withNonNillItem, itemIdentifiers)
		}
	}

	itms, err := h.Items.ListByIds(items.ListByIds{
		AccountId: accountId,
		RequestParams: items.ListItemsByIdsParams{
			Ids: itemsIds,
		},
	})
	if err != nil {
		return err
	}

	idItemMap := make(map[uuid.UUID]*items.Item, 0)
	for _, item := range itms {
		idItemMap[*item.ID] = item
	}

	for _, itemIdentifiers := range withNonNillItem {
		if item, ok := idItemMap[itemIdentifiers.Item.ID.UUID]; ok {
			itemIdentifiers.Item.Resource = item
		}
	}

	return nil
}

func (h *ItemsHandler) ExpandFieldsList(fields *[]string, items []*items.Item, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "group") {
		err := h.expandGroups(items, accountId)
		if err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "inventory") {
		err := h.expandInventories(items, accountId)
		if err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "identifiers") {
		err := h.expandIdentifiers(items, accountId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *ItemsHandler) ExpandFields(fields *[]string, item *items.Item, accountId uuid.UUID) error {
	if fields != nil && slices.Contains(*fields, "group") {
		getParams := groups.Get{
			AccountId:     accountId,
			GroupId:       item.Group.ID.UUID,
			RequestParams: groups.RetrieveGroupParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&item.Group, h.Groups.Get, getParams); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "inventory") {
		getParams := inventories.Get{
			AccountId:     accountId,
			InventoryId:   item.Inventory.ID.UUID,
			RequestParams: inventories.RetrieveInventoryParams{},
			OmitBase:      true,
		}
		if _, err := api.ExpandField(&item.Inventory, h.Inventories.Get, getParams); err != nil {
			return err
		}
	}
	if fields != nil && slices.Contains(*fields, "identifiers") {
		getParams := itemidentifiers.Get{
			AccountId:         accountId,
			ItemIdentifiersId: item.Identifiers.ID.UUID,
			RequestParams:     itemidentifiers.RetrieveItemIdentifiersParams{},
			OmitBase:          true,
		}
		if _, err := api.ExpandField(&item.Identifiers, h.ItemIdentifiers.Get, getParams); err != nil {
			return err
		}
	}
	return nil
}

func (h *ItemsHandler) expandGroups(itms []*items.Item, accountId uuid.UUID) error {
	groupsIds := make([]uuid.UUID, 0, len(itms))

	withNonNillGroup := make([]*items.Item, 0)
	for _, item := range itms {
		if item.Group.ID.Valid {
			groupsIds = append(groupsIds, item.Group.ID.UUID)
			withNonNillGroup = append(withNonNillGroup, item)
		}
	}

	grps, err := h.Groups.ListByIds(groups.ListByIds{
		AccountId: accountId,
		RequestParams: groups.ListGroupsByIdsParams{
			Ids: groupsIds,
		},
	})
	if err != nil {
		return err
	}

	idGroupMap := make(map[uuid.UUID]*groups.Group, 0)
	for _, group := range grps {
		idGroupMap[*group.ID] = group
	}

	for _, item := range withNonNillGroup {
		if group, ok := idGroupMap[item.Group.ID.UUID]; ok {
			item.Group.Resource = group
		}
	}

	return nil
}

func (h *ItemsHandler) expandInventories(itms []*items.Item, accountId uuid.UUID) error {
	inventoriesIds := make([]uuid.UUID, 0, len(itms))

	withNonNillInventory := make([]*items.Item, 0)
	for _, item := range itms {
		if item.Inventory.ID.Valid {
			inventoriesIds = append(inventoriesIds, item.Inventory.ID.UUID)
			withNonNillInventory = append(withNonNillInventory, item)
		}
	}

	invs, err := h.Inventories.ListByIds(inventories.ListByIds{
		AccountId: accountId,
		RequestParams: inventories.ListInventoriesByIdsParams{
			Ids: inventoriesIds,
		},
	})
	if err != nil {
		return err
	}

	idInventoryMap := make(map[uuid.UUID]*inventories.Inventory, 0)
	for _, inventory := range invs {
		idInventoryMap[*inventory.ID] = inventory
	}

	for _, item := range withNonNillInventory {
		if inventory, ok := idInventoryMap[item.Inventory.ID.UUID]; ok {
			item.Inventory.Resource = inventory
		}
	}

	return nil
}

func (h *ItemsHandler) expandIdentifiers(itms []*items.Item, accountId uuid.UUID) error {
	identifiersIds := make([]uuid.UUID, 0, len(itms))

	withNonNillIdentifiers := make([]*items.Item, 0)
	for _, item := range itms {
		if item.Identifiers.ID.Valid {
			identifiersIds = append(identifiersIds, item.Identifiers.ID.UUID)
			withNonNillIdentifiers = append(withNonNillIdentifiers, item)
		}
	}

	invs, err := h.ItemIdentifiers.ListByIds(itemidentifiers.ListByIds{
		AccountId: accountId,
		RequestParams: itemidentifiers.ListItemIdentifiersByIdsParams{
			Ids: identifiersIds,
		},
	})
	if err != nil {
		return err
	}

	idIdentifiersMap := make(map[uuid.UUID]*itemidentifiers.ItemIdentifiers, 0)
	for _, identifiers := range invs {
		idIdentifiersMap[*identifiers.ID] = identifiers
	}

	for _, item := range withNonNillIdentifiers {
		if itemIdentifiers, ok := idIdentifiersMap[item.Identifiers.ID.UUID]; ok {
			item.Identifiers.Resource = itemIdentifiers
		}
	}

	return nil
}
