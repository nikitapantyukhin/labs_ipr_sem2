# Лабораторная работа №7: Observability

## Цель работы

Добавить в проект `sport-platform` стек наблюдаемости:

- Prometheus для сбора метрик;
- Grafana для визуализации;
- Grafana Tempo для distributed tracing;
- OpenTelemetry-инструментирование backend;
- конфигурацию для Docker Compose и Kubernetes.

## Кратко о проекте

`sport-platform` - fullstack-приложение для спортивных секций.

Состав проекта:

- `backend` - Go + Gin API, JWT, PostgreSQL через sqlc, MinIO для вложений;
- `frontend` - Next.js интерфейс;
- `docker-compose.yml` - локальный запуск приложения;
- `docker-compose.observability.yml` - локальный стек наблюдаемости;
- `observability/` - конфиги Prometheus, Grafana и Tempo;
- `k8s/` - Kubernetes-манифесты приложения и observability-стека.

## Что сделано

| Требование | Реализация |
| --- | --- |
| Endpoint `/metrics` | `backend/application/app.go`, `backend/internal/observability/metrics.go` |
| HTTP counter | `sport_platform_backend_http_requests_total` |
| HTTP histogram | `sport_platform_backend_http_request_duration_seconds` |
| Business metrics | `sport_platform_backend_business_events_total` |
| OTLP tracing | `backend/internal/observability/tracing.go` |
| Prometheus scrape | `observability/prometheus/prometheus.yml` |
| Grafana datasource | `observability/grafana/provisioning/datasources/datasources.yml` |
| Grafana dashboard | `observability/grafana/dashboards/sport-platform-backend.json` |
| Kubernetes scrape | `k8s/monitoring/prometheus.yaml` |
| ServiceMonitor | `k8s/observability/prometheus-operator/servicemonitor.yaml` |

## Самостоятельная работа

Для проекта реализована observability-часть:

1. В backend добавлен endpoint `/metrics`.
2. Настроены HTTP-метрики:
   - counter количества запросов;
   - histogram длительности запросов;
   - gauge активных запросов.
3. Labels у HTTP-метрик имеют умеренную кардинальность:
   - `method`;
   - `route`;
   - `status`.
4. Добавлена бизнес-метрика событий приложения:
   - регистрация пользователя;
   - создание секции;
   - создание заявки на вступление;
   - создание тренировки.
5. В backend включён OpenTelemetry OTLP exporter в Tempo.
6. Для Grafana подключены datasource Prometheus и Tempo.
7. Подготовлен dashboard с PromQL-запросами.
8. Добавлены конфиги для Docker Compose и Kubernetes.

## Запуск

### Локальный запуск приложения

```bash
cp .env.example .env
docker compose up -d --build
```

### Локальный запуск observability-стека

```bash
docker compose -f docker-compose.yml -f docker-compose.observability.yml up -d --build
```

После запуска доступны:

- Frontend: `http://localhost:3000`
- Backend health: `http://localhost:8080/health`
- Backend metrics: `http://localhost:8080/metrics`
- Frontend metrics: `http://localhost:3000/api/metrics`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001`
- Tempo: `http://localhost:3200/ready`

## Ход выполнения

### 1. Метрики

В backend реализован endpoint `/metrics`, который отдаёт Prometheus-метрики в текстовом формате.

Используемые метрики:

- `sport_platform_backend_http_requests_total`
- `sport_platform_backend_http_request_duration_seconds`
- `sport_platform_backend_http_requests_in_flight`
- `sport_platform_backend_business_events_total`

### 2. Скрейпинг

Prometheus настроен на сбор метрик:

- backend по `/metrics`;
- frontend по `/api/metrics`;
- собственных метрик Prometheus.

### 3. Трейсинг

Backend экспортирует traces через OTLP HTTP exporter в Grafana Tempo.
Если `OTEL_EXPORTER_OTLP_ENDPOINT` не задан, приложение продолжает работать без трейсинга.

### 4. Grafana

В Grafana подключены datasource:

- Prometheus;
- Tempo.

Dashboard `Sport Platform Backend Observability` содержит панели:

- RPS по маршрутам;
- p95 latency;
- количество бизнес-событий;
- число активных запросов;
- статус backend target.

## Скриншоты

Скриншоты для отчёта должны лежать в `docs/screenshots/lab7/`.

## Вывод

В ходе лабораторной работы в проект `sport-platform` был добавлен полный набор observability-компонентов:
метрики, Prometheus scrape, Grafana dashboard и distributed tracing через Tempo.
Реализация подходит и для локального запуска через Docker Compose, и для развёртывания в Kubernetes.

