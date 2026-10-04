#!/bin/sh
# Stage 5: rm and mv on the deep VFS. Afterwards the directory on disk
# is listed to show that it was not changed.
cd "$(dirname "$0")/.." || exit 1
make build >/dev/null || exit 1
./bin/emulator --vfs vfs/deep --script scripts/stage5.txt
echo "=== vfs/deep on disk after the run (unchanged):"
find vfs/deep | sort
