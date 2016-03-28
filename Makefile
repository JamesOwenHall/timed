export TEST_CASSANDRA=127.0.0.1

all: bin/timed-server

clean:
	rm -rf bin

test:
	go test -v ./...

bin:
	mkdir -p bin

bin/timed-server: bin cmd/timed-server/*.go cassandra/*.go executor/*.go query/*.go server/*.go
	go build -o bin/timed-server cmd/timed-server/main.go

fmt:
	go fmt ./...

run:
	go run cmd/timed-server/main.go -c timed_default.yml

.PHONY: clean test fmt run
