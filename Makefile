.PHONY: test coverage coverage-html proto certs server client build

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

test:
	go test ./... -coverprofile=cover

coverage: test
	go tool cover -func=cover

coverage-html: test
	go tool cover -html=cover

proto: $(PROTO_FILES)
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

certs:
	cd certs && ./generate.sh > /dev/null
.PHONY: certs

server:
	go build -o builds/$@ cmd/$@/*.go
.PHONY: server

client:
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

build: client server
.PHONY: build