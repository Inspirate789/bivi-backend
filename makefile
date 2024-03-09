all: docker-app
test: docker-build docker-integration-test docker-lint

ARCH=amd64
TAG=local

.PHONY: docker-app
docker-app:
	@docker build . --target app --platform linux/${ARCH} -t bivi/backend:${TAG}

.PHONY: docker-build
docker-build:
	@docker build . --target build --platform linux/${ARCH} -t bivi/backend:build

.PHONY: docker-integration-test
docker-integration-test:
	@docker build . --target integration-test --platform linux/${ARCH} -t bivi/backend:integration-test

.PHONY: docker-lint
docker-lint:
	@docker build . --target lint --platform linux/${ARCH} -t bivi/backend:lint
