set dotenv-load

# Golang setup:
export GOARCH := "amd64"
export GOOS := "linux"

# Protoc plugins:
protoc_gen_go_version  := "v1.31.0"
protoc_gen_go_grpc_version  := "v1.3.0"

# Api protobuf scheme:
api_version_tag := "v0.0.2"
api_scheme_destination := "./api"
api_file_name := "mergerapi.proto"

generated_pb_package_destination := "./internal/api"

# build settings:
build_output_file := "./bin/tg-client-go"

default: run

run *FLAGS:
    go run  {{FLAGS}} ./cmd/bot \
    --host='localhost:32256' \
    --tg-token=$TGTOKEN \
    --tg-chat-id=-4052858574 \
    --tg-x-api-key=e9a8a738-fc38-4b4c-be05-9b9267372b94 \
    --vk-token=$VKTOKEN \
    --vk-peer-id=2000000002 \
    --vk-x-api-key=c58f1554-8284-4926-97cb-8caa9e3e234f

build *FLAGS:
    go build -o {{build_output_file}} {{FLAGS}}  ./cmd/bot

# required:
#   0. go programming language;
#   1. proto compiler - protoc;
#   2. add to .bashrc PATH="$PATH:$(go env GOPATH)/bin"
init:
    go mod tidy
    just install-deps
    just get-api
    just gen-pb

gen-pb out=generated_pb_package_destination scheme=(api_scheme_destination+"/"+api_file_name):
    mkdir -p {{out}}
    protoc --go_out={{out}} --go-grpc_out={{out}} {{scheme}}

get-api branch=api_version_tag dest=api_scheme_destination file=api_file_name:
	./scripts/download-api-scheme-v2.sh -b {{branch}} -d {{dest}} -f {{file}}

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@{{protoc_gen_go_version}}
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@{{protoc_gen_go_grpc_version}}