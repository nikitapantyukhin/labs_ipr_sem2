# Лабораторная работа N6: Kustomize и Helm

Проект: fullstack-приложение для спортивных секций. Backend написан на Go + Gin, frontend на Next.js, база данных PostgreSQL, объектное хранилище MinIO.

В этой лабораторной работа из лабы N5 разнесена на два контура:

- **инфраструктура**: PostgreSQL в отдельном каталоге, со StatefulSet, Headless Service и PVC;
- **приложение**: backend/frontend и их конфигурация через Kustomize и Helm без манифестов базы данных.

## Структура Kubernetes

```text
k8s/
  infra/
    postgres/
      README.md
      kustomization/
        base/
        overlays/
          dev/
          prod/
      helm/
        postgres-infra/
          values-dev.yaml
          values-prod.yaml
          templates/
  kustomization/
    base/
    overlays/
      dev/
      prod/
  helm/
    sport-platform-app/
      values-dev.yaml
      values-prod.yaml
      templates/
```

Каталоги `k8s/kustomization` и `k8s/helm/sport-platform-app` содержат только ресурсы приложения: backend, frontend, Service, Ingress, HPA, ConfigMap и Secret. PostgreSQL там отсутствует.

Старые плоские манифесты `k8s/00-*.yaml` оставлены как материал лабы N5. Для лабы N6 используются новые каталоги выше.

## Контракт PostgreSQL

Приложение подключается к БД по контракту инфраструктуры:

| Окружение | Namespace БД | Host | Port | DB | User |
| --- | --- | --- | --- | --- | --- |
| dev | `sport-platform-infra` | `postgres-0.postgres.sport-platform-infra.svc.cluster.local` | `5432` | `sport_platform` | `sport` |
| prod | `sport-platform-infra-prod` | `postgres-0.postgres.sport-platform-infra-prod.svc.cluster.local` | `5432` | `sport_platform` | `sport` |

Пароль лежит в Secret соответствующего окружения.

## Подготовка образов

Для локального Kubernetes в Docker Desktop соберите образы:

```bash
docker build -t sport-platform-backend:1.0.0 ./backend
docker build --build-arg NEXT_PUBLIC_API_URL=http://api.sport.localhost -t sport-platform-frontend:1.0.0 ./frontend
```

## Вариант 1: запуск через Kustomize

Сначала поднимается PostgreSQL:

```bash
kubectl apply -k k8s/infra/postgres/kustomization/overlays/dev
kubectl wait --for=condition=ready pod/postgres-0 -n sport-platform-infra --timeout=180s
```

Затем приложение:

```bash
kubectl apply -k k8s/kustomization/overlays/dev
kubectl get pods,svc,hpa,ingress -n sport-platform
```

Проверка health-checks:

```bash
kubectl port-forward service/backend-service 18080:8080 -n sport-platform
curl http://localhost:18080/health
```

```bash
kubectl port-forward service/frontend-service 13080:80 -n sport-platform
curl http://localhost:13080/api/health
```

Dev-overlay также публикует frontend как NodePort `30080`, поэтому в Docker Desktop можно открыть `http://localhost:30080`.

## Вариант 2: запуск через Helm

Сначала PostgreSQL:

```bash
helm upgrade --install sport-platform-db ./k8s/infra/postgres/helm/postgres-infra \
  --namespace sport-platform-infra --create-namespace \
  -f ./k8s/infra/postgres/helm/postgres-infra/values-dev.yaml
```

Затем приложение:

```bash
helm upgrade --install sport-platform-app ./k8s/helm/sport-platform-app \
  --namespace sport-platform --create-namespace \
  -f ./k8s/helm/sport-platform-app/values-dev.yaml
```

Проверить итоговые YAML без установки:

```bash
kubectl kustomize k8s/kustomization/overlays/dev
helm template sport-platform-db ./k8s/infra/postgres/helm/postgres-infra -f ./k8s/infra/postgres/helm/postgres-infra/values-dev.yaml
helm template sport-platform-app ./k8s/helm/sport-platform-app -f ./k8s/helm/sport-platform-app/values-dev.yaml
```

## Production overlays и values

Kustomize:

```bash
kubectl apply -k k8s/infra/postgres/kustomization/overlays/prod
kubectl apply -k k8s/kustomization/overlays/prod
```

Helm:

```bash
helm upgrade --install sport-platform-db ./k8s/infra/postgres/helm/postgres-infra \
  --namespace sport-platform-infra-prod --create-namespace \
  -f ./k8s/infra/postgres/helm/postgres-infra/values-prod.yaml

helm upgrade --install sport-platform-app ./k8s/helm/sport-platform-app \
  --namespace sport-platform-prod --create-namespace \
  -f ./k8s/helm/sport-platform-app/values-prod.yaml
```

Prod-настройки отличаются namespace, DNS-именем PostgreSQL, количеством реплик, лимитами и заглушками секретов.

## Что реализовано по заданию

| Требование | Где сделано |
| --- | --- |
| PostgreSQL вынесен из приложения | `k8s/infra/postgres` |
| StatefulSet + Headless Service + PVC | `k8s/infra/postgres/kustomization/base`, `k8s/infra/postgres/helm/postgres-infra/templates` |
| Dev/prod для инфраструктуры | `k8s/infra/postgres/kustomization/overlays/*`, `values-dev.yaml`, `values-prod.yaml` |
| Kustomize для приложения | `k8s/kustomization/base`, `k8s/kustomization/overlays/dev`, `k8s/kustomization/overlays/prod` |
| Helm chart для приложения | `k8s/helm/sport-platform-app` |
| В app-манифестах нет PostgreSQL | `k8s/kustomization` и `k8s/helm/sport-platform-app` |
| Контракт БД совпадает с конфигурацией приложения | `HOST`, `PORT`, `DB_NAME`, `USERNAME`, `PASSWORD` в overlays/values |
