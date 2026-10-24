package main

import (
	"fmt"
	"os"
)

func main() {
	usersFile, err := os.ReadFile("examples/user-system/db/tables/UserDiscordAccounts.json")
	if err != nil {
		panic(err)
	}

	usersTable, err := table.UnmarshalTable(usersFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(usersTable.Represent())
}
