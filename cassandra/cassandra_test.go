package cassandra

import (
	"testing"
)

func TestSourceMakeQuery(t *testing.T) {
	source := Source{
		Name:    "foo",
		TimeKey: "tk",
	}
	expected := "SELECT * FROM foo WHERE tk >= ? AND tk < ?"
	actual := source.makeQuery()

	if actual != expected {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}
}
