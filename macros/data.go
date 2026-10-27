package macros

import (
	"database/sql"
	"errors"

	"github.com/snackbag/mysql-architect/types"
)

type Parameter struct {
	Name     string
	Fallback any `json:"fallback"`

	Type    types.DBType
	RawType any `json:"type"`
}

// Result is returned after a macro finished executing and contains its raw values.
// For safe reading, it is recommended to use Get or GetOrFallback.
type Result struct {
	Raw map[string]any // Raw is a map containing all raw information the macro was tasked to extract.
}

// Get returns the value assigned to the specified key of a macro's return result,
// in addition to a boolean indicator showing whether the key was valid. It automatically
// converts the return value into the requested type.
func Get[S any](result Result, key string) (S, bool) {
	var zero S
	return GetOrFallback[S](result, key, zero)
}

// GetOrFallback returns the value assigned to the specified key of a macro's return result,
// in addition to a boolean indicator showing whether the key was valid. If no value was assigned
// to the key, it returns the fallback value. It automatically converts the return value into
// the requested type.
func GetOrFallback[S any](result Result, key string, fallback S) (S, bool) {
	d, ok := result.Raw[key]
	if !ok {
		return fallback, false
	}

	converted, err := convertGet[S](d)
	if err != nil {
		return fallback, false
	}
	return converted, ok
}

func convertGet[S any](val any) (S, error) {
	var result sql.Null[S]
	var zero S

	if err := result.Scan(val); err != nil {
		return zero, err
	}

	if !result.Valid {
		return zero, errors.New("cannot convert null")
	}

	return result.V, nil
}
