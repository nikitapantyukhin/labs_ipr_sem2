# Kubernetes manifests for lab 6

Для лабораторной работы N6 используются новые разделенные каталоги:

- `infra/postgres` - инфраструктурный PostgreSQL: StatefulSet, Headless Service, PVC, dev/prod через Kustomize и Helm.
- `kustomization` - приложение backend/frontend через Kustomize, без манифестов PostgreSQL.
- `helm/sport-platform-app` - приложение backend/frontend через Helm, без манифестов PostgreSQL.

Плоские файлы `00-namespace.yaml` ... `08-hpa.yaml` оставлены как вариант из лабораторной N5.

## Kustomize

```bash
kubectl apply -k k8s/infra/postgres/kustomization/overlays/dev
kubectl wait --for=condition=ready pod/postgres-0 -n sport-platform-infra --timeout=180s
kubectl apply -k k8s/kustomization/overlays/dev
```

## Helm

```bash
helm upgrade --install sport-platform-db ./k8s/infra/postgres/helm/postgres-infra \
  --namespace sport-platform-infra --create-namespace \
  -f ./k8s/infra/postgres/helm/postgres-infra/values-dev.yaml

helm upgrade --install sport-platform-app ./k8s/helm/sport-platform-app \
  --namespace sport-platform --create-namespace \
  -f ./k8s/helm/sport-platform-app/values-dev.yaml
```

## Проверка

```bash
kubectl get pods,svc,hpa,ingress -n sport-platform
kubectl port-forward service/backend-service 18080:8080 -n sport-platform
curl http://localhost:18080/health
```

Подробное описание, контракт БД и ответы на контрольные вопросы находятся в корневом `README.md`.
