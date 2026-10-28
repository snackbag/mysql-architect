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
	// text
	case "text":
		return Text{}, nil
	case "varchar":
		return Varchar{Length: raw.Length}, nil

	// integers
	case "int":
		return IntType[int]{Size: Int, Unsigned: false}, nil
	case "uint":
		return IntType[uint]{Size: Int, Unsigned: true}, nil
	case "tinyint":
		return IntType[int8]{Size: TinyInt, Unsigned: false}, nil
	case "utinyint":
		return IntType[uint8]{Size: TinyInt, Unsigned: true}, nil
	case "smallint":
		return IntType[int16]{Size: SmallInt, Unsigned: false}, nil
	case "usmallint":
		return IntType[uint16]{Size: SmallInt, Unsigned: true}, nil
	case "mediumint":
		return IntType[int32]{Size: MediumInt, Unsigned: false}, nil
	case "umediumint":
		return IntType[uint32]{Size: MediumInt, Unsigned: true}, nil
	case "bigint":
		return IntType[int64]{Size: BigInt, Unsigned: false}, nil
	case "ubigint":
		return IntType[uint64]{Size: BigInt, Unsigned: true}, nil

	// bool
	case "bool", "boolean":
		return Boolean{}, nil

	default:
		return nil, fmt.Errorf("invalid DBType: %s", raw.Name)
	}
}
