#!/bin/bash

set -e

kubectl create configmap scripts --from-file=entrypoint.sh=hack/entrypoint.sh --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f hack/config/services.yaml
kubectl apply -f hack/config/statefulset.yaml

sudo chown -R $(id -u):$(id -g) mariadb