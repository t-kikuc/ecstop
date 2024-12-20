.PHONY: run
run: 
	go run cmd/main.go ${ARGS}

.PNONY: test
test:
	go test ./cmd/...
	go test ./pkg/...

.PNONY: snapshot
snapshot:
	goreleaser release --snapshot --clean