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
    timekey: field1
    consistency: one
  - name: table2
    timekey: field_foo
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
			Timekey     string
			Consistency string
		}{
			{Name: "table1", Timekey: "field1", Consistency: "one"},
			{Name: "table2", Timekey: "field_foo", Consistency: ""},
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
