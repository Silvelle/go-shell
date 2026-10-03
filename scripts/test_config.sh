#!/bin/sh
# Stage 2: config file and priority of command-line flags over it.
cd "$(dirname "$0")/.." || exit 1
make build >/dev/null || exit 1
BIN=./bin/emulator
TMP=$(mktemp -d)

echo "=== --config only"
$BIN --config configs/emulator.ini --log "$TMP/log.csv"

echo "=== flags override the config file (vfs and log)"
$BIN --config configs/emulator.ini --vfs vfs/deep --log "$TMP/override.csv"

echo "=== config file that does not exist"
$BIN --config configs/missing.ini
echo "exit code: $?"

echo "=== broken config file"
$BIN --config configs/broken.ini
echo "exit code: $?"

echo "=== unknown key in config file"
$BIN --config configs/unknown_key.ini
echo "exit code: $?"

rm -rf "$TMP"
