package main

import (
	"fmt"

	architect "github.com/snackbag/mysql-architect"
	"github.com/snackbag/mysql-architect/macros"
)

const username = "architectUser"
const password = "dev"
const dbName = "architectExample"

func main() {
	arch, err := architect.LoadArchitecture("db")
	if err != nil {
		panic(err)
	}

	db, err := arch.Init(username, password, dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create users
	newUser := arch.Macros["macros/new_user"]
	userRes, err := newUser.Exec(db, map[string]any{"Username": "Test", "Password": "TestPassword", "Salt": "TestSalt"})
	if err != nil {
		panic(err)
	}

	newAdmin := arch.Macros["macros/new_admin"]
	adminRes, err := newAdmin.Exec(db, map[string]any{"Username": "Testwo", "Password": "TestPassword", "Salt": "TestSalt"})
	if err != nil {
		panic(err)
	}

	// Extract values from queries
	newUserId, _ := macros.Get[int](userRes, "UserId") // user id of normal user
	fmt.Println(fmt.Sprintf("user id (of the normal user): %d", newUserId))

	newUserId, _ = macros.Get[int](adminRes, "UserId") // user id of admin user
	fmt.Println(fmt.Sprintf("user id (of the admin user): %d", newUserId))

	// Get username of admin
	grabUser := arch.Macros["macros/grab_user"]
	grabRes, err := grabUser.Exec(db, map[string]any{"UserId": newUserId})
	if err != nil {
		panic(err)
	}

	newUsername, _ := macros.GetOrFallback[string](grabRes, "Username", "(no username found)")
	fmt.Println(fmt.Sprintf("admin username: %s", newUsername))
}
