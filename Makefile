install:
	@go mod vendor

test:
	@go test $(shell go list ./... | grep -v '/vendor/') --cover

package:
	@go clean
	@OS="darwin"
	@CGO_ENABLED=0 GOOS=$$OS go build ./cmd/twothy/

all: install test package
