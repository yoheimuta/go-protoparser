#!/usr/bin/env bash

set -euxo pipefail

go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
go install github.com/kisielk/errcheck@latest
go install github.com/gordonklaus/ineffassign@latest
go install github.com/opennota/check/cmd/varcheck@latest
go install github.com/opennota/check/cmd/aligncheck@latest
# Comment out because of the error: internal error: package "fmt" without types was imported from
# go install github.com/mdempsky/unconvert@latest