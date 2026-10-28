package types

import (
	"fmt"
	"time"
)

type Date struct {
}

func (d Date) Kind() string {
	return "date"
}

func (d Date) Represent() string {
	return "DATE"
}

func (d Date) Validate(value any) bool {
	_, ok := value.(time.Time)
	return ok
}

type Time struct {
	Fsp uint8
}

func (t Time) Kind() string {
	return "time"
}

func (t Time) Represent() string {
	if t.Fsp > 0 {
		return fmt.Sprintf("TIME(%d)", t.Fsp)
	}
	return "TIME"
}

func (t Time) Validate(value any) bool {
	_, ok := value.(time.Time)
	return ok
}

type Datetime struct {
	Fsp uint8
}

func (d Datetime) Kind() string {
	return "datetime"
}

func (d Datetime) Represent() string {
	if d.Fsp > 0 {
		return fmt.Sprintf("DATETIME(%d)", d.Fsp)
	}
	return "DATETIME"
}

func (d Datetime) Validate(value any) bool {
	_, ok := value.(time.Time)
	return ok
}

type Timestamp struct {
	Fsp uint8
}

func (d Timestamp) Kind() string {
	return "timestamp"
}

func (d Timestamp) Represent() string {
	if d.Fsp > 0 {
		return fmt.Sprintf("TIMESTAMP(%d)", d.Fsp)
	}
	return "TIMESTAMP"
}

func (d Timestamp) Validate(value any) bool {
	_, ok := value.(time.Time)
	return ok
}
