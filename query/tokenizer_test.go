package query

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	input := `
	COMPUTE
		foo(f1),
		bar(f2)
	FROM
		t1
	SINCE
		2000-01-01T00:00:00Z
	UNTIL
		2000-01-10T00:00:00+01:00
	`

	actual := NewTokenizer(input).Collect()
	expected := []Token{
		{Type: Compute, Raw: "COMPUTE"},
		{Type: Identifier, Raw: "foo"},
		{Type: OpenParen, Raw: "("},
		{Type: Identifier, Raw: "f1"},
		{Type: CloseParen, Raw: ")"},
		{Type: Comma, Raw: ","},
		{Type: Identifier, Raw: "bar"},
		{Type: OpenParen, Raw: "("},
		{Type: Identifier, Raw: "f2"},
		{Type: CloseParen, Raw: ")"},
		{Type: From, Raw: "FROM"},
		{Type: Identifier, Raw: "t1"},
		{Type: Since, Raw: "SINCE"},
		{Type: Timestamp, Raw: "2000-01-01T00:00:00Z"},
		{Type: Until, Raw: "UNTIL"},
		{Type: Timestamp, Raw: "2000-01-10T00:00:00+01:00"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\n     Got: %v", expected, actual)
	}
}
