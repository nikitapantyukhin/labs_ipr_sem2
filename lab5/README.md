# Лабораторная работа №5 Kubernetes

Fullstack-приложение для спортивных секций: Go API, PostgreSQL, MinIO и Next.js интерфейс. Для лабораторной проект подготовлен к запуску в Docker Compose и развертыванию в Kubernetes.

## Что реализовано

- `backend` - Go + Gin API, JWT, PostgreSQL через sqlc, миграции goose, MinIO для вложений.
- `frontend` - Next.js, NextAuth Credentials, React Query, Tailwind CSS.
- `docker-compose.yml` - локальный запуск frontend, backend, PostgreSQL и MinIO.
- `k8s/` - Kubernetes-манифесты приложения и инфраструктуры.
- `k8s/monitoring/` - Prometheus, Loki и Promtail для мониторинга и логирования.

## Закрытие требований лабораторной

| Требование | Где сделано |
| --- | --- |
| Развернуть свое приложение из прошлой лабораторной | `k8s/05-backend.yaml`, `k8s/06-frontend.yaml` |
| Deployment для основных компонентов | backend, frontend, PostgreSQL, MinIO, Prometheus, Loki |
| Service для компонентов | `backend-service`, `frontend-service`, `postgres-service`, `minio-service` |
| Зависимости приложения в кластере | PostgreSQL PVC + Deployment, MinIO PVC + Deployment |
| ConfigMap | `k8s/01-configmap.yaml` |
| Secret | `k8s/02-secret.yaml` |
| Ingress | `k8s/07-ingress.yaml` + nginx ingress controller |
| Monitoring | `/metrics`, `/api/metrics`, `k8s/monitoring/prometheus.yaml` |
| Logging | `k8s/monitoring/loki.yaml`, `k8s/monitoring/promtail.yaml` |
| Horizontal Pod Autoscaler | `k8s/08-hpa.yaml` + Metrics Server |
| Health checks | `/health`, `/api/health`, readiness/liveness probes |


## Быстрый запуск через Docker Compose

```bash
cp .env.example .env
docker compose up -d --build
```

После запуска:

- Frontend: http://localhost:3000
- Backend health-check: http://localhost:8080/health
- Backend metrics: http://localhost:8080/metrics
- Frontend health-check: http://localhost:3000/api/health
- MinIO console: http://localhost:9001
- PostgreSQL на хосте: `localhost:5433`

## Запуск в Kubernetes

Нужны Docker Desktop с включенным Kubernetes и `kubectl`.

Выберите контекст Docker Desktop:

```bash
kubectl config use-context docker-desktop
```

Соберите локальные образы:

```bash
docker build -t sport-platform-backend:1.0.0 ./backend
docker build --build-arg NEXT_PUBLIC_API_URL=http://api.sport.localhost -t sport-platform-frontend:1.0.0 ./frontend
```

Установите Ingress Controller:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.1/deploy/static/provider/cloud/deploy.yaml
kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout=240s
```

Установите Metrics Server для HPA:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type=strategic --patch-file k8s/support/metrics-server-patch.json
kubectl rollout status deployment/metrics-server -n kube-system --timeout=240s
```

Если нужно создать patch вручную, используйте PowerShell:

```powershell
@'
{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server","args":["--cert-dir=/tmp","--secure-port=10250","--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname","--kubelet-use-node-status-port","--metric-resolution=15s","--kubelet-insecure-tls"]}]}}}}
'@ | Set-Content -Encoding UTF8 -NoNewline k8s/support/metrics-server-patch.json

kubectl patch deployment metrics-server -n kube-system --type=strategic --patch-file k8s/support/metrics-server-patch.json
kubectl rollout status deployment/metrics-server -n kube-system --timeout=240s
```

Разверните приложение:

```bash
kubectl apply -f k8s/
kubectl apply -f k8s/monitoring/
```

Проверка статуса:

```bash
kubectl get pods,svc,hpa,ingress -n sport-platform
kubectl top pods -n sport-platform
```

## Проверка Kubernetes

Проверить frontend:

```bash
kubectl port-forward service/frontend-service 13080:80 -n sport-platform
```

Откройте http://localhost:13080 или проверьте:

```bash
curl http://localhost:13080/api/health
```

Проверить backend:

```bash
kubectl port-forward service/backend-service 18080:8080 -n sport-platform
curl http://localhost:18080/health
curl http://localhost:18080/metrics
```

Проверить Ingress:

```powershell
Invoke-WebRequest -UseBasicParsing -Headers @{Host='sport.localhost'} http://localhost/api/health
Invoke-WebRequest -UseBasicParsing -Headers @{Host='api.sport.localhost'} http://localhost/health
```

Проверить Prometheus:

```bash
kubectl port-forward service/prometheus-service 19090:9090 -n sport-platform
```

Откройте http://localhost:19090/targets. Targets `backend` и `frontend` должны быть `up`.

Проверить Loki:

```bash
kubectl exec deployment/loki -n sport-platform -- wget -qO- http://localhost:3100/ready
```

Проверить HPA:

```bash
kubectl get hpa -n sport-platform
```

В колонке `TARGETS` должны быть значения CPU, например `1%/60%`, а не `<unknown>`.

## Что было проверено

Проверено локально:

- `docker compose up -d --build` запускает все сервисы.
- `GET http://localhost:8080/health` возвращает `{"status":"ok"}`.
- `GET http://localhost:3000/api/health` возвращает `{"status":"ok","service":"frontend"}`.
- Демо-логин `coach@sport.local / Demo12345!` работает.
- `GET /clubs/` с JWT возвращает 3 демо-секции.

Проверено в Kubernetes:

- Все pod в namespace `sport-platform` находятся в `Running`.
- Backend и frontend доступны через Service port-forward.
- Ingress Controller установлен и маршрутизирует запросы по Host header.
- Prometheus видит targets `backend` и `frontend` в состоянии `up`.
- Loki отвечает `ready`, Promtail запущен как DaemonSet.
- Metrics Server установлен, `kubectl top pods` работает.
- HPA для backend и frontend получает CPU-метрики.
