##@ Tooling

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GORELEASER ?= $(LOCALBIN)/goreleaser

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.52.2
GORELEASER_VERSION ?= v1.18.2

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: goreleaser
goreleaser: $(GORELEASER) ## Download goreleaser locally if necessary.
$(GORELEASER): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION)