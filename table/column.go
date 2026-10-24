package table

import (
	"db-architect/types"
	"strings"
)

type Column struct {
	Name     string `json:"name"`
	Primary  bool   `json:"primary"`
	Nullable bool   `json:"nullable"`
	AutoIncr bool   `json:"autoIncrement"`

	References ForeignKey `json:"references"`

	Type    types.DBType
	RawType any `json:"type"`
}

func (c Column) Represent() string {
	b := strings.Builder{}
	b.WriteString(c.Name)
	b.WriteString(" ")
	b.WriteString(c.Type.Represent())

	if !c.Nullable && !c.Primary && !c.AutoIncr {
		b.WriteString(" NOT NULL")
	}

	if c.AutoIncr {
		b.WriteString(" AUTO_INCREMENT")
	}

	if c.Primary {
		b.WriteString(" PRIMARY KEY")
	}

	return b.String()
}

type ForeignKey struct {
	Self string

	Table   string `json:"table"`
	Column  string `json:"column"`
	Cascade bool   `json:"cascade"`
}

func (f ForeignKey) Represent() string {
	b := strings.Builder{}
	// FOREIGN KEY (UserId)
	b.WriteString("FOREIGN KEY (")
	b.WriteString(f.Self)
	b.WriteString(") ")

	// ) REFERENCES Users(UserId)
	b.WriteString("REFERENCES ")
	b.WriteString(f.Table)
	b.WriteString("(")
	b.WriteString(f.Column)
	b.WriteString(")")

	if f.Cascade {
		b.WriteString(" ON DELETE CASCADE")
	}

	return b.String()
}
