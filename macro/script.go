package macro

import (
	"regexp"
	"strings"
)

var paramRegex = regexp.MustCompile("@([A-Za-z_0-9]+)")

type Script struct {
	File   string            `json:"file"`
	Values map[string]string `json:"values"`
	Return []string          `json:"return"`

	Statements []Statement
}

type Statement struct {
	ParamMapping []string
	Value        string
}

func ParseScript(data []byte) []Statement {
	stmts := make([]Statement, 0)

	lines := strings.Split(string(data), ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		stmt := Statement{
			Value:        paramRegex.ReplaceAllString(line, "?"),
			ParamMapping: make([]string, 0),
		}

		for _, submatch := range paramRegex.FindAllStringSubmatch(line, -1) {
			stmt.ParamMapping = append(stmt.ParamMapping, submatch[1])
		}

		stmts = append(stmts, stmt)
	}

	return stmts
}
