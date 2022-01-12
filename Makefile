include variables.mk

.PHONY: build

.DEFAULT_GOAL=build

version:
	echo $(FULL_VERSION)

build:
	${GO_ENV_VARS} go build ${BUILD_FLAGS} -o build/${BINARY} .

build/${DARWIN}:
	${GO_ENV_VARS} GOOS=${DARWIN} go build ${BUILD_FLAGS} -o build/${BINARY}-${DARWIN} .

build/${LINUX}:
	${GO_ENV_VARS} GOOS=${LINUX} go build ${BUILD_FLAGS} -o build/${BINARY}-${LINUX} .

run-uts:
	${GO_ENV_VARS} go test `go list ./... | grep -vE "mocks|/cmd|/test"` -race -cover -coverprofile=$(COVERAGE_OUT)

cover:
	# test coverage report to the console
	go tool cover -func=$(COVERAGE_OUT)
	# html file with code coverage report
	go tool cover -html=$(COVERAGE_OUT) -o coverage.html

lint:
ifeq (,$(wildcard ./bin/golangci-lint))
	@echo ">>> installing golangci-lint into bin directory"
	@go build -v -o=./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
endif
	$(GO_ENV_VARS) bin/golangci-lint -v run ./...

tools:
	go generate ./tools.go

clean:
	@rm -rf build/*
