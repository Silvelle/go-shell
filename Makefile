BIN := bin/repl
SRC := ./cmd/emulator

.PHONY: build

build:
	mkdir -p bin
	go build -o $(BIN) $(SRC)

run:
	go run $(SRC)

clean:
	rm -rf bin
