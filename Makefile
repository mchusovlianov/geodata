SHELL := /bin/bash

VERSION := 1.0

# ==============================================================================
# Run all tests
test:
	go test ./...


# ==============================================================================
# Run coverages for business logic and for web api app
business_coverage:
	go test -v -coverpkg=github.com/mchusovlianov/geodata/business/... -coverprofile=business_profile.cov  github.com/mchusovlianov/geodata/business/...
	go tool cover -func business_profile.cov

web_coverage:
	go test -v -coverpkg=github.com/mchusovlianov/geodata/app/services/... -coverprofile=web_profile.cov  github.com/mchusovlianov/geodata/app/services/...
	go tool cover -func web_profile.cov

# ==============================================================================
# Build all app
all: geoimport geoapi

geoimport:
	docker build \
		-f infra/docker/dockerfile.importer \
		-t geoimport-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

	docker image prune -f --filter label=stage=builder

geoapi:
	docker build \
		-f infra/docker/dockerfile.geoapi \
		-t geoapi-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

	docker image prune -f --filter label=stage=builder

# ==============================================================================
# Run importer
import:
	docker-compose -f infra/docker-compose/docker-compose.yaml up -d mysql
	docker-compose -f infra/docker-compose/docker-compose.yaml up geoimport

# ==============================================================================
# Run geoapi
api:
	docker-compose -f infra/docker-compose/docker-compose.yaml up -d mysql
	docker-compose -f infra/docker-compose/docker-compose.yaml up geoapi

# ==============================================================================
# Stop all services
down:
	docker-compose -f infra/docker-compose/docker-compose.yaml down
