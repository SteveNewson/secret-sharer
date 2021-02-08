# Project Parameters
VERSION?=`git rev-parse --short HEAD`
BUILD?=`date +%FT%T%z`
ID?=`id -u`
GO_CMD?="go"
GO_FMT_CMD?="gofmt"
ROOT_DIR?=$(CURDIR)
PLATFORM_ROOT?=$(CURDIR)/../..
IN_CI?=""
GITHUB_TOKEN?=

# Configuration
LDFLAGS=-ldflags "-w -s -X github.com/SteveNewson/secret-sharer/internal.Version=${VERSION} -X github.com/SteveNewson/secret-sharer/internal.Build=${BUILD} -X main.Version=${VERSION} -X main.Build=${BUILD}"

# Shortcuts
BIN_PATH=$(ROOT_DIR)/bin
TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')

default: test build

clean:
	$(GO_CMD) clean
	rm -f $(BIN_PATH)/secret-sharer

test: fmtcheck
	$(GO_CMD) test $(TEST)

build: fmtcheck
	CGO_ENABLED=0 GOOS=linux $(GO_CMD) build -installsuffix 'static' $(LDFLAGS) -o ./bin/secret-sharer -v ./cmd/secret-sharer

fmt:
	$(GO_FMT_CMD) -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "GO_FMT_CMD='$(GO_FMT_CMD)' '$(CURDIR)/scripts/gofmtcheck.sh'"

install: build
	@echo "==> Installing to $(GOPATH)/bin ..."
	go install $(LDFLAGS) -v $(ROOT_DIR)/cmd/secret-sharer

.PHONY: all build clean default fmt fmtcheck install test
