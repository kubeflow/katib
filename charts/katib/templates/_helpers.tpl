{{/*
Expand the name of the chart.
*/}}
{{- define "katib.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "katib.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "katib.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "katib.labels" -}}
helm.sh/chart: {{ include "katib.chart" . }}
{{ include "katib.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "katib.selectorLabels" -}}
app.kubernetes.io/name: {{ include "katib.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the controller service account to use
*/}}
{{- define "katib.controller.serviceAccountName" -}}
{{- if .Values.controller.serviceAccount.create }}
{{- default (printf "%s-controller" (include "katib.fullname" .)) .Values.controller.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.controller.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Controller labels
*/}}
{{- define "katib.controller.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: controller
{{- with .Values.controller.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Controller selector labels
*/}}
{{- define "katib.controller.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: controller
{{- end }}

{{/*
UI labels
*/}}
{{- define "katib.ui.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: ui
{{- with .Values.ui.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
UI selector labels
*/}}
{{- define "katib.ui.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: ui
{{- end }}

{{/*
DB Manager labels
*/}}
{{- define "katib.dbManager.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: db-manager
{{- with .Values.dbManager.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
DB Manager selector labels
*/}}
{{- define "katib.dbManager.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: db-manager
{{- end }}

{{/*
MySQL labels
*/}}
{{- define "katib.mysql.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: mysql
{{- end }}

{{/*
MySQL selector labels
*/}}
{{- define "katib.mysql.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: mysql
{{- end }}

{{/*
PostgreSQL labels
*/}}
{{- define "katib.postgres.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: postgres
{{- end }}

{{/*
PostgreSQL selector labels
*/}}
{{- define "katib.postgres.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: postgres
{{- end }}

{{/*
Webhook labels
*/}}
{{- define "katib.webhook.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: webhook
{{- end }}

{{/*
Webhook selector labels
*/}}
{{- define "katib.webhook.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: webhook
{{- end }}

{{/*
Webhook service name
*/}}
{{- define "katib.webhook.serviceName" -}}
{{ include "katib.fullname" . }}-controller
{{- end }}

{{/*
Webhook secret name
*/}}
{{- define "katib.webhook.secretName" -}}
{{ include "katib.fullname" . }}-webhook-cert
{{- end }}

{{/*
Database labels
*/}}
{{- define "katib.db.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: database
{{- end }}

{{/*
Database selector labels
*/}}
{{- define "katib.db.selectorLabels" -}}
{{ include "katib.selectorLabels" . }}
app.kubernetes.io/component: database
{{- end }}

{{/*
Resolve image repository and tag
*/}}
{{- define "katib.image" -}}
{{- $registry := .registry | default .global.imageRegistry -}}
{{- $repository := .repository -}}
{{- $tag := .tag | default .global.imageTag -}}
{{- if $registry -}}
{{ printf "%s/%s:%s" $registry $repository $tag }}
{{- else -}}
{{ printf "%s:%s" $repository $tag }}
{{- end -}}
{{- end -}}

{{/*
Database host helper
*/}}
{{- define "katib.database.host" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.enabled -}}
{{ printf "%s-mysql" (include "katib.fullname" .) }}
{{- else -}}
{{ .Values.database.external.host }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.enabled -}}
{{ printf "%s-postgres" (include "katib.fullname" .) }}
{{- else -}}
{{ .Values.database.external.host }}
{{- end -}}
{{- else -}}
{{ .Values.database.external.host }}
{{- end -}}
{{- end -}}

{{/*
Database port helper
*/}}
{{- define "katib.database.port" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.enabled -}}
{{ .Values.database.mysql.service.port }}
{{- else -}}
{{ .Values.database.external.port }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.enabled -}}
{{ .Values.database.postgres.service.port }}
{{- else -}}
{{ .Values.database.external.port }}
{{- end -}}
{{- else -}}
{{ .Values.database.external.port }}
{{- end -}}
{{- end -}}

{{/*
Database name helper
*/}}
{{- define "katib.database.name" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.enabled -}}
{{ .Values.database.mysql.auth.database }}
{{- else -}}
{{ .Values.database.external.database }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.enabled -}}
{{ .Values.database.postgres.auth.database }}
{{- else -}}
{{ .Values.database.external.database }}
{{- end -}}
{{- else -}}
{{ .Values.database.external.database }}
{{- end -}}
{{- end -}}

{{/*
Database username helper
*/}}
{{- define "katib.database.username" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.enabled -}}
{{ .Values.database.mysql.auth.username }}
{{- else -}}
{{ .Values.database.external.username }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.enabled -}}
{{ .Values.database.postgres.auth.username }}
{{- else -}}
{{ .Values.database.external.username }}
{{- end -}}
{{- else -}}
{{ .Values.database.external.username }}
{{- end -}}
{{- end -}}

{{/*
Database secret name helper
*/}}
{{- define "katib.database.secretName" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.auth.existingSecret -}}
{{ .Values.database.mysql.auth.existingSecret }}
{{- else -}}
{{ printf "%s-mysql" (include "katib.fullname" .) }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.auth.existingSecret -}}
{{ .Values.database.postgres.auth.existingSecret }}
{{- else -}}
{{ printf "%s-postgres" (include "katib.fullname" .) }}
{{- end -}}
{{- else -}}
{{- if .Values.database.external.existingSecret -}}
{{ .Values.database.external.existingSecret }}
{{- else -}}
{{ printf "%s-external-db" (include "katib.fullname" .) }}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Database environment variables helper
*/}}
{{- define "katib.database.env" -}}
{{- if eq .Values.database.type "mysql" }}
- name: DB_NAME
  value: "mysql"
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: MYSQL_ROOT_PASSWORD
{{- else if eq .Values.database.type "postgres" }}
- name: DB_NAME
  value: "postgres"
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: POSTGRES_PASSWORD
{{- else }}
- name: DB_NAME
  value: {{ .Values.database.external.type | default "mysql" | quote }}
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: DB_PASSWORD
{{- end }}
{{- end -}}

{{/*
Self-signed certificate for webhook (development only)
*/}}
{{- define "katib.webhook.selfSignedCert" -}}
{{- $altNames := list (printf "%s.%s.svc" (include "katib.webhook.serviceName" .) (include "katib.namespace" .)) (printf "%s.%s.svc.cluster.local" (include "katib.webhook.serviceName" .) (include "katib.namespace" .)) -}}
{{- $ca := genCA "katib-ca" 365 -}}
{{- $cert := genSignedCert "katib-webhook" nil $altNames 365 $ca -}}
{{- if not .Values._generatedCert -}}
{{- $_ := set .Values "_generatedCert" $cert -}}
{{- end -}}
{{- .Values._generatedCert.Cert | b64enc -}}
{{- end -}}

{{/*
Self-signed private key for webhook (development only)
*/}}
{{- define "katib.webhook.selfSignedKey" -}}
{{- $altNames := list (printf "%s.%s.svc" (include "katib.webhook.serviceName" .) (include "katib.namespace" .)) (printf "%s.%s.svc.cluster.local" (include "katib.webhook.serviceName" .) (include "katib.namespace" .)) -}}
{{- $ca := genCA "katib-ca" 365 -}}
{{- $cert := genSignedCert "katib-webhook" nil $altNames 365 $ca -}}
{{- if not .Values._generatedCert -}}
{{- $_ := set .Values "_generatedCert" $cert -}}
{{- end -}}
{{- .Values._generatedCert.Key | b64enc -}}
{{- end -}}

{{/*
Namespace helper
*/}}
{{- define "katib.namespace" -}}
{{- .Values.global.namespace | default .Release.Namespace -}}
{{- end -}} 