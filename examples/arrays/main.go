package main

import (
	"database/sql"
	"fmt"
	"time"

	architect "github.com/snackbag/mysql-architect"
	"github.com/snackbag/mysql-architect/macros"
)

const username = "architectUser"
const password = "dev"
const dbName = "architectExample"

var architecture architect.Architecture
var db *sql.DB

func main() {
	var err error

	// initialisation
	architecture, err = architect.LoadArchitecture("db")
	if err != nil {
		panic(err)
	}

	db, err = architecture.Init(username, password, dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// database actions
	john := registerCustomer("John")
	mary := registerCustomer("Mary")

	registerInvoice(john, 4000)
	registerInvoice(john, 2000)
	registerInvoice(john, 8500)

	registerInvoice(mary, 200)

	// reading from database
	fmt.Println(fmt.Sprintf("John's payments: %d", getInvoices(john)))
	fmt.Println(fmt.Sprintf("Mary's payments: %d", getInvoices(mary)))
}

func registerCustomer(name string) int {
	macro, _ := architecture.Macros["macros/new_customer"]
	result, err := macro.Exec(db, map[string]any{"Name": name, "CurrentTime": time.Now()})
	if err != nil {
		panic(err)
	}

	customerId, _ := macros.Get[int](result, "CustomerId")
	return customerId
}

func registerInvoice(customerId int, amount int) {
	macro, _ := architecture.Macros["macros/new_invoice"]
	_, err := macro.Exec(db, map[string]any{"CustomerId": customerId, "Amount": amount, "CurrentTime": time.Now()})
	if err != nil {
		panic(err)
	}
}

func getInvoices(customerId int) []int {
	macro, _ := architecture.Macros["macros/get_invoices"]
	result, err := macro.Exec(db, map[string]any{"CustomerId": customerId})
	if err != nil {
		panic(err)
	}

	invoices, _ := macros.Get[[]int](result, "InvoiceAmounts")
	return invoices
}
