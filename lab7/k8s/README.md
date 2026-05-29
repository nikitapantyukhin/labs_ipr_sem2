# Kubernetes manifests for lab 5

## 1. Build local images

Run from the project root while Docker Desktop is running:

```bash
docker build -t sport-platform-backend:1.0.0 ./backend
docker build --build-arg NEXT_PUBLIC_API_URL=http://api.sport.localhost -t sport-platform-frontend:1.0.0 ./frontend
```

Docker Desktop Kubernetes can use these local images because the manifests set `imagePullPolicy: IfNotPresent`.

## 2. Install Nginx Ingress Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.1/deploy/static/provider/cloud/deploy.yaml
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=180s
```

## 3. Deploy the application

```bash
kubectl apply -f k8s/
kubectl apply -f k8s/monitoring/
kubectl get pods -n sport-platform
kubectl get services -n sport-platform
```

Application URLs:

- Frontend through NodePort: http://localhost:30080
- Frontend through Ingress: http://sport.localhost
- Backend API through Ingress: http://api.sport.localhost/health
- Prometheus: http://localhost:30090
- Grafana: http://localhost:30091 (`admin` / `admin`)

Demo login:

```text
coach@sport.local
Demo12345!
```

## 4. HPA notes

The HPA manifests require Metrics Server. If it is not installed in Docker Desktop Kubernetes, add it first:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type=strategic --patch-file k8s/support/metrics-server-patch.json
kubectl rollout status deployment/metrics-server -n kube-system --timeout=240s
```

Check autoscaling:

```bash
kubectl get hpa -n sport-platform
```

## 5. Useful checks

```bash
kubectl logs deployment/backend -n sport-platform
kubectl logs deployment/frontend -n sport-platform
kubectl port-forward service/grafana-service 13001:3000 -n sport-platform
kubectl port-forward service/tempo-service 13200:3200 -n sport-platform
kubectl port-forward service/loki-service 3100:3100 -n sport-platform
kubectl delete namespace sport-platform
```

## 6. Lab 7 observability

The backend exposes Prometheus metrics on `/metrics` and exports OTLP traces to
`http://tempo-service:4318` through `OTEL_EXPORTER_OTLP_ENDPOINT`.

The standalone Prometheus manifest scrapes:

- `backend-service:8080` with `metrics_path: /metrics`
- `frontend-service:80` with `metrics_path: /api/metrics`

If Prometheus Operator is installed, apply the optional ServiceMonitor manifests:

```bash
kubectl apply -f k8s/observability/prometheus-operator/servicemonitor.yaml
```
