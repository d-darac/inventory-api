package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d-darac/inventory-assets/sql"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: cli <connection_string> <action>")
		os.Exit(1)
	}

	connectionString := os.Args[1]

	migrations := sql.DbMigrations{
		DbURL: connectionString,
	}

	action := os.Args[2]

	switch action {
	case "up":
		{
			migrations.Up()
		}
	case "down":
		{
			migrations.Down()
		}
	case "up-to":
		if len(os.Args) < 4 {
			fmt.Println("missing <version> argument")
			fmt.Println("usage: cli <connection_string> up-to <version>")
			os.Exit(1)
		}
		version := os.Args[3]
		v, err := strconv.Atoi(version)
		if err != nil {
			fmt.Println("value of version argument must be numeric")
			os.Exit(1)
		}
		migrations.UpTo(int64(v))
	case "down-to":
		if len(os.Args) < 4 {
			fmt.Println("missing <version> argument")
			fmt.Println("usage: cli <connection_string> down-to <version>")
			os.Exit(1)
		}
		version := os.Args[3]
		v, err := strconv.Atoi(version)
		if err != nil {
			fmt.Println("value of version argument must be numeric")
			os.Exit(1)
		}
		migrations.DownTo(int64(v))
	default:
		fmt.Printf("command not found: %s\n", action)
		fmt.Println("available commands:")
		fmt.Print("- up\n- up-to\n- down\n- down-to\n")
		os.Exit(1)
	}
}
