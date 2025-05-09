GO_OUT_DIR = ./pkg/gen
PROTO_DIR = ./api/task
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)


proto:
	@echo "🛠️ Generating Go gRPC + Gateway code for all .proto"
	protoc --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)
