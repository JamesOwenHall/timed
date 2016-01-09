package timed

type ValueType int

const (
	Null ValueType = iota
	Boolean
	Number
	String
)

type Value struct {
	Type ValueType
	Data interface{}
}

type Record map[string]Value
