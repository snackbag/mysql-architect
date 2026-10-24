package types

import (
	"encoding/json"
	"fmt"
)

func ToFormat(v any) Format {
	switch v.(type) {
	case string:
		return Format{Name: v.(string)}
	case Format:
		return v.(Format)
	case map[string]any:
		var f Format
		m, _ := json.Marshal(v)
		json.Unmarshal(m, &f)
		return f
	default:
		return Format{Name: "failed ToFormat conversion"}
	}
}

type Format struct {
	Name   string `json:"name"`
	Length int    `json:"length"`

	Link   string `json:"link"`
	Column string `json:"column"`
}

func GetType(raw Format) (DBType, error) {
	switch raw.Name {
	case "text":
		return Text{}, nil
	case "int":
		return Int{}, nil
	case "varchar":
		return Varchar{Length: raw.Length}, nil
	default:
		return nil, fmt.Errorf("invalid DBType: %s", raw.Name)
	}
}
