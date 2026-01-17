package groups

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/d-darac/inventory-assets/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func TestUnitNewListGroupsParams(t *testing.T) {
	lp := NewListGroupsParams()
	if lp.PaginationParams == nil {
		t.Fatal("PaginationParams is nil")
	}
	if lp.PaginationParams.Limit == nil {
		t.Fatal("Limit is nil")
	}
	if *lp.PaginationParams.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", *lp.PaginationParams.Limit)
	}
}

func TestUnitMapCreateGroupParams(t *testing.T) {
	acc := uuid.New()
	name := "test-group"
	desc := "description"
	parent := uuid.New()

	cp := CreateGroupParams{
		Description: &desc,
		Name:        name,
		ParentGroup: func() *string { s := parent.String(); return &s }(),
	}

	dbp := MapCreateGroupParams(Create{AccountId: acc, RequestParams: cp})
	if dbp.AccountID != acc {
		t.Fatalf("expected account id %v, got %v", acc, dbp.AccountID)
	}
	if dbp.Name != name {
		t.Fatalf("expected name %s, got %s", name, dbp.Name)
	}
	if (!dbp.Description.Valid) || dbp.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, dbp.Description.String)
	}
	if (!dbp.ParentID.Valid) || dbp.ParentID.UUID != parent {
		t.Fatalf("expected parent id %v, got %v", parent, dbp.ParentID.UUID)
	}
}

func TestUnitMapListGroupsParams(t *testing.T) {
	acc := uuid.New()
	lp := NewListGroupsParams()
	name := "test-group"
	desc := "description"
	lp.Name = &name
	lp.Description = &desc

	dbp := MapListGroupsParams(List{AccountId: acc, RequestParams: lp})
	if dbp.AccountID != acc {
		t.Fatalf("expected account id %v, got %v", acc, dbp.AccountID)
	}
	if (!dbp.Name.Valid) || dbp.Name.String != name {
		t.Fatalf("expected name %s, got %s", name, dbp.Name.String)
	}
	if (!dbp.Description.Valid) || dbp.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, dbp.Description.String)
	}
}

func TestUnitMapUpdateGroupParams(t *testing.T) {
	id := uuid.New()
	acc := uuid.New()
	name := "test-group"
	up := UpdateGroupParams{
		Name: &name,
	}

	dbp := MapUpdateGroupParams(Update{AccountId: acc, GroupId: id, RequestParams: up})
	if dbp.ID != id {
		t.Fatalf("expected id %v, got %v", id, dbp.ID)
	}
	if dbp.AccountID != acc {
		t.Fatalf("expected account id %v, got %v", acc, dbp.AccountID)
	}
	if (!dbp.Name.Valid) || dbp.Name.String != name {
		t.Fatalf("expected name %s, got %s", name, dbp.Name.String)
	}
}

func TestIntegrationCreate(t *testing.T) {
	godotenv.Load("../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("couldn't create test account: %v", err)
	}
	defer db.Exec("DELETE FROM accounts WHERE id = $1;", acc.ID)

	s := NewGroupsService(q)
	parentName := "test-group-parent"
	parentDesc := "parent-description"
	childName := "test-group-child"
	childDesc := "child-description"

	parentCp := CreateGroupParams{
		Description: &parentDesc,
		Name:        parentName,
	}
	parentGroup, err := s.Create(Create{AccountId: acc.ID, RequestParams: parentCp})
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}
	if parentGroup.Name != parentName {
		t.Fatalf("expected name %v, got %v", parentName, parentGroup.Name)
	}
	if (!parentGroup.Description.Valid) || parentGroup.Description.String != parentDesc {
		t.Fatalf("expected description %s, got %s", parentDesc, parentGroup.Description.String)
	}

	parentGroupRow, err := q.GetGroup(context.Background(), database.GetGroupParams{ID: *parentGroup.ID, AccountID: acc.ID})
	if err != nil {
		t.Fatalf("error retrieving group: %v", err)
	}
	if parentGroupRow.Name != parentName {
		t.Fatalf("expected name %v, got %v", parentName, parentGroupRow.Name)
	}
	if (!parentGroupRow.Description.Valid) || parentGroupRow.Description.String != parentDesc {
		t.Fatalf("expected description %s, got %s", parentDesc, parentGroupRow.Description.String)
	}

	childCp := CreateGroupParams{
		Description: &childDesc,
		Name:        childName,
		ParentGroup: func() *string { s := parentGroupRow.ID.String(); return &s }(),
	}
	childGroup, err := s.Create(Create{AccountId: acc.ID, RequestParams: childCp})
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}
	if childGroup.Name != childName {
		t.Fatalf("expected name %v, got %v", childName, childGroup.Name)
	}
	if (!childGroup.Description.Valid) || childGroup.Description.String != childDesc {
		t.Fatalf("expected description %s, got %s", childDesc, childGroup.Description.String)
	}
	if (!childGroup.ParentGroup.ID.Valid) || childGroup.ParentGroup.ID.UUID != *parentGroup.ID {
		t.Fatalf("expected parent group %v, got %v", parentGroup.ID, childGroup.ParentGroup.ID.UUID)
	}

	childGroupRow, err := q.GetGroup(context.Background(), database.GetGroupParams{ID: *childGroup.ID, AccountID: acc.ID})
	if err != nil {
		t.Fatalf("error retrieving group: %v", err)
	}
	if childGroupRow.Name != childName {
		t.Fatalf("expected name %v, got %v", childName, childGroupRow.Name)
	}
	if (!childGroupRow.Description.Valid) || childGroupRow.Description.String != childDesc {
		t.Fatalf("expected description %s, got %s", childDesc, childGroupRow.Description.String)
	}
	if (!childGroupRow.ParentGroup.Valid) || childGroupRow.ParentGroup.UUID != parentGroupRow.ID {
		t.Fatalf("expected parent group %v, got %s", parentGroupRow.ID, childGroupRow.ParentGroup.UUID)
	}
}

func TestIntegrationDelete(t *testing.T) {
	godotenv.Load("../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("couldn't create test account: %v", err)
	}
	defer db.Exec("DELETE FROM accounts WHERE id = $1;", acc.ID)

	s := NewGroupsService(q)
	name := "test-group"
	desc := "description"

	grp, err := q.CreateGroup(context.Background(), database.CreateGroupParams{
		Description: sql.NullString{
			String: desc,
			Valid:  true,
		},
		Name:      name,
		AccountID: acc.ID,
	})
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}

	err = s.Delete(Delete{AccountId: acc.ID, GroupId: grp.ID})
	if err != nil {
		t.Fatalf("error deleting group: %v", err)
	}

	_, err = q.GetGroup(context.Background(), database.GetGroupParams{ID: grp.ID, AccountID: acc.ID})
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
}

func TestIntegrationList(t *testing.T) {
	godotenv.Load("../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("couldn't create test account: %v", err)
	}
	defer db.Exec("DELETE FROM accounts WHERE id = $1;", acc.ID)

	s := NewGroupsService(q)
	name := "test-group"
	desc := "description"
	rows := make([]database.CreateGroupRow, 0, 10)

	for i := range 10 {
		grp, err := q.CreateGroup(context.Background(), database.CreateGroupParams{
			Description: sql.NullString{
				String: fmt.Sprintf("%s %d", desc, i),
				Valid:  true,
			},
			Name:      fmt.Sprintf("%s %d", name, i),
			AccountID: acc.ID,
		})
		if err != nil {
			t.Fatalf("couldn't create test group: %v", err)
		}
		rows = append(rows, grp)
	}

	groups, hasMore, err := s.List(List{AccountId: acc.ID, RequestParams: NewListGroupsParams()})
	if err != nil {
		t.Fatalf("error listing groups: %v", err)
	}

	if hasMore {
		t.Fatalf("expected hasMore %v, got %v", false, hasMore)
	}

	if len(groups) != 10 {
		t.Fatalf("expected len %v got %v", 10, len(groups))
	}

	i := 9

	for _, group := range groups {
		if group.Name != rows[i].Name {
			t.Fatalf("expected name %v, got %v", rows[i].Name, group.Name)
		}
		if (!group.Description.Valid) || group.Description.String != rows[i].Description.String {
			t.Fatalf("expected description %s, got %s", rows[i].Description.String, group.Description.String)
		}
		i--
	}
}

func TestIntegrationRetrieve(t *testing.T) {
	godotenv.Load("../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("couldn't create test account: %v", err)
	}
	defer db.Exec("DELETE FROM accounts WHERE id = $1;", acc.ID)

	s := NewGroupsService(q)
	name := "test-group"
	desc := "description"

	grp, err := q.CreateGroup(context.Background(), database.CreateGroupParams{
		Description: sql.NullString{
			String: desc,
			Valid:  true,
		},
		Name:      name,
		AccountID: acc.ID,
	})
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}

	group, err := s.Get(Get{
		AccountId:     acc.ID,
		GroupId:       grp.ID,
		RequestParams: RetrieveGroupParams{},
		OmitBase:      false,
	})
	if err != nil {
		t.Fatalf("couldn't get test group: %v", err)
	}

	if group.Name != name {
		t.Fatalf("expected name %v, got %v", name, group.Name)
	}
	if (!group.Description.Valid) || group.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, group.Description.String)
	}

	row, err := q.GetGroup(context.Background(), database.GetGroupParams{ID: grp.ID, AccountID: acc.ID})
	if err != nil {
		t.Fatalf("error retrieving group: %v", err)
	}

	if row.Name != name {
		t.Fatalf("expected name %v, got %v", name, row.Name)
	}
	if (!row.Description.Valid) || row.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, row.Description.String)
	}
}

func TestIntegrationUpdate(t *testing.T) {
	godotenv.Load("../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("couldn't connect to database: %v", err)
	}

	q := database.New(db)

	acc, err := q.CreateAccount(context.Background(), database.CreateAccountParams{
		Country: database.CountryIE,
		Nickname: sql.NullString{
			String: "Test Account",
			Valid:  true,
		},
	})
	if err != nil {
		t.Fatalf("couldn't create test account: %v", err)
	}

	s := NewGroupsService(q)
	name := "test-group"
	desc := "description"
	updatedName := "test-group-updated"
	updatedDesc := "description-updated"

	grp, err := q.CreateGroup(context.Background(), database.CreateGroupParams{
		Description: sql.NullString{
			String: desc,
			Valid:  true,
		},
		Name:      name,
		AccountID: acc.ID,
	})
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}
	defer db.Exec("DELETE FROM accounts WHERE id = $1;", acc.ID)

	up := UpdateGroupParams{
		Name:        &updatedName,
		Description: &updatedDesc,
	}
	group, err := s.Update(Update{AccountId: acc.ID, GroupId: grp.ID, RequestParams: up})
	if err != nil {
		t.Fatalf("couldn't update test group: %v", err)
	}

	if group.Name != updatedName {
		t.Fatalf("expected name %v, got %v", updatedName, group.Name)
	}
	if (!group.Description.Valid) || group.Description.String != updatedDesc {
		t.Fatalf("expected description %s, got %s", updatedDesc, group.Description.String)
	}

	row, err := q.GetGroup(context.Background(), database.GetGroupParams{ID: grp.ID, AccountID: acc.ID})
	if err != nil {
		t.Fatalf("error retrieving group: %v", err)
	}

	if row.Name != updatedName {
		t.Fatalf("expected name %v, got %v", updatedName, row.Name)
	}
	if (!row.Description.Valid) || row.Description.String != updatedDesc {
		t.Fatalf("expected description %s, got %s", updatedDesc, row.Description.String)
	}
}
