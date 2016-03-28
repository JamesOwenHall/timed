package cassandra

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/JamesOwenHall/timed/executor"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	if cass := os.Getenv("TEST_CASSANDRA"); cass != "" {
		cluster := gocql.NewCluster(cass)
		cluster.Keyspace = "timed"
		cluster.ProtoVersion = 4

		var err error
		session, err = cluster.CreateSession()
		if err != nil {
			panic(err)
		}
	}
}

func TestIterator(t *testing.T) {
	if session == nil {
		t.Skip("Cassandra is not initialized.")
	}

	source := Source{
		Session:     session,
		Consistency: gocql.One,
		Name:        "data",
		TimeKey:     "period",
	}
	iter := source.Iterator(
		time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
	)

	expected := []executor.Record{
		{
			"partition": 1,
			"period":    time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC),
			"shop_id":   int64(100),
			"sample":    "hola",
		},
		{
			"partition": 1,
			"period":    time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC),
			"shop_id":   int64(100),
			"sample":    "hola",
		},
		{
			"partition": 1,
			"period":    time.Date(2015, 1, 4, 0, 0, 0, 0, time.UTC),
			"shop_id":   int64(100),
			"sample":    "hola",
		},
	}

	actual := []executor.Record{}
	for {
		rec, err := iter.Next()
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		if rec == nil {
			break
		}

		actual = append(actual, rec)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}

func TestSourceMakeQuery(t *testing.T) {
	source := Source{
		Name:    "foo",
		TimeKey: "tk",
	}
	expected := "SELECT * FROM foo WHERE partition = 1 AND tk >= ? AND tk < ?"
	actual := source.makeQuery()

	if actual != expected {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}
}
