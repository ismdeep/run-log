SHELL=/bin/bash


COMMIT_VERSION:=$(shell bash commit-version.bash)

help:
	@echo "$(COMMIT_VERSION)"

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/$(COMMIT_VERSION)/run-log-$(GOOS)-$(GOARCH) main.go

build-all:
	make build GOOS=linux GOARCH=386
	make build GOOS=linux GOARCH=amd64
	make build GOOS=linux GOARCH=arm64
	make build GOOS=linux GOARCH=arm
	make build GOOS=linux GOARCH=mips64le
	make build GOOS=darwin GOARCH=amd64
	make build GOOS=darwin GOARCH=arm64
