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

	cp := &CreateGroupParams{
		Description: &desc,
		Name:        name,
		ParentGroup: &parent,
	}

	dbp := MapCreateGroupParams(acc, cp)
	if dbp.AccountID != acc {
		t.Fatalf("expected account id %v, got %v", acc, dbp.AccountID)
	}
	if dbp.Name != name {
		t.Fatalf("expected name %s, got %s", name, dbp.Name)
	}
	if (!dbp.Description.Valid) || dbp.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, dbp.Description.String)
	}
}

func TestUnitMapListGroupsParams(t *testing.T) {
	acc := uuid.New()
	lp := NewListGroupsParams()
	name := "test-group"
	desc := "description"
	lp.Name = &name
	lp.Description = &desc

	dbp := MapListGroupsParams(acc, lp)
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
	up := &UpdateGroupParams{
		Name: &name,
	}

	dbp := MapUpdateGroupParams(id, acc, up)
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
	godotenv.Load("../../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()
	defer db.Exec("TRUNCATE accounts CASCADE;")

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

	cp := &CreateGroupParams{
		Description: &desc,
		Name:        name,
	}

	group, err := s.Create(acc.ID, cp)
	if err != nil {
		t.Fatalf("couldn't create test group: %v", err)
	}

	if group.Name != name {
		t.Fatalf("expected name %v, got %v", name, group.Name)
	}
	if (!group.Description.Valid) || group.Description.String != desc {
		t.Fatalf("expected description %s, got %s", desc, group.Description.String)
	}

	row, err := q.GetGroup(context.Background(), database.GetGroupParams{ID: *group.ID, AccountID: acc.ID})
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

func TestIntegrationDelete(t *testing.T) {
	godotenv.Load("../../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()
	defer db.Exec("TRUNCATE accounts CASCADE;")

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

	err = s.Delete(grp.ID, acc.ID)
	if err != nil {
		t.Fatalf("error deleting group: %v", err)
	}

	_, err = q.GetGroup(context.Background(), database.GetGroupParams{ID: grp.ID, AccountID: acc.ID})
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}
}

func TestIntegrationList(t *testing.T) {
	godotenv.Load("../../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()
	defer db.Exec("TRUNCATE accounts CASCADE;")

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

	groups, hasMore, err := s.List(acc.ID, NewListGroupsParams())
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
	godotenv.Load("../../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()
	defer db.Exec("TRUNCATE accounts CASCADE;")

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

	group, err := s.Get(grp.ID, acc.ID, &RetrieveGroupParams{}, false)
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
	godotenv.Load("../../../.env")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatalf("couldn't open database: %v", err)
	}
	defer db.Close()
	defer db.Exec("TRUNCATE accounts CASCADE;")

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

	group, err := s.Update(grp.ID, acc.ID, &UpdateGroupParams{
		Name:        &updatedName,
		Description: &updatedDesc,
	})
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
