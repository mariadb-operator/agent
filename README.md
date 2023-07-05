<p align="center">
<img src="https://mariadb-operator.github.io/mariadb-operator/assets/mariadb-operator.png" alt="mariadb" width="250"/>
</p>

<p align="center">
<a href="https://github.com/mariadb-operator/agent/actions/workflows/ci.yml"><img src="https://github.com/mariadb-operator/agent/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
<a href="https://github.com/mariadb-operator/agent/actions/workflows/release.yml"><img src="https://github.com/mariadb-operator/agent/actions/workflows/release.yml/badge.svg" alt="Release"></a>
<a href="https://goreportcard.com/report/github.com/mariadb-operator/agent"><img src="https://goreportcard.com/badge/github.com/mariadb-operator/agent" alt="Go Report Card"></a>
<a href="https://pkg.go.dev/github.com/mariadb-operator/agent"><img src="https://pkg.go.dev/badge/github.com/mariadb-operator/agent.svg" alt="Go Reference"></a>
</p>


# 🤖 agent
Sidecar agent for MariaDB that co-operates with [mariadb-operator](https://github.com/mariadb-operator/mariadb-operator). Remotely manage Galera via HTTP instead of configuration `*.cnf` files.
- HTTP API to manage Galera and expose the MariaDB state to `mariadb-operator`
- Query and update Galera state without mounting `/var/lib/mysql/grastate.dat`
- Perform [Galera cluster recovery](https://galeracluster.com/library/documentation/crash-recovery.html) remotely 
- Bootstrap new Galera cluster as a result of the cluster recovery
- Idiomatic Go HTTP client [pkg/client](./pkg/client/)
- Authentication using Kubernetes service accounts via [TokenReview](https://kubernetes.io/docs/reference/kubernetes-api/authentication-resources/token-review-v1/) API


### How to use it

Specify the agent image in the `MariaDB` `spec.galera.agent` field.

```yaml
apiVersion: mariadb.mmontes.io/v1alpha1
kind: MariaDB
metadata:
  name: mariadb-galera
spec:
  ...
  image:
    repository: mariadb
    tag: "10.11.3"
    pullPolicy: IfNotPresent
  port: 3306
  replicas: 3

  galera:
    sst: mariabackup
    replicaThreads: 1

    agent:
      image:
        repository: ghcr.io/mariadb-operator/agent
        tag: "v0.0.2"
        pullPolicy: IfNotPresent
      port: 5555
      gracefulShutdownTimeout: 5s
  ...
```

### HTTP API

You can consume the agent API using the [pkg/client](./pkg/client/). Alternatively, take a look at our Postman collection.

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/9776-cbdc1706-5e01-423a-822a-ed46daff6abd?action=collection%2Ffork&collection-url=entityId%3D9776-cbdc1706-5e01-423a-822a-ed46daff6abd%26entityType%3Dcollection%26workspaceId%3Da184b7e4-b1f7-405e-b9ec-ec62ed36dd27#?env%5BKubernetes%5D=W3sia2V5IjoidXJsIiwidmFsdWUiOiJodHRwOi8vbWFyaWFkYi1nYWxlcmEtMC5tYXJpYWRiLWdhbGVyYS1pbnRlcm5hbC5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsOjU1NTUiLCJlbmFibGVkIjp0cnVlLCJ0eXBlIjoiZGVmYXVsdCIsInNlc3Npb25WYWx1ZSI6Imh0dHA6Ly9tYXJpYWRiLWdhbGVyYS0wLm1hcmlhZGItZ2FsZXJhLWludGVybmFsLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWw6NTU1NSIsInNlc3Npb25JbmRleCI6MH1d)
