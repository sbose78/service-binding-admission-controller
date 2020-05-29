.DEFAULT_GOAL := build

IMAGE ?= quay.io/shbose/service-binding-admission-controller:v0.2

bin/webhook-server: $(shell find . -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./cmd/service-binding-webhook-server

.PHONY: build
build: bin/webhook-server
	docker build -t $(IMAGE) --file image/Dockerfile bin/

.PHONY: push
push: build
	docker push $(IMAGE)