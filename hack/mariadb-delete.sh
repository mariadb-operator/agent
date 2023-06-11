#!/bin/bash

set -e

kubectl delete configmap scripts
kubectl delete -f hack/manifests/services.yaml
kubectl delete -f hack/manifests/statefulset.yaml

rm -rf mariadb/config/*
rm -rf mariadb/state/*