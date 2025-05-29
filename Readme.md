# GoKeeper

## Описание

GoKeeper — сервис для хранения и управления секретами (логины, пароли, заметки и т.д.).

## Стэк

-   Backend: Go (Gin, PostgreSQL)
-   Frontend: React + Vite + TailwindCSS
-   БД: PostgreSQL
-   Docker, docker-compose

## Как запустить проект (docker-compose)

**Перед запуском создай файл `.env` или скопируй `example.env` в `.env` в корне проекта.**

Пример содержимого:

```
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gokeeper
DB_SSL_MODE=disable
PORT=8080
JWT_SECRET=testtest
SECRET_AES_KEY=12345678901234567890123456789012
```

1. Убедись, что установлен Docker и docker-compose
2. В корне проекта:
    ```sh
    docker-compose up --build
    ```
3. Фронт доступен на [http://localhost:5173](http://localhost:5173)
4. Бэкенд (API) на [http://localhost:8080](http://localhost:8080)

## Как запустить фронт отдельно

```sh
cd client
npm install
npm run dev
```

## Как запустить бэкенд отдельно

```sh
# В корне проекта
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=gokeeper
export DB_SSL_MODE=disable
export JWT_SECRET=testtest
export SECRET_AES_KEY=12345678901234567890123456789012

go run ./cmd/server
```

## Примечания

-   Для локального запуска БД можно использовать docker-compose только с сервисом postgres.
-   Все переменные окружения для бэкенда указаны в docker-compose.yml.
