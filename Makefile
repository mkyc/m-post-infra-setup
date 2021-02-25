ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION ?= dev
IMAGE_REPOSITORY := epiphanyplatform/pis

IMAGE_NAME := $(IMAGE_REPOSITORY):$(VERSION)

export

#used for correctly setting shared folder permissions
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

.PHONY: all

all: build

.PHONY: build test pipeline-test release prepare-service-principal

build: guard-IMAGE_NAME
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		--build-arg ARG_HOST_UID=$(HOST_UID) \
		--build-arg ARG_HOST_GID=$(HOST_GID) \
		-t $(IMAGE_NAME) \
		.

print-%:
	@echo "$($*)"
	

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

doctor:
	go mod tidy
	go fmt ./...
	go vet ./...
	goimports -l -w .
