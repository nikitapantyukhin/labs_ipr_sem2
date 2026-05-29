# Lab 7 observability

Этот каталог содержит локальную платформенную часть для лабораторной работы №7:
Prometheus, Grafana, Tempo, provisioning datasources и готовый dashboard.

## Локальный запуск

Из корня проекта:

```powershell
docker compose -f docker-compose.yml -f docker-compose.observability.yml up -d --build
```

Сервисы:

- Backend: http://localhost:8080/health
- Backend metrics: http://localhost:8080/metrics
- Prometheus targets: http://localhost:9090/targets
- Grafana: http://localhost:3001, логин `admin`, пароль `admin`
- Tempo API: http://localhost:3200/ready

Prometheus должен видеть targets:

- `sport-platform-backend` с `metrics_path: /metrics`
- `sport-platform-frontend` с `metrics_path: /api/metrics`
- `prometheus`

## Проверка метрик

Сгенерируйте несколько запросов:

```powershell
Invoke-WebRequest -UseBasicParsing http://localhost:8080/health
Invoke-WebRequest -UseBasicParsing http://localhost:8080/metrics
```

В Grafana откройте папку `Lab7` и dashboard
`Sport Platform Backend Observability`.

Полезные PromQL-запросы:

```promql
sum by (route) (rate(sport_platform_backend_http_requests_total[1m]))
histogram_quantile(0.95, sum by (le, route) (rate(sport_platform_backend_http_request_duration_seconds_bucket[5m])))
sum by (event) (increase(sport_platform_backend_business_events_total[15m]))
```

## Проверка трейсов

Backend получает OTLP endpoint из overlay:

```text
OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318
OTEL_SERVICE_NAME=sport-platform-backend
```

После HTTP-запросов к API откройте Grafana:
`Explore -> Tempo -> Search`, выберите service `sport-platform-backend`
и найдите trace по маршрутам `/health`, `/users/login` или другим API endpoint.

## Kubernetes

Базовый стек без Prometheus Operator:

```powershell
kubectl apply -f k8s/
kubectl apply -f k8s/monitoring/
```

Проверка:

```powershell
kubectl port-forward service/prometheus-service 19090:9090 -n sport-platform
kubectl port-forward service/grafana-service 13001:3000 -n sport-platform
kubectl port-forward service/tempo-service 13200:3200 -n sport-platform
```

Если в кластере установлен Prometheus Operator или `kube-prometheus-stack`,
дополнительно примените ServiceMonitor:

```powershell
kubectl apply -f k8s/observability/prometheus-operator/servicemonitor.yaml
```

## Скриншоты для отчёта

Сложите файлы в `docs/screenshots/lab7/`:

- `prometheus-targets.png` - Prometheus Targets, backend/frontend jobs `UP`;
- `grafana-dashboard.png` - dashboard `Sport Platform Backend Observability`;
- `tempo-trace.png` - Grafana Explore Tempo с раскрытым trace и span-ами.
