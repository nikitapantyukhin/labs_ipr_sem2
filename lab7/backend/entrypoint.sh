#!/bin/sh
set -e

goose -dir ./migrations postgres "host=${HOST} port=${PORT} database=${DB_NAME} user=${USERNAME} password=${PASSWORD} sslmode=disable" up

exec /app/go-server
