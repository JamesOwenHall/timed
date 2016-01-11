package executor

// ValueType describes the type of a Value.
type ValueType int

const (
	Unknown ValueType = iota
	Int64
	Int
	Float64
	Float32
	Boolean
	String
)

// Value is arbitrary data with type information.
type Value struct {
	Type ValueType
	Data interface{}
}

// Record is a map from string keys to Values.
type Record map[string]Value
