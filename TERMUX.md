# Vivo Y35: быстрый запуск

Полная инструкция находится в [`DEPLOYMENT.md`](./DEPLOYMENT.md). Этот файл оставлен как короткая памятка.

## Первый запуск

```bash
termux-setup-storage
cd ~
git clone YOUR_REPOSITORY_URL chirik
cd ~/chirik
bash scripts/termux-install.sh
bash scripts/termux-start.sh
```

После запуска Cloudpub должен проксировать только `127.0.0.1:8080`.

Открываемый пользователями адрес: `https://ВАШ-ПУБЛИЧНЫЙ-АДРЕС-CLOUDPUB`.

## Управление

```bash
bash ~/chirik/scripts/termux-start.sh
bash ~/chirik/scripts/termux-stop.sh
tail -f ~/storage/shared/Chirik/data/logs/server.log
tail -f ~/storage/shared/Chirik/data/logs/realtime.log
```

Данные находятся в `Внутренняя память/Chirik/data/`.
