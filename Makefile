.PHONY: test coverage coverage-html proto certs server client build help

PROTO_SRC = proto
PROTO_FILES = notification secrets users
PROTO_DST = pkg/$(PROTO_SRC)

BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date +%d.%m.%y)
BUILD_VERSION ?= 0.0.2

PLATFORMS = \
    darwin/amd64 \
    darwin/arm64 \
    linux/amd64 \
    windows/amd64

help: ## Display help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

test: ## Running tests with generation of coverage report
	go test ./... -coverprofile=cover

coverage: test ## Generate test coverage report in text format
	go tool cover -func=cover

coverage-html: test ## Generate test coverage report in HTML format
	go tool cover -html=cover

proto: $(PROTO_FILES) ## Generating code from protobuf files
.PHONY: proto

$(PROTO_FILES): %: $(PROTO_DST)/%

$(PROTO_DST)/%:
	protoc \
		--proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_DST) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DST) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/$(notdir $@).proto

certs: ## Generating certificates for secure connection
	cd certs && ./generate.sh > /dev/null
.PHONY: certs

server: ## Building a server application
	go build -o builds/$@ cmd/$@/*.go
.PHONY: server

client: ## Building a client application for all platforms
	go build \
		-ldflags "\
			-X 'main.buildVersion=$(BUILD_VERSION)' \
			-X 'main.buildDate=$(BUILD_DATE)' \
			-X 'main.buildCommit=$(BUILD_COMMIT)' \
		" \
		-o cmd/$@/$@ \
		cmd/$@/*.go

	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT=builds/client-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build \
			-ldflags "\
				-X 'main.buildVersion=$(BUILD_VERSION)' \
				-X 'main.buildDate=$(BUILD_DATE)' \
				-X 'main.buildCommit=$(BUILD_COMMIT)' \
			" \
			-o $$OUTPUT \
			cmd/client/*.go || exit 1; \
	done
	rm cmd/client/client

.PHONY: client

build: client server ## Building client and server applications
.PHONY: build