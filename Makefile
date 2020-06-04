COMMIT_SHA := $(shell git rev-parse --short HEAD)
TAG ?= 'clean'

.PHONY: push
push:
	okteto build -t ${TAG} --build-arg COMMIT=${COMMIT_SHA} .

multi:
	# docker buildx create --name mbuilder
	docker buildx use mbuilder
	docker buildx build  --platform linux/amd64,linux/arm64 -t ${TAG} --build-arg COMMIT=${COMMIT_SHA}  --push .