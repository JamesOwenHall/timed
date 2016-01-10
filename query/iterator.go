package query

// Iterator represents any data source that can be iterated over.  Next returns
// nil once the end of the data has been reached.
type Iterator interface {
	Next() (Record, error)
}

// Aggregator represents an aggregation function.
type Aggregator interface {
	// Name returns the name of the aggregator.
	Name() string
	// Next is called over each value of a specified key.
	Next(Value) error
	// Final returns the final value of the aggregation.
	Final() Value
}

// AggregatorCall represents a aggregator function call.
type AggregatorCall struct {
	Key        string
	Aggregator Aggregator
	Alias      string
}
