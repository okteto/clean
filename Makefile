COMMIT_SHA := $(shell git rev-parse --short HEAD)
VERSION := '0.1.0'

.PHONY: build
build:
	okteto build -t okteto/clean:${VERSION} --build-arg COMMIT=${COMMIT_SHA} .