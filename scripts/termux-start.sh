#!/data/data/com.termux/files/usr/bin/bash
set -euo pipefail

APP_DIR="$HOME/chirik"
DATA_DIR="$HOME/storage/shared/Chirik/data"
ENV_FILE="$APP_DIR/.env"
BIN_DIR="$APP_DIR/bin"

if [ ! -f "$ENV_FILE" ]; then
  echo "Нет $ENV_FILE. Сначала запустите termux-install.sh" >&2
  exit 1
fi

# Keep the Android process alive while Termux is in the background.
command -v termux-wake-lock >/dev/null 2>&1 && termux-wake-lock || true
mkdir -p "$DATA_DIR/logs"

if [ -f "$DATA_DIR/realtime.pid" ] && kill -0 "$(cat "$DATA_DIR/realtime.pid")" 2>/dev/null; then
  echo "Realtime уже запущен"
else
  set -a
  . "$ENV_FILE"
  set +a
  nohup "$BIN_DIR/chirik-realtime" >> "$DATA_DIR/logs/realtime.log" 2>&1 &
  echo $! > "$DATA_DIR/realtime.pid"
fi

if [ -f "$DATA_DIR/server.pid" ] && kill -0 "$(cat "$DATA_DIR/server.pid")" 2>/dev/null; then
  echo "Chirik уже запущен"
else
  set -a
  . "$ENV_FILE"
  set +a
  nohup "$BIN_DIR/chirik-server" >> "$DATA_DIR/logs/server.log" 2>&1 &
  echo $! > "$DATA_DIR/server.pid"
fi

echo "Chirik доступен по адресу: http://IP_ТЕЛЕФОНА:8080"
