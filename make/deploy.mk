##@ Docker

PLATFORM ?= linux/amd64,linux/arm64
IMG ?= ghcr.io/mariadb-operator/agent:latest
BUILDX ?= docker buildx build --platform $(PLATFORM) -t $(IMG) 
BUILDER ?= agent

.PHONY: docker-builder
docker-builder: ## Configure docker builder.
	docker buildx create --name $(BUILDER) --use --platform $(PLATFORM)

.PHONY: docker-build
docker-build: ## Build docker image.
	docker build -t $(IMG) .  

.PHONY: docker-buildx
docker-buildx: ## Build multi-arch docker image.
	$(BUILDX) .

.PHONY: docker-push
docker-push: ## Build multi-arch docker image and push it to the registry.
	$(BUILDX) --push .

.PHONY: docker-inspect
docker-inspect: ## Inspect docker image.
	docker buildx imagetools inspect $(IMG)

##@ MariaDB

.PHONY: mariadb
mariadb: ## Create a MariaDB galera cluster using docker compose.
	docker compose up -d

.PHONY: mariadb-delete
mariadb-delete: ## Delete the MariaDB galera cluster.
	docker compose rm --stop --force
	sudo rm -rf mariadb

.PHONY: mariadb-logs
mariadb-logs: ## Check the MariaDB galera cluster logs.
	docker compose logs --follow

.PHONY: mariadb-ps
mariadb-ps: ## Check the MariaDB processes.
	ps -ef | grep mariadbd