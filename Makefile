
.PHONY: install
install:
	go install github.com/dewep-online/devtool@latest

.PHONY: setup
setup:
	devtool setup-lib

.PHONY: lint
lint:
	devtool lint

.PHONY: build
build:
	devtool build --arch=amd64

.PHONY: tests
tests:
	devtool test

.PHONY: pre-commite
pre-commite: setup lint build tests

.PHONY: ci
ci: install setup lint build tests

.PHONY: run
run:
	go run ./cmd/artifactory/main.go --config=./config/config.dev.yaml

docker:
	docker-compose up -d

deb:
	deb-builder build
