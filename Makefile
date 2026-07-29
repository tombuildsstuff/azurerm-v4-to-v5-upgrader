MODULE        := github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader
HAS_GOIMPORTS := $(shell command -v goimports 2>/dev/null)

.PHONY: build install test fmt tools

# build tidies the module, refreshes the (git-ignored) vendor tree, and
# compiles the upgrader binary into ./upgrader.
build:
	go mod tidy
	go mod vendor
	go build -o upgrader .

# install builds, then installs the binary into $(go env GOPATH)/bin.
install: build
	go install .

# test runs the full Go test suite. The opt-in terraform-validate layer
# stays skipped unless AZURERM_ACC=1 is set (it needs terraform + network).
test:
	go test ./...

# fmt formats Go sources (gofmt + goimports) and Terraform fixtures.
fmt:
	go fmt ./...
ifdef HAS_GOIMPORTS
	goimports -w -local $(MODULE) $(shell go list -f '{{.Dir}}' ./...)
else
	@echo ">> goimports not found; run 'make tools' (skipping import grouping)"
endif
	terraform fmt -recursive

# tools installs the developer dependencies the other targets use.
# (terraform is assumed to already be on PATH.)
tools:
	go install golang.org/x/tools/cmd/goimports@latest
