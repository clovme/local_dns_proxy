package enums

type Enums struct {
	Key  string
	Name string
	Desc string
}

type Enum[T any] interface {
	Key() string
	Name() string
	Desc() string
	Int() int
	Is(v T) bool
}
