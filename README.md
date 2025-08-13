![coverage](https://raw.githubusercontent.com/SneaksAndData/nexus-receiver/badges/.badges/main/coverage.svg)

# Nexus Receiver
Nexus Receiver is an essential component of a Nexus deployment, responsible for result accounting. It is deployed in the same network as algorithm container and provides decoupling between scheduling and result submission API.
Receiver can be deployed anywhere, it only needs outbound access to Nexus checkpoint store host.

## Quickstart

Receiver requires a connection to Apache Cassandra or similar backend that is used by Nexus schedulers. Example is provided in the [helm values](.helm/values.yaml). Please refer to [Nexus QuickStart](https://github.com/SneaksAndData/nexus?tab=readme-ov-file#quickstart) for additional information. Once you have a secret with Cassandra connection details, install the receiver:
```shell
helm install nexus-receiver --namespace nexus --create-namespace oci://ghcr.io/sneaksanddata/helm/nexus-receiver \
--set receiver.config.cqlStore.secretName=nexus-cassandra \
--set ginMode=release
```

### API Management
Adding new API paths must be reflected in Swagger docs, even though the app doesn't serve Swagger. Update the generated docs:
```shell
./swag init --parseDependency --parseInternal -g main.go
```

This is required for the API clients (Go and Python) to be updated correctly. Note that until Swag 2.0 is released OpenAPI v3 model must be updated using [Swagger converter](https://converter.swagger.io/#/Converter/convertByContent)
