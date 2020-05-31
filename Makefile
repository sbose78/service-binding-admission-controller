.DEFAULT_GOAL := build

IMAGE ?= quay.io/shbose/service-binding-admission-controller:local

.PHONY: build
bin/service-binding-webhook-server: $(shell find . -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./cmd/service-binding-admission-controller

.PHONY: build-image
build-image: build
	docker build -t $(IMAGE) --file Dockerfile .

.PHONY: push
push: build
	docker push $(IMAGE)