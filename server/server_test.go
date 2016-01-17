package server

import (
	"os"
	"testing"

	"github.com/gocql/gocql"
)

var cassandraAddr string

func init() {
	cassandraAddr = os.Getenv("TEST_CASSANDRA")
}

func TestNewServer(t *testing.T) {
	if cassandraAddr == "" {
		t.Skip("Cassandra is not initialized.")
	}

	config, err := NewConfigFromYAML([]byte(`
listen: ":1234"
cassandra:
  keyspace: timed
  addresses:
    - ` + cassandraAddr + `
sources:
  - name: table1
    timekey: time
    consistency: one
`))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	if server.server.Addr != ":1234" {
		t.Errorf("listen: %s", server.server.Addr)
	}
	if len(server.sources) != 1 {
		t.Fatalf("sources: %s", server.sources)
	}
	if server.sources[0].Name != "table1" {
		t.Errorf("source name: %s", server.sources[0].Name)
	}
	if server.sources[0].TimeKey != "time" {
		t.Errorf("source timekey: %s", server.sources[0].TimeKey)
	}
	if server.sources[0].Consistency != gocql.One {
		t.Errorf("source consistency: %s", server.sources[0].Consistency)
	}
}
