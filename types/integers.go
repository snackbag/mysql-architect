package types

import "strings"

type IntConstraint interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int
}

type IntSize int

const (
	TinyInt IntSize = iota
	SmallInt
	MediumInt
	Int
	BigInt
)

var intName = map[IntSize]string{
	TinyInt:   "TINY",
	SmallInt:  "SMALL",
	MediumInt: "MEDIUM",
	Int:       "",
	BigInt:    "BIG",
}

func (s IntSize) Name() string {
	return intName[s]
}

type IntType[T IntConstraint] struct {
	Unsigned bool
	Size     IntSize
}

func (i IntType[T]) Kind() string {
	name := strings.ToLower(i.Size.Name()) + " integer"
	if i.Size == Int { // remove space
		name = "integer"
	}

	if i.Unsigned {
		return name + " (unsigned)"
	}
	return name
}

func (i IntType[T]) Represent() string {
	name := i.Size.Name() + "INT"

	if i.Unsigned {
		return name + " UNSIGNED"
	}
	return name
}

func (i IntType[T]) Validate(value any) bool {
	_, ok := value.(T)
	return ok
}
