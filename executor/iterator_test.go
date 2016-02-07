package executor

import (
	"fmt"
)

type RecordIterator struct {
	records []Record
	index   int
}

func NewRecordIterator(recs []Record) *RecordIterator {
	return &RecordIterator{records: recs}
}

func (r *RecordIterator) Next() (Record, error) {
	if r.index >= len(r.records) {
		return nil, nil
	}

	rec := r.records[r.index]
	r.index++
	return rec, nil
}

type CountAggregator struct {
	count int
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

var ErrErrorAggregator = fmt.Errorf("error from ErrorAggregator")

type ErrorAggregator struct{}

func (e *ErrorAggregator) Name() string {
	return "error"
}

func (e *ErrorAggregator) Next(v interface{}) error {
	return ErrErrorAggregator
}

func (e *ErrorAggregator) Final() interface{} {
	return nil
}
