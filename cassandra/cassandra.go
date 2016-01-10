package cassandra

import (
	"fmt"
	"time"

	"github.com/JamesOwenHall/timed/executor"

	"github.com/gocql/gocql"
)

type Source struct {
	Session     *gocql.Session
	Consistency gocql.Consistency
	Name        string
	TimeKey     string
}

func (s *Source) Iterator(start, end time.Time) executor.Iterator {
	return &iterator{
		iter: s.Session.Query(s.makeQuery(), start, end).
			Consistency(s.Consistency).Iter(),
	}
}

func (s *Source) makeQuery() string {
	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s >= ? AND %s < ? ALLOW FILTERING",
		s.Name, s.TimeKey, s.TimeKey,
	)
}

type iterator struct {
	iter *gocql.Iter
}

func (i *iterator) Next() (executor.Record, error) {
	m := make(map[string]interface{})
	if ok := i.iter.MapScan(m); !ok {
		return nil, i.iter.Close()
	}

	rec := make(executor.Record)
	for k, v := range m {
		switch v := v.(type) {
		case string:
			rec[k] = executor.Value{
				Type: executor.String,
				Data: v,
			}
		case bool:
			rec[k] = executor.Value{
				Type: executor.Boolean,
				Data: v,
			}
		default:
			rec[k] = executor.Value{
				Type: executor.Null,
				Data: nil,
			}
		}
	}

	return rec, nil
}
