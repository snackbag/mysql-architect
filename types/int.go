package types

type Int struct {
}

func (i Int) Kind() string {
	return "integer"
}

func (i Int) Represent() string {
	return "INT"
}

func (i Int) Validate(value any) bool {
	_, ok := value.(int)
	return ok
}
