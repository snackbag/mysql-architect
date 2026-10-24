package db_architect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/snackbag/mysql-architect/macro"
	"github.com/snackbag/mysql-architect/table"

	_ "github.com/go-sql-driver/mysql"
)

type Architecture struct {
	Tables map[string]table.Table
	Macros map[string]macro.Macro
}

func LoadArchitecture(dir string) (Architecture, error) {
	// load file
	var raw struct {
		Tables []string `json:"tables"`
		Macros []string `json:"macros"`
	}

	data, err := os.ReadFile(filepath.Join(dir, "architecture.json"))
	if err != nil {
		return Architecture{}, fmt.Errorf("failed to read architecture.json: %w", err)
	}

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return Architecture{}, fmt.Errorf("failed to unmarshal architecture: %w", err)
	}

	// load tables
	arch := Architecture{Tables: make(map[string]table.Table), Macros: make(map[string]macro.Macro)}
	for _, t := range raw.Tables {
		data, err = os.ReadFile(filepath.Join(dir, filepath.FromSlash(t)+".json"))

		unmarshaled, err := table.UnmarshalTable(data)
		if err != nil {
			return arch, err
		}
		unmarshaled.Name = t[strings.LastIndex(t, "/")+1:]
		arch.Tables[unmarshaled.Name] = unmarshaled
	}

	// load macros
	for _, m := range raw.Macros {
		data, err = os.ReadFile(filepath.Join(dir, filepath.FromSlash(m)+".json"))

		loaded, err := macro.LoadMacro(dir, m[strings.LastIndex(m, "/")+1:], data)
		if err != nil {
			return arch, err
		}
		arch.Macros[loaded.Name] = loaded
	}

	return arch, nil
}

// InitDB opens a connection to the database and creates the defined tables. Returns the database
// connection and an error. An error does not necessarily mean the database wasn't connected.
// DO NOT FORGET to close the connection!
func (a Architecture) InitDB(username string, password string, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", username, password, dbName))

	if db == nil || err != nil {
		return db, err
	}

	tableErrs := make([]string, 0)
	for _, t := range a.Tables {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + t.Represent())
		if err != nil {
			tableErrs = append(tableErrs, "\""+err.Error()+"\"")
		}
	}

	if len(tableErrs) > 0 {
		return db, fmt.Errorf("failed to create all tables: %s", strings.Join(tableErrs, ", "))
	}

	return db, nil
}
