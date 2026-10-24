package types

type DBType interface {
	Kind() string
	Represent() string
	Validate(value any) bool
}
