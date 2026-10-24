package macro

import (
	"database/sql"
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

func (m Macro) Exec(db *sql.DB, globalParams map[string]any) (map[string]any, error) {
	rv := make(map[string]any) // return values
	returned := make([]any, 0)

	tx, err := db.Begin()
	if err = handleTxError(tx, err); err != nil {
		return rv, err
	}

	for _, script := range m.Program {
		sqlParams := make(map[string]any)
		for param, link := range script.Values {
			sqlParams[param] = globalParams[link]
		}

		for _, stmt := range script.Statements {
			args := orderArgs(stmt, sqlParams)
			values, err := stmt.Exec(tx, args)
			if err = handleTxError(tx, err); err != nil {
				return rv, fmt.Errorf("execute script in %q: %w", script.File, err)
			}
			returned = append(returned, values...)
		}

		if len(returned) < len(script.Return) {
			if err = handleTxError(tx, err); err != nil {
				return rv, fmt.Errorf(
					"script %q returned %d values, but %d were expected",
					script.File, len(returned), len(script.Return),
				)
			}
		}

		for i, name := range script.Return {
			rv[name] = returned[i]
		}
	}

	err = tx.Commit()
	if err = handleTxError(tx, err); err != nil {
		return rv, err
	}

	return rv, nil
}

func handleTxError(tx *sql.Tx, err error) error {
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback Tx (%w): %w", err, rollbackErr)
		}

		return err
	}

	return nil
}

func orderArgs(stmt Statement, params map[string]any) []any {
	args := make([]any, 0)

	for _, name := range stmt.ParamMapping {
		args = append(args, params[name])
	}

	return args
}
