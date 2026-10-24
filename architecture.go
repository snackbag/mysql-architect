package db_architect

import (
	"database/sql"
	"db-architect/table"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Architecture struct {
	Tables []table.Table
}

func LoadArchitecture(dir string) (Architecture, error) {
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

	arch := Architecture{Tables: make([]table.Table, 0)}
	for _, t := range raw.Tables {
		data, err = os.ReadFile(filepath.Join(dir, filepath.FromSlash(t)+".json"))

		unmarshaled, err := table.UnmarshalTable(data)
		if err != nil {
			return arch, err
		}
		unmarshaled.Name = t[strings.LastIndex(t, "/")+1:]
		arch.Tables = append(arch.Tables, unmarshaled)
	}

	return arch, nil
}

// InitDB opens a connection to the database. DO NOT FORGET to close this connection to avoid memory leaks!
func (a Architecture) InitDB(username string, password string, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", username, password, dbName))
	return db, err
}
