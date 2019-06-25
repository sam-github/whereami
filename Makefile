default: test

test:
	go test -v

bench:
	go test -v -bench=.

build:
	go build .
