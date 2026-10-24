package main

import (
	architect "db-architect"
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
}
