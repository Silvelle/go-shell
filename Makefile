# Arguments passed to the emulator by "make run", for example:
#   make run ARGS="--vfs vfs/deep --log log.csv"
ARGS ?=

# Window backends on Linux. If the Wayland development headers are
# installed, both X11 and Wayland are built; otherwise only X11, which
# also works on Wayland desktops through XWayland.
# Override with TAGS=x11 (X11 only) or TAGS= (both).
TAGS ?= $(if $(wildcard /usr/include/wayland-client-core.h),,x11)

# The program gets the .exe suffix when built on Windows.
EXT := $(if $(filter windows,$(shell go env GOOS)),.exe,)
BIN := bin/emulator$(EXT)
SRC := ./src

.PHONY: build run test lint clean

build:
	mkdir -p bin
	go build -tags "$(TAGS)" -o $(BIN) $(SRC)

run: build
	./$(BIN) $(ARGS)

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin
