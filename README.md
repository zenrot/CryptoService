# CryptoService

REST-сервис для отслеживания криптовалют, получения цен, истории, статистики и управления расписанием обновлений. Поддерживает внутреннюю авторизацию (JWT) и два типа хранилища: RAM и PostgreSQL.

## Возможности

- Регистрация и логин пользователей с выдачей JWT
- Добавление/удаление криптовалют в трекинг
- Получение списка, карточки, истории и статистики цены
- Ручное обновление цены и массовое обновление всех монет
- Расписание автоподкачки цен (включение/выключение, интервал)

## Стек

- Go 1.24
- HTTP API на Gin
- PostgreSQL (опционально)
- Coingecko API для получения цен

## Быстрый старт

1. Установите Go 1.24+
2. Проверьте конфигурацию в [config/config.yaml](config/config.yaml)
3. Запустите сервис:

```bash
go run cryptoserver.go -configPath config/config.yaml
```

Сервер по умолчанию слушает адрес из `http-config.address`.

## Конфигурация

Файл: [config/config.yaml](config/config.yaml)

```yaml
coingeckoKey: "<ваш_ключ>"
authorizer_type: "internal"
http-config:
	jwt_key: "<секрет>"
	address: "localhost:8090"
postgres-storage:
	host: "localhost"
	port: "5432"
	user: "postgres"
	password: "<пароль>"
	dbname: "crypto_service"
```

Параметры:

- `coingeckoKey` — ключ Coingecko API
- `authorizer_type` — тип авторизации, сейчас используется `internal`
- `http-config.jwt_key` — ключ подписи JWT
- `http-config.address` — адрес HTTP сервера
- `storage_type` (в YAML — `storage_type`) — `ram` или `postgres` (по умолчанию `ram`)
- `postgres-storage.*` — параметры подключения к PostgreSQL

## Запуск с PostgreSQL

1. Создайте базу `crypto_service`.
2. Укажите параметры подключения в `config/config.yaml`.
3. Установите `storage_type: "postgres"`.
4. Запустите сервис. Таблицы создаются автоматически при старте.

## API

Все эндпоинты, кроме `/auth/*`, требуют заголовок:

```
Authorization: Bearer <token>
```

### Аутентификация

- `POST /auth/register`
	- Body: `{ "username": "user", "password": "pass" }`
	- Ответ: `{ "token": "..." }`

- `POST /auth/login`
	- Body: `{ "username": "user", "password": "pass" }`
	- Ответ: `{ "token": "..." }`

### Криптовалюты

- `GET /crypto` — список трекаемых монет
- `GET /crypto/:symbol` — информация по монете
- `GET /crypto/:symbol/history` — история цены
- `GET /crypto/:symbol/stats` — статистика (min/max/avg/count)
- `POST /crypto` — добавить монету
	- Body: `{ "symbol": "BTC" }`
- `PUT /crypto/:symbol/refresh` — обновить цену вручную
- `DELETE /crypto/:symbol` — удалить монету из трекинга

### Расписание обновлений

- `GET /schedule` — текущие настройки
- `PUT /schedule` — изменить настройки
	- Body: `{ "enabled": true, "interval_seconds": 60 }`
- `POST /schedule/trigger` — принудительное обновление всех цен

## Примеры запросов

```bash
curl -X POST http://localhost:8090/auth/register \
	-H "Content-Type: application/json" \
	-d '{"username":"demo","password":"secret"}'

curl -X POST http://localhost:8090/crypto \
	-H "Authorization: Bearer <token>" \
	-H "Content-Type: application/json" \
	-d '{"symbol":"BTC"}'

curl -X GET http://localhost:8090/crypto/BTC \
	-H "Authorization: Bearer <token>"
```

## Тесты

В проекте есть интеграционные тесты на Python.

```bash
make install
make test
```

Дополнительные тесты расписания:

```bash
make test SCHEDULE=1
```

## Замечания

- Для корректной работы нужны доступ к интернету и валидный `coingeckoKey`.
- При использовании `ram` данные не сохраняются между перезапусками.
- В `postgres` режиме данные сохраняются и используются при старте.
