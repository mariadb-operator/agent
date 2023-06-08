#!/bin/bash

set -e

kubectl delete configmap scripts
kubectl delete -f hack/config/services.yaml
kubectl delete -f hack/config/statefulset.yaml

rm -rf mariadb/config/*
rm -rf mariadb/state/*