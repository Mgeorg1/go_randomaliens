PROTO_DIR=proto
GEN_DIR=gen/pb
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client

proto:
	rm -rf $(GEN_DIR)/*

	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

clean-proto:
	rm -rf (GEN_DIR)/*

install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build-server:
	go build ./$(SERVER_DIR)

build-client:
	go build ./$(CLIENT_DIR)
clean-server:
	rm -rf ./server
.PHONY: proto clean-proto install-proto build-server build-client clean-server
