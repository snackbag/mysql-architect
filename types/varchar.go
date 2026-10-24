package types

import "fmt"

type Varchar struct {
	Length int
}

func (v Varchar) Name() string {
	return fmt.Sprintf("varchar(length=%d)", v.Length)
}

func (v Varchar) Represent() string {
	return fmt.Sprintf("VARCHAR(%d)", v.Length)
}

func (v Varchar) Validate(value any) bool {
	s, ok := value.(string)
	return ok && len(s) <= v.Length
}
