SHELL=/bin/bash

COMMIT_DATE=$(shell git log -1 --pretty=format:"%ci" | awk '{print $$1}' | sed 's/-//g')
COMMIT_TIME=$(shell git log -1 --pretty=format:"%ci" | awk '{print $$2}' | sed 's/://g')
COMMIT_DATETIME_H=$(shell git log -1 --pretty=format:"%ci")
COMMIT_ID=$(shell git log -1 --pretty=format:"%h" | awk '{print $$1}')
COMMIT_ID_LONG=$(shell git log -1 --pretty=format:"%H" | awk '{print $$1}')
COMMIT_VERSION=${COMMIT_DATE}$(COMMIT_TIME)-$(COMMIT_ID)

help:
	@echo $(COMMIT_VERSION)

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.goVersion=$(shell go version)' -X 'main.gitHash=$(COMMIT_ID_LONG)' -X 'main.gitCommitDateTime=$(COMMIT_DATETIME_H)' -X 'main.version=$(COMMIT_VERSION)'" \
		-o ./bin/$(COMMIT_VERSION)/run-log-$(GOOS)-$(GOARCH) main.go

build-all:
	make build GOOS=linux GOARCH=386
	make build GOOS=linux GOARCH=amd64
	make build GOOS=linux GOARCH=arm64
	make build GOOS=linux GOARCH=arm
	make build GOOS=linux GOARCH=mips64le
	make build GOOS=darwin GOARCH=amd64
	make build GOOS=darwin GOARCH=arm64
