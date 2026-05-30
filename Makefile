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

get-proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

get-gorm:
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres

build-server:
	go build ./$(SERVER_DIR)

build-client:
	go build ./$(CLIENT_DIR)
clean-server:
	rm -rf ./server

docker-compose-up:
	docker compose up -d

all: get-proto get-gorm proto build-server build-client docker-compose-up

.PHONY: proto clean-proto install-proto build-server build-client clean-server get-proto get-gorm docker-compose-up
