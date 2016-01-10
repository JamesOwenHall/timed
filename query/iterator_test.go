package query

type RecordIterator struct {
	records []Record
	index   int
}

func NewRecordIterator(recs []Record) *RecordIterator {
	return &RecordIterator{records: recs}
}

func (r *RecordIterator) Next() *Record {
	if r.index >= len(r.records) {
		return nil
	}

	rec := &r.records[r.index]
	r.index++
	return rec
}

type CountAggregator struct {
	count int
}

func (c *CountAggregator) Name() string {
	return "count"
}

func (c *CountAggregator) Next(v Value) {
	c.count++
}

func (c *CountAggregator) Final() Value {
	return Value{
		Type: Number,
		Data: float64(c.count),
	}
}
