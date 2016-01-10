package query

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestQueryUnmarshal(t *testing.T) {
	input := `
{
	"source": "foo",
	"start": "2015-01-01T00:00:00Z",
	"end": "2015-01-10T00:00:00Z",
	"compute": [
		{"function": "count", "argument": "foo"},
		{"function": "min", "argument": "bar"}
	]
}`
	expected := Query{
		Source: "foo",
		Start:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2015, 1, 10, 0, 0, 0, 0, time.UTC),
		Calls: []FunctionCall{
			{Function: "count", Argument: "foo"},
			{Function: "min", Argument: "bar"},
		},
	}

	var actual Query
	if err := json.Unmarshal([]byte(input), &actual); err != nil {
		t.Fatalf("Unexpected error: %e", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}
