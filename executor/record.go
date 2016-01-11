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
	Time
)

func (v ValueType) String() string {
	switch v {
	case Int64:
		return "int64"
	case Int:
		return "int"
	case Float64:
		return "float64"
	case Float32:
		return "float32"
	case Boolean:
		return "bool"
	case String:
		return "string"
	case Time:
		return "time"
	default:
		return "unknown"
	}
}

// Value is arbitrary data with type information.
type Value struct {
	Type ValueType
	Data interface{}
}

// Record is a map from string keys to Values.
type Record map[string]Value
