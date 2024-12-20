.PHONY: run
run: 
	go run cmd/main.go ${ARGS}

.PNONY: test
test:
	go test -v ./cmd/...
	go test -v ./pkg/...