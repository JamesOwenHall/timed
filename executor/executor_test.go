package executor

import (
	"reflect"
	"testing"
)

func TestExecutorExecute(t *testing.T) {
	q := Executor{
		Source: NewRecordIterator([]Record{
			{"foo": Value{Number, 1}, "bar": Value{Null, nil}},
			{"foo": Value{Number, 1}, "bar": Value{Null, nil}},
			{"foo": Value{Number, 1}, "bar": Value{Null, nil}},
			{"foo": Value{Number, 1}, "bar": Value{Null, nil}},
			{"foo": Value{Number, 1}, "bar": Value{Null, nil}},
			{"foo": Value{Number, 1}},
			{"foo": Value{Number, 1}},
			{"foo": Value{Number, 1}},
			{"foo": Value{Number, 1}},
			{"foo": Value{Number, 1}},
		}),
		Calls: []AggregatorCall{
			{Key: "foo", Aggregator: new(CountAggregator)},
			{Key: "bar", Aggregator: new(CountAggregator)},
			{Key: "baz", Aggregator: new(CountAggregator), Alias: "dummy"},
		},
	}

	expected := Record{
		"count(foo)": Value{Number, 10.0},
		"count(bar)": Value{Number, 5.0},
		"dummy":      Value{Number, 0.0},
	}
	actual, err := q.Execute()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}

func TestExecutorExecuteError(t *testing.T) {
	q := Executor{
		Source: NewRecordIterator([]Record{
			{"foo": Value{Number, 1}},
		}),
		Calls: []AggregatorCall{
			{Key: "foo", Aggregator: new(CountAggregator)},
			{Key: "foo", Aggregator: new(ErrorAggregator)},
		},
	}

	_, err := q.Execute()
	if err != ErrErrorAggregator {
		t.Fatalf("\nExpected: %s\n     Got: %v", ErrErrorAggregator, err)
	}
}

func TestNameGenerator(t *testing.T) {
	call := AggregatorCall{Key: "foo", Aggregator: new(CountAggregator)}
	ng := make(nameGenerator)

	expected := []string{
		"count(foo)",
		"count(foo) - 2",
		"count(foo) - 3",
		"count(foo) - 4",
	}

	for _, exp := range expected {
		actual := ng.name(call)
		if actual != exp {
			t.Fatalf("\nExpected: %s\n     Got: %s", exp, actual)
		}
	}
}