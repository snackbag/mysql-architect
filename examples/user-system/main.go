package main

import (
	architect "db-architect"
)

func main() {
	arch, err := architect.LoadArchitecture("db")
	if err != nil {
		panic(err)
	}

	for _, table := range arch.Tables {
		println(table.Represent())
	}
}
