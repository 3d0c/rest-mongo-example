// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vektra/mockery/v2"
)

//go:generate go build -v -o=./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
