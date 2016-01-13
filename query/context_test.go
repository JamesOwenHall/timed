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

func TestExecuteQuery(t *testing.T) {
	if session == nil {
		t.Skip("Cassandra is not initialized.")
	}

	q := &Query{
		Source: "data",
		Start:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
		Calls: []FunctionCall{
			FunctionCall{
				Function: "count",
				Argument: "shop_id",
			},
			FunctionCall{
				Function: "count",
				Argument: "fake_field",
			},
		},
	}

	context := Context{
		Sources: []cassandra.Source{
			cassandra.Source{
				Session:     session,
				Consistency: gocql.One,
				Name:        "data",
				TimeKey:     "period",
			},
		},
	}

	expected := executor.Record{
		"count(shop_id)":    executor.Value{executor.Int64, int64(4)},
		"count(fake_field)": executor.Value{executor.Int64, int64(0)},
	}

	actual, err := context.ExecuteQuery(q)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}
