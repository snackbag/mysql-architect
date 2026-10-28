package types

type Boolean struct {
}

func (b Boolean) Kind() string {
	return "bool"
}

func (b Boolean) Represent() string {
	return "BOOLEAN"
}

func (b Boolean) Validate(value any) bool {
	_, ok := value.(bool)
	return ok
}
