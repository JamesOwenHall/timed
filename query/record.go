package query

// ValueType describes the type of a Value.
type ValueType int

const (
	Null ValueType = iota
	Boolean
	Number
	String
)

// Value is arbitrary data with type information.
type Value struct {
	Type ValueType
	Data interface{}
}

// Record is a map from string keys to Values.
type Record map[string]Value
