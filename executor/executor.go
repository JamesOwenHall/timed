package executor

import (
	"fmt"
)

// Executor runs aggregations over a data source.
type Executor struct {
	Source Iterator
	Calls  []AggregatorCall
}

// Execute runs the query and returns the results of the aggregations in a
// record.
func (e *Executor) Execute() (Record, error) {
	for {
		rec, err := e.Source.Next()
		if err != nil {
			return nil, err
		}
		if rec == nil {
			break
		}

		for _, call := range e.Calls {
			val, exists := rec[call.Key]
			if !exists {
				continue
			}

			if err := call.Aggregator.Next(val); err != nil {
				return nil, err
			}
		}
	}

	res := make(Record)
	ng := make(nameGenerator)
	for _, call := range e.Calls {
		if call.Alias != "" {
			res[call.Alias] = call.Aggregator.Final()
		} else {
			name := ng.name(call)
			res[name] = call.Aggregator.Final()
		}
	}

	return res, nil
}

// nameGenerator generates unique names from aggregate calls.
type nameGenerator map[string]bool

// name returns a new unique name.
func (n nameGenerator) name(call AggregatorCall) string {
	base := fmt.Sprintf("%s(%s)", call.Aggregator.Name(), call.Key)

	if !n[base] {
		n[base] = true
		return base
	}

	for i := 2; ; i++ {
		name := fmt.Sprintf("%s - %d", base, i)
		if !n[name] {
			n[name] = true
			return name
		}
	}
}
