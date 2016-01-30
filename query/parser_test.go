package query

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	input := "COMPUTE a(b), c(d), e(f) FROM t1 SINCE 2000-01-01T00:00:00Z UNTIL 2000-01-10T00:00:00-05:00"
	actual, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	until, _ := time.Parse(time.RFC3339, "2000-01-10T00:00:00-05:00")
	expected := &Query{
		Source: "t1",
		Since:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:  until,
		Calls: []FunctionCall{
			FunctionCall{Function: "a", Argument: "b"},
			FunctionCall{Function: "c", Argument: "d"},
			FunctionCall{Function: "e", Argument: "f"},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}
