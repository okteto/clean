COMMIT_SHA := $(shell git rev-parse --short HEAD)
VERSION := '0.1.2-dev'

.PHONY: push
push:
	okteto build -t okteto/clean:${VERSION} --build-arg COMMIT=${COMMIT_SHA} .

multi:
	# docker buildx create --name mbuilder
	docker buildx use mbuilder
	docker buildx build  --platform linux/amd64,linux/arm64 -t okteto/clean:${VERSION} --build-arg COMMIT=${COMMIT_SHA}  --push .