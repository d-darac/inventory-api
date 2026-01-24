package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func argToKeyVal(arg string) (key, value string, err error) {
	keyValue := strings.Split(arg, "=")
	if len(keyValue) != 2 {
		return "", "", fmt.Errorf("format of optional arguments must be `key=value`")
	}
	key = keyValue[0]
	value = keyValue[1]
	return key, value, nil
}

func exitWithErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printIDs(ids []uuid.UUID) {
	fmt.Println("----------------------------------------------------------------")
	n := min(len(ids), 5)
	if n < len(ids) {
		fmt.Printf("Displaying first %d IDs...\n", n)
	}
	for _, item := range ids[:n] {
		fmt.Printf("%v\n", item)
	}
	fmt.Println("----------------------------------------------------------------")
}
