{{- define "kogger.name" -}}
{{- default "kogger" .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Helm required labels */}}
{{- define "kogger.labels" -}}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
app: "{{ template "kogger.name" . }}"
layer: vault
{{- end -}}

{{/* matchLabels */}}
{{- define "kogger.matchLabels" -}}
release: {{ .Release.Name }}
app: "{{ template "kogger.name" . }}"
{{- end -}}