#!/usr/bin/env bash

set -euxo pipefail

go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
go install github.com/kisielk/errcheck@latest
go install github.com/gordonklaus/ineffassign@latest
# I got Error: ../../../go/pkg/mod/golang.org/x/tools@v0.20.0/go/types/objectpath/objectpath.go:397:10: meth.Origin undefined (type *types.Func has no field or method Origin)
# go install github.com/opennota/check/cmd/varcheck@latest
# go install github.com/opennota/check/cmd/aligncheck@latest
# Comment out because of the error: internal error: package "fmt" without types was imported from
# go install github.com/mdempsky/unconvert@latest
