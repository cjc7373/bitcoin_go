test:
	go test ./... -cover

test-verbose:
	go test ./... -cover -v

test-cover:
	go test ./... -coverprofile cover.out
	go tool cover -html cover.out
