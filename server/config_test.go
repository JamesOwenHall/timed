package server

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	input := []byte(`
listen: ":8231"
cassandra:
  keyspace: timed
  addresses:
    - "1.2.3.4"
    - "5.6.7.8:5678"
sources:
  - name: table1
    consistency: one
  - name: table2
`)

	expected := &Config{
		Listen: ":8231",
		Cassandra: struct {
			Keyspace  string
			Addresses []string
		}{
			Keyspace:  "timed",
			Addresses: []string{"1.2.3.4", "5.6.7.8:5678"},
		},
		Sources: []struct {
			Name        string
			Consistency string
		}{
			{Name: "table1", Consistency: "one"},
			{Name: "table2", Consistency: ""},
		},
	}

	actual, err := NewConfigFromYAML(input)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}
