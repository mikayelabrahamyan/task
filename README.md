export PATH="/Users/mikayelabrahamyan/Downloads/protoc-26.0-osx-x86_64/bin:$PATH"

protoc --proto_path=. --go_out=go --go-grpc_out=go marketplace.proto

(cd go && go run server.go)

(cd go && go test)

npm install -g grpc-tools grpc_tools_node_protoc_ts

grpc_tools_node_protoc \
    --proto_path=. \
    --js_out=import_style=commonjs,binary:./next/src/grpc \
    --grpc_out=grpc_js:./next/src/grpc \
    --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` \
    marketplace.proto

grpc_tools_node_protoc \
    --proto_path=. \
    --plugin=protoc-gen-ts=`which protoc-gen-ts` \
    --ts_out=grpc_js:./next/src/grpc \
    marketplace.proto

(cd next && npm run dev)