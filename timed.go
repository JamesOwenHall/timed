package timed

type Iterator interface {
	Next() *Record
}

type Aggregator interface {
	Name() string
	Next(Value)
	Final() Value
}

type AggregatorCall struct {
	Key        string
	Aggregator Aggregator
}
