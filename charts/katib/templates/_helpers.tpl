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
{{- default "katib-controller" .Values.controller.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.controller.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Controller service name - matches Kustomize
*/}}
{{- define "katib.controller.serviceName" -}}
katib-controller
{{- end }}

{{/*
UI service account name - matches Kustomize  
*/}}
{{- define "katib.ui.serviceAccountName" -}}
{{- if .Values.ui.serviceAccount.create }}
{{- default "katib-ui" .Values.ui.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.ui.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
UI service name - matches Kustomize
*/}}
{{- define "katib.ui.serviceName" -}}
katib-ui
{{- end }}

{{/*
DB Manager service name - matches Kustomize
*/}}
{{- define "katib.dbManager.serviceName" -}}
katib-db-manager
{{- end }}

{{/*
MySQL service name - matches Kustomize
*/}}
{{- define "katib.mysql.serviceName" -}}
katib-mysql
{{- end }}

{{/*
PostgreSQL service name - matches Kustomize
*/}}
{{- define "katib.postgres.serviceName" -}}
katib-postgres
{{- end }}

{{/*
Validating webhook configuration name - matches Kustomize
*/}}
{{- define "katib.webhook.validatingName" -}}
katib.kubeflow.org
{{- end }}

{{/*
Mutating webhook configuration name - matches Kustomize
*/}}
{{- define "katib.webhook.mutatingName" -}}
katib.kubeflow.org
{{- end }}

{{/*
Trial templates ConfigMap name - matches Kustomize
*/}}
{{- define "katib.trialTemplates.configMapName" -}}
trial-templates
{{- end }}

{{/*
Katib config ConfigMap name - matches Kustomize
*/}}
{{- define "katib.config.configMapName" -}}
katib-config
{{- end }}

{{/*
MySQL secret name - matches Kustomize
*/}}
{{- define "katib.mysql.secretName" -}}
katib-mysql-secrets
{{- end }}

{{/*
Webhook secret name - matches Kustomize
*/}}
{{- define "katib.webhook.secretName" -}}
katib-webhook-cert
{{- end }}

{{/*
Controller labels
*/}}
{{- define "katib.controller.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: controller
katib.kubeflow.org/component: controller
{{- with .Values.controller.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Controller selector labels - Kustomize compatible
*/}}
{{- define "katib.controller.selectorLabels" -}}
katib.kubeflow.org/component: controller
{{- end }}

{{/*
UI labels
*/}}
{{- define "katib.ui.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: ui
katib.kubeflow.org/component: ui
{{- with .Values.ui.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
UI selector labels - Kustomize compatible
*/}}
{{- define "katib.ui.selectorLabels" -}}
katib.kubeflow.org/component: ui
{{- end }}

{{/*
DB Manager labels
*/}}
{{- define "katib.dbManager.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: db-manager
katib.kubeflow.org/component: db-manager
{{- with .Values.dbManager.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
DB Manager selector labels - Kustomize compatible
*/}}
{{- define "katib.dbManager.selectorLabels" -}}
katib.kubeflow.org/component: db-manager
{{- end }}

{{/*
MySQL labels
*/}}
{{- define "katib.mysql.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: mysql
katib.kubeflow.org/component: mysql
{{- end }}

{{/*
MySQL selector labels - Kustomize compatible
*/}}
{{- define "katib.mysql.selectorLabels" -}}
katib.kubeflow.org/component: mysql
{{- end }}

{{/*
PostgreSQL labels
*/}}
{{- define "katib.postgres.labels" -}}
{{ include "katib.labels" . }}
app.kubernetes.io/component: postgres
katib.kubeflow.org/component: postgres
{{- end }}

{{/*
PostgreSQL selector labels - Kustomize compatible
*/}}
{{- define "katib.postgres.selectorLabels" -}}
katib.kubeflow.org/component: postgres
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
{{ include "katib.controller.serviceName" . }}
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
{{ include "katib.mysql.serviceName" . }}
{{- else -}}
{{ .Values.database.external.host }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.enabled -}}
{{ include "katib.postgres.serviceName" . }}
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
Database secret name helper - matches Kustomize
*/}}
{{- define "katib.database.secretName" -}}
{{- if eq .Values.database.type "mysql" -}}
{{- if .Values.database.mysql.auth.existingSecret -}}
{{ .Values.database.mysql.auth.existingSecret }}
{{- else -}}
{{ include "katib.mysql.secretName" . }}
{{- end -}}
{{- else if eq .Values.database.type "postgres" -}}
{{- if .Values.database.postgres.auth.existingSecret -}}
{{ .Values.database.postgres.auth.existingSecret }}
{{- else -}}
katib-postgres-secrets
{{- end -}}
{{- else -}}
{{- if .Values.database.external.existingSecret -}}
{{ .Values.database.external.existingSecret }}
{{- else -}}
katib-mysql-secrets
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
  value: "katib"
{{- else }}
- name: DB_NAME
  value: {{ .Values.database.external.type | default "mysql" | quote }}
- name: DB_USER
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: DB_USER
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: DB_PASSWORD
- name: KATIB_MYSQL_DB_DATABASE
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: KATIB_MYSQL_DB_DATABASE
- name: KATIB_MYSQL_DB_HOST
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: KATIB_MYSQL_DB_HOST
- name: KATIB_MYSQL_DB_PORT
  valueFrom:
    secretKeyRef:
      name: {{ include "katib.database.secretName" . }}
      key: KATIB_MYSQL_DB_PORT
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

{{/*
Image pull policy helper
*/}}
{{- define "katib.imagePullPolicy" -}}
{{- $policy := .pullPolicy | default .Values.global.imagePullPolicy -}}
{{- if and .Values.global.kustomizeMode.omitDefaultImagePullPolicy (eq $policy "IfNotPresent") -}}
{{- else -}}
imagePullPolicy: {{ $policy }}
{{- end -}}
{{- end -}}

{{/*
Protocol helper
*/}}
{{- define "katib.protocol" -}}
{{- $protocol := .protocol | default "TCP" -}}
{{- if and .Values.global.kustomizeMode.omitDefaultProtocol (eq $protocol "TCP") -}}
{{- else -}}
protocol: {{ $protocol }}
{{- end -}}
{{- end -}}

{{/*
Service type helper
*/}}
{{- define "katib.serviceType" -}}
{{- $type := .type | default "ClusterIP" -}}
{{- if and .Values.global.kustomizeMode.omitDefaultServiceType (eq $type "ClusterIP") -}}
{{- else -}}
type: {{ $type }}
{{- end -}}
{{- end -}}

{{/*
Failure policy helper
*/}}
{{- define "katib.failurePolicy" -}}
{{- $policy := .failurePolicy | default "Fail" -}}
{{- if and .Values.global.kustomizeMode.omitDefaultFailurePolicy (eq $policy "Fail") -}}
{{- else -}}
failurePolicy: {{ $policy }}
{{- end -}}
{{- end -}}

{{/*
Controller labels for ClusterRole (conditionally omit component labels)
*/}}
{{- define "katib.controller.clusterRoleLabels" -}}
{{- if .Values.global.kustomizeMode.omitComponentLabels -}}
{{ include "katib.labels" . }}
{{- else -}}
{{ include "katib.controller.labels" . }}
{{- end -}}
{{- end -}}

{{/*
UI labels for ClusterRole (conditionally omit component labels)
*/}}
{{- define "katib.ui.clusterRoleLabels" -}}
{{- if .Values.global.kustomizeMode.omitComponentLabels -}}
{{ include "katib.labels" . }}
{{- else -}}
{{ include "katib.ui.labels" . }}
{{- end -}}
{{- end -}}

{{/*
Controller labels for ServiceAccount (conditionally omit component labels)
*/}}
{{- define "katib.controller.serviceAccountLabels" -}}
{{- if .Values.global.kustomizeMode.omitComponentLabels -}}
{{ include "katib.labels" . }}
{{- else -}}
{{ include "katib.controller.labels" . }}
{{- end -}}
{{- end -}}

{{/*
UI labels for ServiceAccount (conditionally omit component labels)
*/}}
{{- define "katib.ui.serviceAccountLabels" -}}
{{- if .Values.global.kustomizeMode.omitComponentLabels -}}
{{ include "katib.labels" . }}
{{- else -}}
{{ include "katib.ui.labels" . }}
{{- end -}}
{{- end -}}

{{/*
MySQL labels for PVC (conditionally omit component labels)
*/}}
{{- define "katib.mysql.pvcLabels" -}}
{{- if .Values.global.kustomizeMode.omitComponentLabels -}}
{{ include "katib.labels" . }}
{{- else -}}
{{ include "katib.mysql.labels" . }}
{{- end -}}
{{- end -}}

{{/*
Default security context helper
*/}}
{{- define "katib.defaultSecurityContext" -}}
runAsNonRoot: true
allowPrivilegeEscalation: false
runAsUser: 1000
seccompProfile:
  type: RuntimeDefault
capabilities:
  drop:
  - ALL
{{- end -}}

{{/*
Pod annotations helper with sidecar injection
*/}}
{{- define "katib.podAnnotations" -}}
{{- $annotations := dict -}}
{{- if .sidecarInject -}}
{{- $_ := set $annotations "sidecar.istio.io/inject" (.sidecarInject | quote) -}}
{{- end -}}
{{- if .prometheusAnnotations -}}
{{- $_ := set $annotations "prometheus.io/scrape" "true" -}}
{{- $_ := set $annotations "prometheus.io/port" (.prometheusAnnotations.port | quote) -}}
{{- if .prometheusAnnotations.scheme -}}
{{- $_ := set $annotations "prometheus.io/scheme" .prometheusAnnotations.scheme -}}
{{- end -}}
{{- end -}}
{{- if .customAnnotations -}}
{{- range $key, $value := .customAnnotations -}}
{{- $_ := set $annotations $key $value -}}
{{- end -}}
{{- end -}}
{{- if $annotations -}}
{{- range $key, $value := $annotations }}
{{ $key }}: {{ $value }}
{{- end -}}
{{- end -}}
{{- end -}}