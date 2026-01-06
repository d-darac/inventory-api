package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d-darac/inventory-assets/sql"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cli <action>")
		os.Exit(1)
	}

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	migrations := sql.DbMigrations{
		DbURL: dbUrl,
	}

	cmdName := os.Args[1]

	switch cmdName {
	case "up":
		{
			migrations.Up()
		}
	case "down":
		{
			migrations.Down()
		}
	case "up-to":
		if len(os.Args) < 3 {
			fmt.Println("usage: cli up-to <version>")
			os.Exit(1)
		}
		cmdVersion := os.Args[2]
		v, err := strconv.Atoi(cmdVersion)
		if err != nil {
			fmt.Println("value of version argument must be numeric")
		}
		migrations.UpTo(int64(v))
	case "down-to":
		if len(os.Args) < 3 {
			fmt.Println("usage: cli down-to <version>")
			os.Exit(1)
		}
		cmdVersion := os.Args[2]
		v, err := strconv.Atoi(cmdVersion)
		if err != nil {
			fmt.Println("value of version argument must be numeric")
		}
		migrations.DownTo(int64(v))
	default:
		fmt.Printf("command not found: %s\n", cmdName)
		fmt.Println("available commands:")
		fmt.Print("- up\n- up-to\n- down\n- down-to\n")
		os.Exit(1)
	}
}
