SHELL=/bin/bash


COMMIT_VERSION:=$(shell bash commit-version.bash)

help:
	@echo "$(COMMIT_VERSION)"

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/$(COMMIT_VERSION)/run-log-$(GOOS)-$(GOARCH) main.go

build-all:
	make build GOOS=linux GOARCH=amd64
	make build GOOS=linux GOARCH=arm64
	make build GOOS=linux GOARCH=mips64le
	make build GOOS=darwin GOARCH=amd64
	make build GOOS=darwin GOARCH=arm64
