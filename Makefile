test:
	go test ./... -cover

test-verbose:
	go test ./... -cover -v

test-bench:
	go test --run=^$$ --bench=. -benchmem ./...

test-cover:
	go test ./... -coverprofile cover.out
	go tool cover -html cover.out

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 internal/network/proto/protocol.proto
	protoc --go_out=. --go_opt=paths=source_relative \
	 internal/block/proto/block.proto
