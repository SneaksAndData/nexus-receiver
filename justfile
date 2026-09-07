# Default recipe: list available recipes
default:
    @just --list

# Build application binaries
build:
    go build -v ./...

# Start docker compose dependencies
compose-up:
    docker compose up --quiet-pull -d

# Stop and remove docker compose dependencies
compose-down:
    docker compose down

# Build application binaries and start docker compose dependencies
up: build compose-up

# Stop docker compose dependencies
stop: compose-down

# Restart fresh: stop existing dependencies, rebuild, and start dependencies
fresh: stop up

# Run unit tests
test:
    APPLICATION_ENVIRONMENT=units go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

