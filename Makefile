LOCAL_BIN := $(CURDIR)/bin
BUF_BUILD := $(LOCAL_BIN)/buf

include vendor.proto.mk

.PHONY: .bin-deps up

.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)
	go install github.com/bufbuild/buf/cmd/buf@v1.54.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

up: .bin-deps
	@echo "Environment is up. Binaries installed to $(LOCAL_BIN)"

.buf-init:
	./bin/buf config init

.buf-generate:
	$(info Generating protobuf code...)
	PATH="$(LOCAL_BIN):$PATH" $(BUF_BUILD) generate

.tidy:
	$(info Tidying up the module...)
	go mod tidy

generate: .buf-generate .tidy
	$(info Generation complete)
