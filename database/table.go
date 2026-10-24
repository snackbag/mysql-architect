package database

import "strings"

type Table struct {
	Name    string
	Columns []*Column `json:"columns"`
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
