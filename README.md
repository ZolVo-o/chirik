# Чирик

Небольшая социальная сеть на Go, SQLite и React.

## Сервер на Vivo Y35

Полная инструкция по установке на Android через Termux и публикации через Cloudpub находится в [`DEPLOYMENT.md`](./DEPLOYMENT.md). Для быстрого запуска также доступна памятка [`TERMUX.md`](./TERMUX.md).

## Запуск backend

```bash
export CHIRIK_JWT_SECRET="замените-на-длинный-случайный-секрет-минимум-32-символа"
export CHIRIK_ALLOWED_ORIGIN="http://localhost:3000"
export CHIRIK_REALTIME_URL="http://localhost:8090"
export CHIRIK_REALTIME_SECRET="внутренний-секрет-для-связи-сервисов"
go run ./cmd/server
```

Сервер доступен на `http://localhost:8080`. В локальном режиме база данных создаётся в `./chirik.db`; для телефона пути задаются через `CHIRIK_DATA_DIR` и описаны в [`DEPLOYMENT.md`](./DEPLOYMENT.md).

## Запуск frontend

```bash
cd frontend
npm install
npm start
```

Для production-сборки выполните `npm run build`, затем скопируйте содержимое `frontend/build` в `web`.

## Отдельный realtime-сервер

Realtime вынесен в самостоятельный SSE-сервис. Он принимает события только от API по внутреннему секрету и раздаёт их подключённым клиентам.

```bash
export CHIRIK_REALTIME_SECRET="внутренний-секрет-для-связи-сервисов"
go run ./cmd/realtime
```

По умолчанию сервис слушает `:8090` только для внутренней связи. Основной сервер проксирует realtime через `/events`, поэтому наружу достаточно опубликовать один порт `8080`.
