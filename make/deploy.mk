CLUSTER ?= agent
CLUSTER_OPERATOR ?= mdb

##@ Docker

PLATFORM ?= linux/amd64,linux/arm64
IMG ?= ghcr.io/mariadb-operator/agent:v0.0.2
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

.PHONY: docker-load
docker-load: ## Load docker image in KIND.
	kind load docker-image --name ${CLUSTER_OPERATOR} ${IMG}

##@ Cluster

KIND_CONFIG ?= hack/config/kind.yaml
KIND_IMAGE ?= kindest/node:v1.26.0

.PHONY: cluster
cluster: kind ## Create the kind cluster.
	$(KIND) create cluster --name $(CLUSTER) --config $(KIND_CONFIG)

.PHONY: cluster-delete
cluster-delete: kind ## Delete the kind cluster.
	$(KIND) delete cluster --name $(CLUSTER)

.PHONY: cluster-ctx
cluster-ctx: ## Sets cluster context.
	@kubectl config use-context kind-$(CLUSTER)

##@ MariaDB

.PHONY: mariadb
mariadb: ## Create a MariaDB galera in kind.
	@./hack/mariadb.sh

.PHONY: mariadb-delete
mariadb-delete: ## Delete the MariaDB galera cluster.
	@./hack/mariadb-delete.sh
