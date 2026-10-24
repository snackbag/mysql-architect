package macro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snackbag/mysql-architect/types"
)

type Macro struct {
	Name    string
	Params  []Parameter
	Program []Script `json:"execute"`
}

func LoadMacro(baseDir string, name string, data []byte) (Macro, error) {
	var raw struct {
		Params  map[string]Parameter `json:"params"`
		Program []Script             `json:"execute"`
	}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return Macro{}, fmt.Errorf("failed to unmarshal macro %q: %w", name, err)
	}

	macro := Macro{Name: name, Params: make([]Parameter, 0), Program: raw.Program}
	for key, param := range raw.Params {
		t, err := types.GetType(types.ToFormat(param.RawType))
		if err != nil {
			return macro, fmt.Errorf("macro %q contains invalid type: %w", name, err)
		}

		param.Type = t
		param.Name = key
		macro.Params = append(macro.Params, param)
	}

	// load scripts
	newProgram := make([]Script, 0)
	for _, script := range macro.Program {
		data, err = os.ReadFile(filepath.Join(baseDir, filepath.FromSlash(script.File)))
		if err != nil {
			return macro, fmt.Errorf("failed to read script %q from macro %q, because: %w", script.File, name, err)
		}

		script.Statements = ParseScript(data)
		newProgram = append(newProgram, script)
	}

	macro.Program = newProgram

	return macro, nil
}

func (m Macro) Exec(map[string]any) (map[string]any, error) {
	return nil, nil
}
