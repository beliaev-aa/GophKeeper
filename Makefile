.PHONY: test coverage coverage-html proto certs server

PROTO_SRC = proto
PROTO_FILES = notification secrets users
PROTO_DST = pkg/$(PROTO_SRC)

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