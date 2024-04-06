all: docker-app
test: docker-build docker-e2e-test docker-lint

ARCH=amd64
TAG=local

.PHONY: docker-app
docker-app:
	@docker build . --target app --platform linux/${ARCH} -t bivi/backend:${TAG}

.PHONY: docker-build
docker-build:
	@docker build . --target build --platform linux/${ARCH} -t bivi/backend:build

.PHONY: docker-e2e-test
docker-e2e-test:
	@docker build . --target e2e-test --platform linux/${ARCH} -t bivi/backend:e2e-test
	@docker create --name extract bivi/backend:e2e-test
	@docker cp extract:/test-reports .
	@docker rm extract

.PHONY: docker-lint
docker-lint:
	@docker build . --target lint --platform linux/${ARCH} -t bivi/backend:lint
