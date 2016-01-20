export TEST_CASSANDRA=192.168.99.100

all: bin/timed-server

test:
	go test -v ./...

bin:
	mkdir -p bin

bin/timed-server: bin cmd/timed-server/*.go cassandra/*.go executor/*.go query/*.go server/*.go
	go build -o bin/timed-server cmd/timed-server/main.go

fmt:
	go fmt ./...

.PHONY: fmt test
