package query

import (
	"fmt"
	"time"

	"github.com/JamesOwenHall/timed/cassandra"
	"github.com/JamesOwenHall/timed/executor"
)

type Context struct {
	Sources []cassandra.Source
}

func (c *Context) ExecuteQuery(q *Query) (executor.Record, error) {
	ex, err := c.ToExecutor(q)
	if err != nil {
		return nil, err
	}

	return ex.Execute()
}

func (c *Context) ToExecutor(q *Query) (*executor.Executor, error) {
	// Source.
	var source *cassandra.Source
	for i := range c.Sources {
		if c.Sources[i].Name == q.Source {
			source = &c.Sources[i]
			break
		}
	}
	if source == nil {
		return nil, &ErrInvalidQuery{
			Message:   "unknown source",
			Component: "source",
			Value:     q.Source,
		}
	}

	// Time range.
	if !q.Start.Before(q.End) {
		return nil, &ErrInvalidQuery{
			Message:   "invalid range",
			Component: "end",
			Value:     q.End.Format(time.RFC3339),
		}
	}

	// Function calls.
	aggregators := make([]executor.AggregatorCall, 0)
	for _, fc := range q.Calls {
		agg := NewAggregator(fc.Function)
		if agg == nil {
			return nil, &ErrInvalidQuery{
				Message:   "unknown function",
				Component: "compute",
				Value:     fc.Function,
			}
		}

		aggregators = append(aggregators, executor.AggregatorCall{
			Key:        fc.Argument,
			Aggregator: agg,
		})
	}

	return &executor.Executor{
		Source: source.Iterator(q.Start, q.End),
		Calls:  aggregators,
	}, nil
}

type ErrInvalidQuery struct {
	Message   string
	Component string
	Value     string
}

func (e *ErrInvalidQuery) Error() string {
	return fmt.Sprintf(
		"invalid query => %s component: %s, value: %s",
		e.Message, e.Component, e.Value,
	)
}
