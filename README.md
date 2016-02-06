# timed

Use Apache Cassandra as a time-series database and run aggregation queries over time ranges.

timed is a query execution layer that leverages Cassandra as a datastore.  Queries are written in a SQL-like query language that allows you to express the aggregation functions to run.

```
COMPUTE
  sum(subtotal),
  avg(item_count)
FROM
  orders
SINCE
  '2000-01-01T00:00:00Z'
UNTIL
  '2000-03-01T00:00:00Z'
```

### Setup

timed is written in Go and requires the Go compiler to build, as well as Make.

```sh
$ make
```

This command builds the binary in the `bin/` directory.  You'll also need to create a configuration file.  You can use `timed_default.yml` as a starting place.

### Usage

```sh
$ bin/timed-server -c timed_default.yml
```

The previous command starts the timed server using `timed_default.yml` as the configuration file.  You can then query at [http://localhost:8565/query](http://localhost:8565/query).

### Running the tests

Running the usual `go test ./...` will work, but it will skip some of the tests.  This is because some of the tests rely on a running instance of Cassandra.  You can see which ones by adding the `-v` flag.

To run the tests that rely on Cassandra, set the `TEST_CASSANDRA` environment variable to the Cassandra's IP address.
