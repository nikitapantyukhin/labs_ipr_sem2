{{- define "postgres-infra.labels" -}}
app.kubernetes.io/name: postgres
app.kubernetes.io/part-of: sport-platform
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}
