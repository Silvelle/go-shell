BIN := bin/emulator
SRC := ./src

.PHONY: build run test clean

build:
	mkdir -p bin
	go build -o $(BIN) $(SRC)

run:
	go run $(SRC)

test:
	go test -v ./...

clean:
	rm -rf bin