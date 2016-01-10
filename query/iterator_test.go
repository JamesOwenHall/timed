package query

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

func (c *CountAggregator) Next(v Value) error {
	c.count++
	return nil
}

func (c *CountAggregator) Final() Value {
	return Value{
		Type: Number,
		Data: float64(c.count),
	}
}
