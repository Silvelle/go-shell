#!/bin/sh
# Stage 3: runs the emulator with different VFS directories.
cd "$(dirname "$0")/.." || exit 1
make build >/dev/null || exit 1
BIN=./bin/emulator

echo "=== default VFS (no --vfs)"
$BIN --script scripts/stage3.txt

echo "=== minimal VFS: one file"
$BIN --vfs vfs/minimal --script scripts/stage3.txt

echo "=== several files"
$BIN --vfs vfs/files --script scripts/stage3.txt

echo "=== deep VFS: 3+ levels of directories"
$BIN --vfs vfs/deep --script scripts/stage3.txt

echo "=== error: VFS directory does not exist"
$BIN --vfs vfs/no_such_dir --script scripts/stage3.txt

echo "=== error: VFS path is a file, not a directory"
$BIN --vfs README.md --script scripts/stage3.txt
