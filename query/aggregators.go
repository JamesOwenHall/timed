package query

import (
	"strings"

	"github.com/JamesOwenHall/timed/executor"
)

func NewAggregator(name string) executor.Aggregator {
	switch strings.ToLower(name) {
	case "count":
		return new(CountAggregator)
	default:
		return nil
	}
}

// CountAggregator counts the number of occurrences of the field.
type CountAggregator struct {
	count int64
}

func (c *CountAggregator) Name() string {
	return "count"
}

func (c *CountAggregator) Next(v executor.Value) error {
	c.count++
	return nil
}

func (c *CountAggregator) Final() executor.Value {
	return executor.Value{
		Type: executor.Int64,
		Data: c.count,
	}
}
