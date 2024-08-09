{{/* Expand the name of the chart. */}}
{{- define "dashboard.name" -}}
{{-   default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Create default fully qualified name (limited to 63 chars by the DNS) */}}
{{- define "dashboard.fullname" -}}
{{-   if .Values.fullnameOverride }}
{{-     .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{-   else }}
{{-     $name := default .Chart.Name .Values.nameOverride }}
{{-     if contains $name .Release.Name }}
{{-       .Release.Name | trunc 63 | trimSuffix "-" }}
{{-     else }}
{{-       printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{-     end }}
{{-   end }}
{{- end }}

{{/* Create the name of the service account to use */}}
{{- define "dashboard.serviceAccountName" -}}
{{-   if .Values.serviceAccount.create }}
{{-     default (include "dashboard.fullname" .) .Values.serviceAccount.name }}
{{-   else }}
{{-     default "default" .Values.serviceAccount.name }}
{{-   end }}
{{- end }}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "dashboard.chart" -}}
{{-   printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels */}}
{{- define "dashboard.labels" -}}
{{    include "dashboard.selectorLabels" . }}
{{-   if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{-   end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "dashboard.chart" . }}
{{- end }}

{{/* Selector labels */}}
{{- define "dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: dashboard
{{- end }}
