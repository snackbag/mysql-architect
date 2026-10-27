package macro

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snackbag/mysql-architect/types"
)

type Macro struct {
	Name    string
	Params  map[string]Parameter
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

	macro := Macro{Name: name, Params: make(map[string]Parameter), Program: raw.Program}
	for key, param := range raw.Params {
		t, err := types.GetType(types.ToFormat(param.RawType))
		if err != nil {
			return macro, fmt.Errorf("macro %q contains invalid type: %w", name, err)
		}

		param.Type = t
		param.Name = key
		macro.Params[key] = param
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

func (m Macro) Exec(db *sql.DB, globalParams map[string]any) (Result, error) {
	if m.Name == "" {
		return Result{make(map[string]any)}, errors.New("macro not found")
	}

	// ensure all parameters are set
	for _, macroParam := range m.Params {
		_, ok := globalParams[macroParam.Name]
		if !ok {
			if macroParam.Fallback == nil {
				return Result{}, fmt.Errorf("cannot execute macro %s, because parameter %s is not set and has no fallback value", m.Name, macroParam.Name)
			}

			globalParams[macroParam.Name] = macroParam.Fallback
		}
	}

	raw := make(map[string]any) // return values

	// begin tx
	tx, err := db.Begin()
	if err != nil {
		return Result{raw}, err
	}
	defer tx.Rollback()

	// go through each script and execute
	for _, script := range m.Program {
		returned := make([]any, 0)
		sqlParams := make(map[string]any)
		for param, link := range script.Values {
			value, ok := globalParams[link]
			if !ok {
				return Result{raw}, fmt.Errorf("script %s wants to fill in SQL parameter %s with macro parameter %s, but the macro parameter does not exist", script.File, param, link)
			}
			sqlParams[param] = value
		}

		// execute each statement
		for _, stmt := range script.Statements {
			args := orderArgs(stmt, sqlParams)
			values, err := stmt.Exec(tx, args)
			if err != nil {
				return Result{raw}, fmt.Errorf("execute script in %q: %w", script.File, err)
			}
			returned = append(returned, values...)
		}

		if len(returned) < len(script.Return) {
			return Result{raw}, fmt.Errorf(
				"script %q returned %d values, but %d were expected",
				script.File, len(returned), len(script.Return),
			)
		}

		for i, name := range script.Return {
			raw[name] = returned[i]

			if script.OverwriteParams {
				globalParams[name] = returned[i]
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return Result{raw}, err
	}

	return Result{raw}, nil
}

func orderArgs(stmt Statement, params map[string]any) []any {
	args := make([]any, 0)

	for _, name := range stmt.ParamMapping {
		args = append(args, params[name])
	}

	return args
}
