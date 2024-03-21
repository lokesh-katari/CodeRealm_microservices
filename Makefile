run:
	@go run cmd/auth/main.go
	
watcher:
	find . -name '*.go' -o | entr -r make
# generate_protobufs:
# 	protoc --plugin=./frontend/node_modules/.bin/protoc-gen-ts_proto --ts_proto_out=./internal/pkg ./internal/pkg/grpc/proto/auth.proto
#     # protoc --go_out=./internal/pkg/grpc/proto --go-grpc_out=./internal/pkg/grpc/proto  auth.proto
# generate_protobufs:
# 	protoc --proto_path=./internal/pkg/grpc/proto --plugin=./frontend/node_modules/.bin/protoc-gen-ts_proto --ts_proto_out=./internal/pkg/grpc/proto ./internal/pkg/grpc/proto/auth.proto

generate_protobufs:
	mkdir -p ./src
	protoc -I=./internal/pkg/grpc/proto ./internal/pkg/grpc/proto/auth.proto \
	--js_out=import_style=commonjs:./frontend/src/proto \
	--grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/proto \
