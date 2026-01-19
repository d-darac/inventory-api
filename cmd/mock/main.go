package main

import (
	"database/sql"
	"fmt"
	"os"

	"strings"

	"github.com/d-darac/inventory-api/env"

	"github.com/d-darac/inventory-assets/database"

	_ "github.com/lib/pq"
)

type cfg struct {
	db  *sql.DB
	env env.Env
	q   *database.Queries
}

func main() {
	if len(os.Args) < 2 {
		exitWithErr(fmt.Errorf("usage: cli <command> [...args]"))
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	db, env, err := setup()
	if err != nil {
		exitWithErr(err)
	}
	defer db.Close()

	cfg := cfg{
		db:  db,
		env: env,
		q:   database.New(db),
	}

	switch cmd {
	case "all":
		exitWithErr(handleAll(cfg))
	case "groups":
		exitWithErr(handleGroups(cfg, args...))
	case "inventories":
		exitWithErr(handleInventories(cfg, args...))
	case "items":
		exitWithErr(handleItems(cfg, args...))
	case "wipeall":
		exitWithErr(handleWipeAll(cfg))
	case "wipe":
		exitWithErr(handleWipe(cfg, args...))
	case "help":
		exitWithErr(handleHelp(args...))
	default:
		exitWithErr(fmt.Errorf("unknown command: %s\n\nrun command 'help' for more information", cmd))
	}

	fmt.Println()
}

func setup() (*sql.DB, env.Env, error) {
	e := env.GetEnv()
	dbUrl := e.DB_URL
	platform := e.PLATFORM
	if strings.ToLower(platform) != "dev" {
		return nil, env.Env{}, fmt.Errorf("this should only be run in dev environment")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, env.Env{}, fmt.Errorf("couldn't open database: %v\n", err)
	}
	if err := db.Ping(); err != nil {
		return nil, env.Env{}, fmt.Errorf("couldn't connect to database: %v\n", err)
	}
	return db, e, nil
}
