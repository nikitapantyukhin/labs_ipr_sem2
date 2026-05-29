FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@v3.26.0

COPY . .

RUN go build -ldflags="-w -s" -o go-server main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/go-server /app/go-server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/sqlc/migrations ./migrations

COPY entrypoint.sh /app/entrypoint.sh
RUN ls -la /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

CMD ["/app/entrypoint.sh"]
