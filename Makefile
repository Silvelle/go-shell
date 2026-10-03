BIN := bin/emulator
SRC := ./src

.PHONY: build run test lint clean

build:
	mkdir -p bin
	go build -o $(BIN) $(SRC)

run:
	go run $(SRC)

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin