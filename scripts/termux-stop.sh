#!/data/data/com.termux/files/usr/bin/bash
set -euo pipefail

DATA_DIR="$HOME/storage/shared/Chirik/data"
for name in realtime server; do
  pid_file="$DATA_DIR/$name.pid"
  if [ -f "$pid_file" ]; then
    kill "$(cat "$pid_file")" 2>/dev/null || true
    rm -f "$pid_file"
  fi
done
command -v termux-wake-unlock >/dev/null 2>&1 && termux-wake-unlock || true
