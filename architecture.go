package db_architect

import (
	"db-architect/table"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
