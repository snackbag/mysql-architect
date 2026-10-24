package types

import (
	"fmt"
)

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
		return nil, fmt.Errorf("invalid database type: %s", raw.Name)
	}
}
