package executor

import (
	"reflect"
	"testing"
)

func TestExecutorExecute(t *testing.T) {
	q := Executor{
		Source: NewRecordIterator([]Record{
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
			{"foo": Value{Int, 1}, "bar": Value{Unknown, nil}},
		}),
		Calls: []AggregatorCall{
			{Key: "foo", Aggregator: new(CountAggregator)},
			{Key: "bar", Aggregator: new(CountAggregator)},
		},
	}

	expected := Record{
		"count_foo": Value{Int, 10},
		"count_bar": Value{Int, 10},
	}
	actual, err := q.Execute()

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
	}
}

func TestUnknownField(t *testing.T) {
	q := Executor{
		Source: NewRecordIterator([]Record{
			{"foo": Value{Int, 1}},
			{"foo": Value{Int, 2}},
			{"foo": Value{Int, 3}},
		}),
		Calls: []AggregatorCall{
			{Key: "fake_field", Aggregator: new(CountAggregator)},
		},
	}

	_, err := q.Execute()
	eerr := err.(*ErrExecution)
	if eerr.Component != "fake_field" {
		t.Fatalf("Unexpected error: %s", eerr.Error())
	}
}

func TestExecutorAggregatorError(t *testing.T) {
	q := Executor{
		Source: NewRecordIterator([]Record{
			{"foo": Value{Int, 1}},
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
		"count_foo",
		"count_foo_2",
		"count_foo_3",
		"count_foo_4",
	}

	for _, exp := range expected {
		actual := ng.name(call)
		if actual != exp {
			t.Fatalf("\nExpected: %s\n     Got: %s", exp, actual)
		}
	}
}
