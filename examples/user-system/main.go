package main

import (
	"fmt"

	architect "github.com/snackbag/mysql-architect"
)

const username = "architectUser"
const password = "dev"
const dbName = "architectExample"

func main() {
	arch, err := architect.LoadArchitecture("db")
	if err != nil {
		panic(err)
	}

	db, err := arch.InitDB(username, password, dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	macro := arch.Macros[0]
	rv, err := macro.Exec(db, map[string]any{"Username": "Test", "Password": "TestPassword", "Salt": "TestSalt"})
	if err != nil {
		panic(err)
	}
	fmt.Println(rv["UserId"])
}
