SHELL := /bin/bash
.DEFAULT_GOAL := build

BINARY_NAME := svc
BUILD_PATH := cmd/build
SERVICE_NAME := booking_v1

# vendor repo sparse dirs
protos-dirs = common google protoc-gen-openapiv2 account_v1
vendor-proto-dirs = vp-common vp-account_v1

.SILENT:

build:
	mkdir -p $(BUILD_PATH)
	CGO_ENABLED=0 go build -o $(BUILD_PATH)/$(BINARY_NAME) cmd/main.go

clean:
	rm -rf $(BUILD_PATH)

lint:
	golangci-lint run

# Установка плагинов protoc
proto-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

refresh-vendor-proto:
	rm -rf vendor-proto && \
		git clone --depth=1 --filter=blob:none --sparse git@github.com:zhamspace/protos.git vendor-proto
	cd vendor-proto && git sparse-checkout set \
		$(protos-dirs)
	cd vendor-proto && rm -rf .git .idea && find . -type f ! -name '*.proto' -delete

init: refresh-vendor-proto
	mkdir -p docs
	mkdir -p api/proto/$(SERVICE_NAME)
	touch api/proto/$(SERVICE_NAME)/service.proto
	mkdir -p cmd
	touch cmd/main.go
	mkdir -p internal
	mkdir -p migrations
	go mod init github.com/zhamspace/$(patsubst %_v1,%,$(SERVICE_NAME))

vp-%:
	mkdir -p pkg/proto
	protoc -I vendor-proto \
	--go_out pkg/proto --go_opt paths=source_relative \
	--go_opt=Mcommon/common.proto=`go list -m `/pkg/proto/common \
	--go-grpc_out pkg/proto --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/proto --grpc-gateway_opt paths=source_relative \
	vendor-proto/$(subst vp-,,$@)/*.proto

generate-vendor-proto-dirs: $(vendor-proto-dirs)

generate-proto-$(SERVICE_NAME):
	mkdir -p pkg/proto
	protoc -I vendor-proto -I api/proto \
	--go_out pkg/proto --go_opt paths=source_relative \
	--go_opt=Mcommon/common.proto=`go list -m `/pkg/proto/common \
	--go-grpc_out pkg/proto --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/proto --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=json_names_for_fields=false,allow_merge=true,merge_file_name=api:docs \
	api/proto/$(SERVICE_NAME)/*.proto

generate-proto: generate-vendor-proto-dirs generate-proto-$(SERVICE_NAME)
