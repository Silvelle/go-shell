#!/bin/sh
# Stage 2: runs the emulator with every command-line parameter.
# The startup script ends with "exit", so every window closes by itself.
cd "$(dirname "$0")/.." || exit 1
make build >/dev/null || exit 1
BIN=./bin/emulator
TMP=$(mktemp -d)

echo "=== no parameters (close the window to continue)"
$BIN

echo "=== --script"
$BIN --script scripts/stage2.txt

echo "=== --vfs (shown in the settings dump)"
$BIN --vfs vfs/files --script scripts/stage2.txt

echo "=== --log"
$BIN --log "$TMP/log.csv" --script scripts/stage2.txt
echo "--- log file:"
cat "$TMP/log.csv"

echo "=== --log to a directory that does not exist"
$BIN --log /no/such/dir/log.csv --script scripts/stage2.txt

echo "=== --script that does not exist"
$BIN --script no_such_script.txt

echo "=== unknown flag"
$BIN --color red
echo "exit code: $?"

rm -rf "$TMP"
