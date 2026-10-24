package macro

import (
	"database/sql"
	"regexp"
	"strings"
)

var paramRegex = regexp.MustCompile("@([A-Za-z_0-9]+)")

type Script struct {
	File            string            `json:"file"`
	Values          map[string]string `json:"values"`
	Return          []string          `json:"return"`
	OverwriteParams bool              `json:"overwriteParams"`

	Statements []Statement
}

type Statement struct {
	ParamMapping []string
	Query        string
}

func ParseScript(data []byte) []Statement {
	stmts := make([]Statement, 0)

	lines := strings.Split(string(data), ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		stmt := Statement{
			Query:        paramRegex.ReplaceAllString(line, "?"),
			ParamMapping: make([]string, 0),
		}

		for _, submatch := range paramRegex.FindAllStringSubmatch(line, -1) {
			stmt.ParamMapping = append(stmt.ParamMapping, submatch[1])
		}

		stmts = append(stmts, stmt)
	}

	return stmts
}

func (s Statement) Exec(tx *sql.Tx, args []any) ([]any, error) {
	rows, err := tx.Query(s.Query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returned []any

	for {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			values := make([]any, len(columns))
			destinations := make([]any, len(columns))

			for i := range values {
				destinations[i] = &values[i]
			}

			if err := rows.Scan(destinations...); err != nil {
				return nil, err
			}

			returned = append(returned, values...)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		if !rows.NextResultSet() {
			break
		}
	}

	return returned, rows.Err()
}
