# PostgreSQL infrastructure

Отдельный контур инфраструктуры для лабораторной работы N6. В этом каталоге живет только PostgreSQL: StatefulSet, Headless Service, Secret и PVC через `volumeClaimTemplates`.

## Контракт для приложения

Приложение не создает PostgreSQL самостоятельно. Оно получает параметры подключения из своего ConfigMap/Secret и использует стабильное DNS-имя StatefulSet.

| Окружение | Namespace БД | Host | Port | DB | User |
| --- | --- | --- | --- | --- | --- |
| dev | `sport-platform-infra` | `postgres-0.postgres.sport-platform-infra.svc.cluster.local` | `5432` | `sport_platform` | `sport` |
| prod | `sport-platform-infra-prod` | `postgres-0.postgres.sport-platform-infra-prod.svc.cluster.local` | `5432` | `sport_platform` | `sport` |

Пароль задается в Secret инфраструктуры. В учебных манифестах лежат демонстрационные значения, для реального запуска их нужно заменить через overlay, Helm values или CI/CD secret.

## Деплой через Kustomize

```bash
kubectl apply -k k8s/infra/postgres/kustomization/overlays/dev
kubectl get pods,pvc,svc -n sport-platform-infra -l app.kubernetes.io/name=postgres
```

Prod-вариант:

```bash
kubectl apply -k k8s/infra/postgres/kustomization/overlays/prod
kubectl get pods,pvc,svc -n sport-platform-infra-prod -l app.kubernetes.io/name=postgres
```

## Деплой через Helm

```bash
helm upgrade --install sport-platform-db ./k8s/infra/postgres/helm/postgres-infra \
  --namespace sport-platform-infra --create-namespace \
  -f ./k8s/infra/postgres/helm/postgres-infra/values-dev.yaml
```

Prod-вариант:

```bash
helm upgrade --install sport-platform-db ./k8s/infra/postgres/helm/postgres-infra \
  --namespace sport-platform-infra-prod --create-namespace \
  -f ./k8s/infra/postgres/helm/postgres-infra/values-prod.yaml
```

## Проверка

```bash
kubectl wait --for=condition=ready pod/postgres-0 -n sport-platform-infra --timeout=180s
kubectl exec -n sport-platform-infra postgres-0 -- pg_isready -U sport -d sport_platform
```

PVC не удаляется автоматически при удалении StatefulSet. Это важно для данных: перед удалением PVC нужно явно понимать, что данные больше не нужны.
