package macro

import "github.com/snackbag/mysql-architect/types"

type Parameter struct {
	Name     string
	Fallback any `json:"fallback"`

	Type    types.DBType
	RawType any `json:"type"`
}

// Result is returned after a macro finished executing and contains its raw values.
// For safe reading, it is recommended to use Get or GetOrDefault.
type Result struct {
	Raw map[string]any // Raw is a map containing all raw information the macro was tasked to extract.
}

// Get returns the value assigned to the specified key of a macro's return result,
// in addition to a boolean indicator showing whether the key was valid.
func Get[S any](res Result, key string) (S, bool) {
	d, ok := res.Raw[key]
	return d.(S), ok
}

// GetOrDefault returns the value assigned to the specified key of a macro's return result,
// in addition to a boolean indicator showing whether the key was valid. If no value was assigned
// to the key, it returns the fallback value.
func GetOrDefault[S any](res Result, key string, fallback S) (S, bool) {
	d, ok := res.Raw[key]
	if !ok {
		return fallback, ok
	}
	return d.(S), ok
}
