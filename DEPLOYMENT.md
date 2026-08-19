# Развёртывание Чирика на Vivo Y35

Эта инструкция описывает установку сервера на Android-телефон Vivo Y35 через Termux и публикацию через Cloudpub.

## 1. Архитектура

На телефоне работают два процесса:

| Процесс | Внутренний адрес | Назначение |
|---|---|---|
| `chirik-server` | `127.0.0.1:8080` | сайт, API и внешний SSE-маршрут `/events` |
| `chirik-realtime` | `127.0.0.1:8090` | внутренний realtime-сервис |

Пользователи и Cloudpub видят только `8080`. Основной сервер проксирует `/events` на внутренний realtime-процесс.

Пользовательские данные не лежат рядом с бинарниками:

```text
Внутренняя память/Chirik/data/chirik.db   SQLite-база
Внутренняя память/Chirik/data/logs/       логи процессов
Внутренняя память/Chirik/web/             production frontend
```

Исходники и бинарники находятся в Termux home:

```text
~/chirik/                                исходники
~/chirik/bin/                            собранные бинарники
~/chirik/.env                            секреты и настройки
```

Обновление исходников не удаляет базу.

## 2. Установка приложений

1. Установите Termux из F-Droid. Не используйте устаревшую сборку Termux из Google Play.
2. Откройте Termux.
3. Установите Termux:Boot из F-Droid, если нужен автозапуск после перезагрузки.
4. Установите Cloudpub тем способом, который предусмотрен вашим аккаунтом Cloudpub.

В настройках Vivo для Termux отключите оптимизацию батареи и разрешите работу в фоне. Иначе Android может остановить сервер через некоторое время.

## 3. Подготовка памяти телефона

В Termux выполните:

```bash
termux-setup-storage
```

Разрешите доступ к файлам в системном диалоге Android. После этого каталог общей памяти доступен как:

```text
~/storage/shared/
```

## 4. Загрузка проекта

### Вариант A: через Git

Замените `YOUR_REPOSITORY_URL` на URL репозитория:

```bash
cd ~
git clone YOUR_REPOSITORY_URL chirik
cd ~/chirik
```

### Вариант B: через архив

Если архив проекта уже находится в `Download`:

```bash
cd ~
tar -xzf ~/storage/shared/Download/chirik.tar.gz
mv chirik-* chirik
cd ~/chirik
```

Если архив распаковывается сразу в каталог `chirik`, команду `mv` выполнять не нужно.

Проверьте наличие главных файлов:

```bash
ls
ls scripts
```

В каталоге должны быть `go.mod`, `frontend`, `cmd`, `internal` и `scripts`.

## 5. Установка и первая сборка

Запустите установочный скрипт:

```bash
cd ~/chirik
bash scripts/termux-install.sh
```

Скрипт:

- обновляет пакеты Termux;
- устанавливает Go, Node.js, Git, OpenSSL и инструменты сборки SQLite;
- создаёт каталоги на памяти телефона;
- генерирует секреты JWT и realtime;
- собирает backend и realtime-бинарники;
- собирает production frontend;
- копирует frontend в `Внутренняя память/Chirik/web/`;
- добавляет автозапуск Termux:Boot.

Секреты сохраняются в:

```bash
cat ~/chirik/.env
```

Файл имеет права `600`. Не публикуйте его и не отправляйте в Git.

## 6. Запуск сервера

```bash
bash ~/chirik/scripts/termux-start.sh
```

Скрипт запускает realtime и основной сервер в фоне. Повторный запуск безопасен: уже работающие процессы не запускаются второй раз.

Локальная проверка сайта:

```bash
curl -i http://127.0.0.1:8080/
```

Ожидаемый HTTP-код: `200`.

Проверка API без авторизации:

```bash
curl -i http://127.0.0.1:8080/api/conversations
```

Ожидаемый HTTP-код: `401`. Это означает, что API работает и требует токен.

Проверка внутреннего realtime:

```bash
curl -i http://127.0.0.1:8090/health
```

Ожидаемый ответ: `200` и `ok`.

## 7. Настройка Cloudpub

Cloudpub должен туннелировать только этот локальный адрес:

```text
127.0.0.1:8080
```

Не указывайте `127.0.0.1:8090`: этот порт внутренний и не должен быть публичным.

В зависимости от версии Cloudpub настройка выполняется в приложении, личном кабинете или CLI. Значение локального target: `localhost:8080` или `127.0.0.1:8080`.

После запуска туннеля Cloudpub выдаст публичный адрес, например:

```text
https://example.cloudpub.ru
```

Именно этот адрес нужно открывать пользователям. Отдельные адреса для API и realtime не нужны:

```text
https://example.cloudpub.ru/             сайт
https://example.cloudpub.ru/api/...       API
https://example.cloudpub.ru/events        realtime
```

HTTPS Cloudpub также нужен для корректной работы браузеров и `EventSource`.

## 8. Остановка

```bash
bash ~/chirik/scripts/termux-stop.sh
```

Проверка процессов:

```bash
ps -A | grep chirik
```

Проверка портов:

```bash
ss -ltn
```

Ожидаются порты `8080` и внутренний `8090`.

## 9. Просмотр логов

Основной сервер:

```bash
tail -f ~/storage/shared/Chirik/data/logs/server.log
```

Realtime:

```bash
tail -f ~/storage/shared/Chirik/data/logs/realtime.log
```

Если сервер не запускается, сначала остановите старые процессы и повторите запуск:

```bash
bash ~/chirik/scripts/termux-stop.sh
bash ~/chirik/scripts/termux-start.sh
```

Если порт занят:

```bash
ss -ltnp
```

## 10. Обновление проекта

### Обновление через Git

```bash
bash ~/chirik/scripts/termux-stop.sh
cd ~/chirik
git pull
bash scripts/termux-install.sh
bash scripts/termux-start.sh
```

### Обновление архивом

Сначала остановите сервер:

```bash
bash ~/chirik/scripts/termux-stop.sh
```

Затем замените исходники новым архивом, не удаляя каталог данных:

```text
Внутренняя память/Chirik/data/
```

После замены исходников выполните:

```bash
cd ~/chirik
bash scripts/termux-install.sh
bash scripts/termux-start.sh
```

`termux-install.sh` не удаляет базу данных и не пересоздаёт существующие секреты в `.env`.

## 11. Резервная копия базы

Перед резервным копированием остановите сервер, чтобы SQLite-файл был закрыт:

```bash
bash ~/chirik/scripts/termux-stop.sh
mkdir -p ~/storage/shared/Chirik-backups
tar -czf ~/storage/shared/Chirik-backups/chirik-$(date +%Y-%m-%d-%H%M).tar.gz \
  -C ~/storage/shared/Chirik data
bash ~/chirik/scripts/termux-start.sh
```

Архив содержит базу и логи. Храните копии вне телефона.

Для восстановления остановите сервер и распакуйте архив обратно в `Chirik`:

```bash
bash ~/chirik/scripts/termux-stop.sh
tar -xzf ~/storage/shared/Chirik-backups/ИМЯ_АРХИВА.tar.gz \
  -C ~/storage/shared/Chirik
bash ~/chirik/scripts/termux-start.sh
```

## 12. Автозапуск

После установки файл автозапуска находится здесь:

```text
~/.termux/boot/chirik
```

Для проверки:

```bash
ls -l ~/.termux/boot/chirik
```

После перезагрузки телефона Termux:Boot подождёт 10 секунд и запустит оба процесса.

Если автозапуск не сработал:

1. Проверьте, установлен ли Termux:Boot.
2. Откройте Termux:Boot один раз.
3. Отключите оптимизацию батареи для Termux и Termux:Boot в настройках Vivo.
4. Проверьте логи после ручного запуска.

## 13. Настройки `.env`

Основной файл настроек:

```bash
nano ~/chirik/.env
```

Используемые значения:

```text
CHIRIK_JWT_SECRET              секрет авторизации, минимум 32 символа
CHIRIK_REALTIME_SECRET         секрет API -> realtime
CHIRIK_REALTIME_URL            внутренний URL публикации событий
CHIRIK_ALLOWED_ORIGIN          обычно * для SSE через Cloudpub
CHIRIK_DATA_DIR                каталог постоянных данных
CHIRIK_WEB_DIR                 каталог production frontend
CHIRIK_SERVER_ADDR             локальный порт, обычно 127.0.0.1:8080
CHIRIK_REALTIME_ADDR           внутренний порт, по умолчанию 127.0.0.1:8090
```

После изменения `.env` обязательно перезапустите процессы:

```bash
bash ~/chirik/scripts/termux-stop.sh
bash ~/chirik/scripts/termux-start.sh
```

## 14. Частые проблемы

### Открывается пустая страница

Проверьте, что frontend собран и скопирован:

```bash
ls ~/storage/shared/Chirik/web/index.html
bash ~/chirik/scripts/termux-install.sh
```

### `401` на `/api/conversations`

Это нормальный результат без JWT-токена. Для проверки самого API достаточно увидеть `401`, а не `404` или ошибку соединения.

### Ошибка `address already in use`

Остановите старые процессы:

```bash
bash ~/chirik/scripts/termux-stop.sh
ss -ltnp
bash ~/chirik/scripts/termux-start.sh
```

### Сообщения не обновляются сразу

Проверьте оба сервиса:

```bash
curl http://127.0.0.1:8090/health
tail -f ~/storage/shared/Chirik/data/logs/realtime.log
```

Cloudpub должен публиковать `127.0.0.1:8080`, а не `8090`.

### Телефон останавливает сервер

Отключите энергосбережение для Termux, запустите `termux-wake-lock` и держите Termux:Boot установленным:

```bash
termux-wake-lock
bash ~/chirik/scripts/termux-start.sh
```

### Потерялась база

Проверьте каталог:

```bash
ls -lah ~/storage/shared/Chirik/data/
```

Если `chirik.db` удалён, восстановите последний архив из `Chirik-backups`.

## 15. Безопасность

- Не публикуйте порт `8090`.
- Не показывайте содержимое `~/chirik/.env`.
- Не коммитьте `.env` в Git.
- Используйте HTTPS-адрес Cloudpub.
- Регулярно копируйте базу за пределы телефона.
- Не удаляйте `Внутренняя память/Chirik/data/` при обновлении проекта.
