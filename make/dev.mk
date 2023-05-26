##@ Dev

.PHONY: lint
lint: golangci-lint ## Lint.
	$(GOLANGCI_LINT) run

.PHONY: build
build: ## Build binary.
	go build -o bin/agent main.go

.PHONY: test
test: ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: cover
cover: test ## Run tests and generate coverage.
	@go tool cover -html=cover.out -o=cover.html

.PHONY: release
release: goreleaser ## Test release locally.
	$(GORELEASER) release --snapshot --rm-dist

RUN_FLAGS ?=
.PHONY: run
run: lint ## Run a controller from your host.
	go run main.go $(RUN_FLAGS)