#!/data/data/com.termux/files/usr/bin/bash
set -euo pipefail

APP_DIR="$HOME/chirik"
DATA_DIR="$HOME/storage/shared/Chirik/data"
WEB_DIR="$HOME/storage/shared/Chirik/web"
BIN_DIR="$APP_DIR/bin"
ENV_FILE="$APP_DIR/.env"

pkg update -y
pkg install -y golang nodejs-lts git rsync openssl clang make
termux-setup-storage

mkdir -p "$DATA_DIR/logs" "$WEB_DIR" "$BIN_DIR"
mkdir -p "$HOME/.termux/boot"
cp "$APP_DIR/scripts/termux-boot.sh" "$HOME/.termux/boot/chirik"
chmod +x "$HOME/.termux/boot/chirik"

# Пересборка должна заменить уже работающий старый процесс.
bash "$APP_DIR/scripts/termux-stop.sh" || true

if [ ! -f "$ENV_FILE" ]; then
  cat > "$ENV_FILE" <<EOF
CHIRIK_JWT_SECRET=$(openssl rand -base64 48 | tr -d '\n')
CHIRIK_REALTIME_SECRET=$(openssl rand -base64 48 | tr -d '\n')
CHIRIK_REALTIME_URL=http://127.0.0.1:8090
CHIRIK_ALLOWED_ORIGIN=*
CHIRIK_DATA_DIR=$DATA_DIR
CHIRIK_WEB_DIR=$WEB_DIR
CHIRIK_SERVER_ADDR=0.0.0.0:8080
CHIRIK_REALTIME_ADDR=127.0.0.1:8090
EOF
  chmod 600 "$ENV_FILE"
fi

cd "$APP_DIR"
go build -o "$BIN_DIR/chirik-server" ./cmd/server
go build -o "$BIN_DIR/chirik-realtime" ./cmd/realtime

cd "$APP_DIR/frontend"
npm install
npm run build
rsync -a --delete build/ "$WEB_DIR/"

echo "Установка завершена. Запуск: $APP_DIR/scripts/termux-start.sh"
