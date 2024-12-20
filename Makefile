.PHONY: run
run: 
	@echo  [${ARGS}]
	go run src/main.go ${ARGS}

.PNONY: test
test:
	go test -v ./src/...