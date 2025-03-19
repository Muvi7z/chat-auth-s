include .env

LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	golangci-lint run ./... --config .golangci.pipeline.yaml

LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/bufbuild/protoc-gen-validate

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make gen-proto

PROTO_DIR=api/user_v1
OUT_DIR=gen/api/user_v1

gen-proto:
	protoc --proto_path api/user_v1 --proto_path vendor.protogen \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		--validate_out lang=go:$(OUT_DIR) --validate_opt=paths=source_relative \
		--plugin=proto-gen-grpc-gateway=./protoc-gen-grpc-gateway \
		$(PROTO_DIR)/*.proto

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

vendor-proto:
			@if [ ! -d vendor.protogen/google ]; then \
						git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
						mkdir -p  vendor.protogen/google/ &&\
						mv vendor.protogen/googleapis/google/api google/api &&\
						rm -rf vendor.protogen/googleapis ;\
			fi
			@if [ ! -d vendor.protogen/validate ]; then \
            			mkdir -p vendor.protogen/validate &&\
            			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
            			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
            			rm -rf vendor.protogen/protoc-gen-validate ;\
			fi