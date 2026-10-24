package database

import (
	"db-architect/types"
	"encoding/json"
	"fmt"
)

func UnmarshalTable(data []byte) (Table, error) {
	var table Table
	err := json.Unmarshal(data, &table)
	if err != nil {
		return table, fmt.Errorf("failed to unmarshal table: %w", err)
	}

	for _, col := range table.Columns {
		t, err := types.GetType(types.ToFormat(col.RawType))
		if err != nil {
			return table, fmt.Errorf("column %s contains invalid type: %w", col.Name, err)
		}
		col.Type = t
	}

	return table, err
}
