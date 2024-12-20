.PHONY: run
run: 
	go run cmd/main.go ${ARGS}

.PNONY: runb
runb: 
	./dist/ecstop_darwin_arm64_v8.0/ecstop ${ARGS}

.PNONY: test
test:
	go test ./cmd/...
	go test ./pkg/...

.PNONY: snapshot
snapshot:
	goreleaser release --snapshot --clean