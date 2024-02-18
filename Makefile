## test/all runs all related tests.
test/all: test/lint test

## test runs `go test`
test:
	bazel test //...

test/ci:
	bazel --bazelrc=.github/workflows/ci.bazelrc --bazelrc=.bazelrc test //...

## test runs `go test -run $(RUN)`
test/run:
	bazel test --test_filter="$(RUN)" //...

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
	go vet -vettool=$(which shadow) ./...
	# checks not to ignore the error.
	errcheck ./...
	# checks unused global variables and constants.
	varcheck ./...
	# checks no used assigned value.
	ineffassign .
	# checks dispensable type conversions.
	## Comment out because of the error: internal error: package "fmt" without types was imported from
	#unconvert -v ./...

## dev/install/dep installs depenencies required for development.
dev/install/dep:
	./.github/workflows/install_dep.sh

## RUN_EXAMPLE_DEBUG is a debug flag argument for run/example.
RUN_EXAMPLE_DEBUG=false

## RUN_EXAMPLE_PERMISSIVE is a permissive flag argument for run/example.
RUN_EXAMPLE_PERMISSIVE=true

## RUN_EXAMPLE_UNORDERED is an unordered flag argument for run/example.
RUN_EXAMPLE_UNORDERED=false

## run/dump/example runs `go run _example/dump/main.go`
run/dump/example:
	bazel run //_example/dump -- -debug=$(RUN_EXAMPLE_DEBUG) -permissive=$(RUN_EXAMPLE_PERMISSIVE) -unordered=${RUN_EXAMPLE_UNORDERED}
