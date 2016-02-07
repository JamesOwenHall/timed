package query

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/JamesOwenHall/timed/cassandra"
	"github.com/JamesOwenHall/timed/executor"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	if cass := os.Getenv("TEST_CASSANDRA"); cass != "" {
		cluster := gocql.NewCluster(cass)
		cluster.Keyspace = "timed"

		var err error
		session, err = cluster.CreateSession()
		if err != nil {
			panic(err)
		}
	}
}

func defaultContext() *Context {
	return &Context{
		Sources: []cassandra.Source{
			cassandra.Source{
				Session:     session,
				Consistency: gocql.One,
				Name:        "data",
				TimeKey:     "period",
			},
		},
	}
}

func TestExecuteQuery(t *testing.T) {
	if session == nil {
		t.Skip("Cassandra is not initialized.")
	}

	context := defaultContext()
	q := &Query{
		Source: "data",
		Since:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:  time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
		Calls: []FunctionCall{
			FunctionCall{
				Function: "count",
				Argument: "shop_id",
			},
		},
	}

	expected := executor.Record{
		"count_shop_id": int64(4),
	}

	actual, err := context.ExecuteQuery(q)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}

func TestUnknownSource(t *testing.T) {
	context := defaultContext()
	q := &Query{
		Source: "FAKE_SOURCE",
	}

	_, err := context.ExecuteQuery(q)
	qerr := err.(*ErrInvalidQuery)
	if qerr.Component != "source" {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

func TestInvalidRange(t *testing.T) {
	context := defaultContext()
	q := &Query{
		Source: "data",
		Since:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	_, err := context.ExecuteQuery(q)
	qerr := err.(*ErrInvalidQuery)
	if qerr.Component != "end" {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

func TestUnknownFunction(t *testing.T) {
	context := defaultContext()
	q := &Query{
		Source: "data",
		Since:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:  time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC),
		Calls: []FunctionCall{
			FunctionCall{
				Function: "FAKE",
				Argument: "period",
			},
		},
	}

	_, err := context.ExecuteQuery(q)
	qerr := err.(*ErrInvalidQuery)
	if qerr.Component != "compute" {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}
