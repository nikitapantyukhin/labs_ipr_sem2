{{- define "sport-platform.labels" -}}
app.kubernetes.io/name: sport-platform
app.kubernetes.io/part-of: sport-platform
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}

{{- define "sport-platform.backendLabels" -}}
app: backend
app.kubernetes.io/component: backend
{{ include "sport-platform.labels" . }}
{{- end -}}

{{- define "sport-platform.frontendLabels" -}}
app: frontend
app.kubernetes.io/component: frontend
{{ include "sport-platform.labels" . }}
{{- end -}}
