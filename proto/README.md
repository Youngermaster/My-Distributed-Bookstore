# Protocol Buffers Definitions

This directory contains Protocol Buffer (protobuf) definitions for gRPC service contracts.

## Files

- `common.proto` - Common messages shared across services
- `catalog.proto` - Catalog Service definitions
- `user.proto` - User Service definitions
- `cart.proto` - Cart Service definitions
- `order.proto` - Order Service definitions
- `payment.proto` - Payment Service definitions
- `inventory.proto` - Inventory Service definitions

## Generating Code

### For Go services:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/*.proto
```

### For TypeScript/Node.js services:
```bash
protoc --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
       --ts_out=. \
       --js_out=import_style=commonjs,binary:. \
       proto/*.proto
```

### For Python services:
```bash
python -m grpc_tools.protoc -I. \
       --python_out=. \
       --grpc_python_out=. \
       proto/*.proto
```

## Usage

Each microservice should copy the relevant proto files to their `proto/` directory and generate code as needed.
