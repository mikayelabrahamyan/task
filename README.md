# Task Repository

## 1. Project Structure Description
### Folder Structure
- `go/`: Contains the Go project which acts as the gRPC server.
- `next/`: Contains the Next.js project which serves both as the gRPC client and the frontend view.
- `marketplace.proto`: Messaging protocol file for gRPC, used by both the Go and Next.js applications simultaneously.

## 2. External Libraries
### Go
- gRPC-related libraries for implementing the gRPC server functionality.
- `bou.ke/monkey`: Used for mocking in tests due to its ability to patch functions at runtime.

### Next.js
- gRPC-related libraries for enabling gRPC client functionality.

## 3. Project Setup

### Prerequisites
- Install Protocol Buffers Compiler: Follow the [gRPC installation guide](https://grpc.io/docs/protoc-installation/).

### Installing Dependencies

download the required go modules:
```bash
(cd go && go mod download)
```
generate go files from the marketplace.proto file in the go/marketplace-gen folder:
```bash
protoc --proto_path=. --go_out=go --go-grpc_out=go marketplace.proto
```

install nextjs dependencies using npm:
```bash
(cd next && npm ci)
```

generate js ts files from the marketplace.proto file in the next/src/grpc folder:
```bash
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
```

## 4. Running the Application

### Start the Go gRPC Server

run the go gRPC server on the hardcoded localhost and port 50051:
```bash
(cd go && go run server.go)
```

### Start the Next.js Application
run the next.js application, which will connect to the gRPC server and also serve the frontend on hardcoded localhost port 3000:
```bash
(cd next && npm run dev)
```

## 5. Running Tests

### Go Tests

to execute the tests in the Go application:
```bash
(cd go && go test)
```