package query

import (
	"fmt"
	"strings"

	"github.com/JamesOwenHall/timed/executor"
)

func NewAggregator(name string) executor.Aggregator {
	switch strings.ToLower(name) {
	case "count":
		return new(CountAggregator)
	case "sum":
		return new(SumAggregator)
	case "avg":
		return new(AvgAggregator)
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

func (c *CountAggregator) Next(v interface{}) error {
	c.count++
	return nil
}

func (c *CountAggregator) Final() interface{} {
	return c.count
}

type SumAggregator struct {
	intSum   int64
	floatSum float64
}

func (s *SumAggregator) Name() string {
	return "sum"
}

func (s *SumAggregator) Next(v interface{}) error {
	switch v := v.(type) {
	case int:
		s.intSum += int64(v)
	case int64:
		s.intSum += v
	case float32:
		s.floatSum += float64(v)
	case float64:
		s.floatSum += v
	default:
		return &ErrInvalidArg{Type: fmt.Sprintf("%T", v), Function: "sum"}
	}

	return nil
}

func (s *SumAggregator) Final() interface{} {
	if s.intSum != 0 {
		return s.intSum
	}

	return s.floatSum
}

type AvgAggregator struct {
	intSum   int64
	floatSum float64
	count    int64
}

func (a *AvgAggregator) Name() string {
	return "avg"
}

func (a *AvgAggregator) Next(v interface{}) error {
	a.count++

	switch v := v.(type) {
	case int:
		a.intSum += int64(v)
	case int64:
		a.intSum += v
	case float32:
		a.floatSum += float64(v)
	case float64:
		a.floatSum += v
	default:
		return &ErrInvalidArg{Type: fmt.Sprintf("%T", v), Function: "avg"}
	}

	return nil
}

func (a *AvgAggregator) Final() interface{} {
	if a.count == 0 {
		return 0
	}

	if a.intSum != 0 {
		return float64(a.intSum) / float64(a.count)
	}

	return a.floatSum / float64(a.count)
}

type ErrInvalidArg struct {
	Type     string
	Function string
}

func (e *ErrInvalidArg) Error() string {
	return fmt.Sprintf("invalid type %s for function %s", e.Type, e.Function)
}
