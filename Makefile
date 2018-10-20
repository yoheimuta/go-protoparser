## test/all runs all related tests.
test/all: test/lint test

## test runs `go test`
test:
	time go test -v -p 2 -count 1 -timeout 240s -race ./...

## test runs `go test -run $(RUN)`
test/run:
	time go test -v -p 2 -count 1 -timeout 240s -race ./... -run $(RUN)

## test/lint runs linter
test/lint:
	# checks the coding style.
	(! gofmt -s -d `find . -name vendor -prune -type f -o -name '*.go'` | grep '^')
	golint -set_exit_status `go list ./...`
	# checks the import format.
	(! goimports -l `find . -name vendor -prune -type f -o -name '*.go'` | grep 'go')
	# checks the error the compiler can't find.
	go vet ./...
	# checks shadowed variables.
	go vet -shadow ./...
	# checks not to ignore the error.
	errcheck ./...
	# checks unused global variables and constants.
	varcheck ./...
	# checks no used assigned value.
	ineffassign .
	# checks dispensable type conversions.
	unconvert -v ./...

## RUN_EXAMPLE_DEBUG is a debug flag argument for run/example.
RUN_EXAMPLE_DEBUG=false

## RUN_EXAMPLE_PERMISSIVE is a permissive flag argument for run/example.
RUN_EXAMPLE_PERMISSIVE=true

## run/dump/example runs `go run _example/dump/main.go`
run/dump/example:
	go run _example/dump/main.go -debug=$(RUN_EXAMPLE_DEBUG) -permissive=$(RUN_EXAMPLE_PERMISSIVE)
