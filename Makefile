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
	go get github.com/bufbuild/protoc-gen-validate
	go get github.com/envoyproxy/protoc-gen-validate
	go install github.com/envoyproxy/protoc-gen-validate
	go install github.com/rakyll/statik@latest
	go install github.com/bojand/ghz/cmd/ghz@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	mkdir -p pkg/swagger
	make gen-proto
	statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

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
		--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
		$(PROTO_DIR)/*.proto

gen-auth-proto:
	mkdir -p gen/api/auth_v1
	protoc --proto_path api/auth_v1  \
		--go_out=gen/api/auth_v1 \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen/api/auth_v1 \
		--go-grpc_opt=paths=source_relative \
		api/auth_v1/*.proto


gen-access-proto:
	mkdir -p gen/api/access_v1
	protoc --proto_path api/access_v1  \
		--go_out=gen/api/access_v1 \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen/api/access_v1 \
		--go-grpc_opt=paths=source_relative \
		api/access_v1/*.proto

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
			@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
            			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
            			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
            			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
            			rm -rf vendor.protogen/openapiv2 ;\
			fi

gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

grpc-load-test:
	ghz \
		--proto api/user_v1/user.proto \
		--call user_v1.UserV1/Get \
		--data '{"id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/user_v1/user.proto \
		--call user_v1.UserV1/Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051