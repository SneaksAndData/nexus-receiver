# Nexus Receiver
Nexus Receiver is an essential component of a Nexus deployment, responsible for result accounting. It is deployed in the same network as algorithm container and provides decoupling between scheduling and client API.
Receiver can be deployed anywhere, it only needs outbound access to Nexus checkpoint store host.

## Quickstart

-- TBD --

### API Management
Adding new API paths must be reflected in Swagger docs, even though the app doesn't serve Swagger. Update the generated docs:
```shell
./swag init --parseDependency --parseInternal -g main.go
```

This is required for the API clients (Go and Python) to be updated correctly. Note that until Swag 2.0 is released OpenAPI v3 model must be updated using [Swagger converter](https://converter.swagger.io/#/Converter/convertByContent)