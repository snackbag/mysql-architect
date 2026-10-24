package macro

import "github.com/snackbag/mysql-architect/types"

type Parameter struct {
	Name     string
	Fallback any `json:"fallback"`

	Type    types.DBType
	RawType any `json:"type"`
}
