package types

type Text struct {
}

func (t Text) Kind() string {
	return "text"
}

func (t Text) Represent() string {
	return "TEXT"
}

func (t Text) Validate(value any) bool {
	_, ok := value.(string)
	return ok
}
