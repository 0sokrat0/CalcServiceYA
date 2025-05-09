GO_OUT_DIR = ./pkg/gen
# PROTO_DIR = ./api/auth
PROTO_DIR = ./api/task
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: proto
proto:
	@echo "Generating code from proto files in $(PROTO_DIR)..."
	protoc --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)

.PHONY: build
build:
	@echo "Building binary..."
	go build -o bin/grpc-auth ./cmd

.PHONY: run
run: build
	@echo "Starting server..."
	./bin/grpc-auth

.PHONY: clean
clean:
	@echo "Cleaning binary..."
	rm -f bin/grpc-auth

docker:
	docker build  -t grpc-auth .