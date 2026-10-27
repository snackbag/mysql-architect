package table

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/snackbag/mysql-architect/types"
)

type Table struct {
	Name    string
	Columns []Column `json:"columns"`
}

func (t Table) Represent() string {
	b := strings.Builder{}
	b.WriteString(t.Name)
	b.WriteString(" (\n")

	columns := make([]string, 0)
	foreignKeys := make([]string, 0)

	for _, col := range t.Columns {
		columns = append(columns, "\t"+col.Represent())

		if col.References.Self != "" {
			foreignKeys = append(foreignKeys, "\t"+col.References.Represent())
		}
	}
	b.WriteString(strings.Join(columns, ",\n"))

	if len(foreignKeys) > 0 {
		b.WriteString(",\n\n")
		b.WriteString(strings.Join(foreignKeys, ",\n"))
	}

	b.WriteString("\n)")

	return b.String()
}

func UnmarshalTable(name string, data []byte) (Table, error) {
	var table Table
	err := json.Unmarshal(data, &table)
	if err != nil {
		return table, fmt.Errorf("failed to unmarshal table %s: %w", name, err)
	}

	fixedCols := make([]Column, 0)
	for _, col := range table.Columns {
		t, err := types.GetType(types.ToFormat(col.RawType))
		if err != nil {
			return table, fmt.Errorf("column %s of table %s contains invalid type: %w", col.Name, name, err)
		}
		col.Type = t

		if col.References.Table != "" {
			col.References.Self = col.Name
		}

		fixedCols = append(fixedCols, col)
	}

	table.Columns = fixedCols
	table.Name = name

	return table, err
}
