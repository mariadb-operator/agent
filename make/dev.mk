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
	go tool cover -html=cover.out -o=cover.html

.PHONY: release
release: goreleaser ## Test release locally.
	$(GORELEASER) release --snapshot --rm-dist

CONFIG_DIR ?= mariadb/config/local
CONFIG_FILE ?= 1-bootstrap.cnf
.PHONY: config
config: ## Copies a example config file for development purposes.
	@mkdir -p $(CONFIG_DIR)
	cp "examples/$(CONFIG_FILE)" $(CONFIG_DIR)

STATE_DIR ?= mariadb/state/local
STATE_FILE ?= grastate-recovery.dat
.PHONY: state
state: ## Copies a example state file for development purposes.
	@mkdir -p $(STATE_DIR)
	cp "examples/$(STATE_FILE)" "$(STATE_DIR)/grastate.dat"

BASE_RUN_FLAGS ?= --log-level=debug --log-dev

RUN_FLAGS ?= $(BASE_RUN_FLAGS) --config-dir=$(CONFIG_DIR) --state-dir=$(STATE_DIR)
.PHONY: run
run: config state ## Run agent from your host.
	go run main.go $(RUN_FLAGS)

AGENT ?= $(LOCALBIN)/agent
SU_RUN_FLAGS ?= $(BASE_RUN_FLAGS) --config-dir=mariadb/config/docker --state-dir=mariadb/state/docker 
.PHONY: su-run
su-run: build ## Run agent from your host as root to be able to access mariadb volumes and process.
	sudo $(AGENT) $(SU_RUN_FLAGS)