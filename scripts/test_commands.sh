#!/bin/sh
# Stage 4: ls, cd, rev and find on the deep VFS, including errors.
cd "$(dirname "$0")/.." || exit 1
make build >/dev/null || exit 1
./bin/emulator --vfs vfs/deep --script scripts/stage4.txt
