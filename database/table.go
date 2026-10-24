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

	for _, col := range t.Columns {
		columns = append(columns, "\t"+col.Represent())
	}
	b.WriteString(strings.Join(columns, ",\n"))

	b.WriteString("\n)")

	return b.String()
}
