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

type SumAggregator struct {
	intSum   int64
	floatSum float64
}

func (s *SumAggregator) Name() string {
	return "sum"
}

func (s *SumAggregator) Next(v executor.Value) error {
	switch v.Type {
	case executor.Int64:
		s.intSum += v.Data.(int64)
	case executor.Int:
		s.intSum += int64(v.Data.(int))
	case executor.Float64:
		s.floatSum += v.Data.(float64)
	case executor.Float32:
		s.floatSum += float64(v.Data.(float32))
	default:
		return &ErrInvalidArg{Type: v.Type, Function: "sum"}
	}

	return nil
}

func (s *SumAggregator) Final() executor.Value {
	if s.intSum != 0 {
		return executor.Value{Type: executor.Int64, Data: s.intSum}
	}

	return executor.Value{Type: executor.Float64, Data: s.floatSum}
}

type AvgAggregator struct {
	intSum   int64
	floatSum float64
	count    int64
}

func (a *AvgAggregator) Name() string {
	return "avg"
}

func (a *AvgAggregator) Next(v executor.Value) error {
	a.count++

	switch v.Type {
	case executor.Int64:
		a.intSum += v.Data.(int64)
	case executor.Int:
		a.intSum += int64(v.Data.(int))
	case executor.Float64:
		a.floatSum += v.Data.(float64)
	case executor.Float32:
		a.floatSum += float64(v.Data.(float32))
	default:
		return &ErrInvalidArg{Type: v.Type, Function: "avg"}
	}

	return nil
}

func (a *AvgAggregator) Final() executor.Value {
	if a.intSum != 0 {
		return executor.Value{
			Type: executor.Float64,
			Data: float64(a.intSum) / float64(a.count),
		}
	}

	if a.floatSum == 0.0 {
		return executor.Value{
			Type: executor.Float64,
			Data: 0.0,
		}
	}

	return executor.Value{
		Type: executor.Float64,
		Data: a.floatSum / float64(a.count),
	}
}

type ErrInvalidArg struct {
	Type     executor.ValueType
	Function string
}

func (e *ErrInvalidArg) Error() string {
	return fmt.Sprintf("invalid type %s for function %s", e.Type.String(), e.Function)
}
