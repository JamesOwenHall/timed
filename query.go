package timed

import (
	"fmt"
)

type Query struct {
	Source Iterator
	Calls  []AggregatorCall
}

func (q *Query) Execute() Record {
	for rec := q.Source.Next(); rec != nil; rec = q.Source.Next() {
		for _, call := range q.Calls {
			val, exists := (*rec)[call.Key]
			if !exists {
				continue
			}

			call.Aggregator.Next(val)
		}
	}

	res := make(Record)
	ng := make(nameGenerator)
	for _, call := range q.Calls {
		name := ng.name(call)
		res[name] = call.Aggregator.Final()
	}

	return res
}

type nameGenerator map[string]bool

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
