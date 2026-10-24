package types

type DBType interface {
	Name() string
	Represent() string
	Validate(value any) bool
}
