
{{- define "chart.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 51 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 51 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 51 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}


{{- define "chart.namespace" -}}
{{- if .Values.fullnameOverride -}}
{{- template "chart.fullname" }}-system
{{- else -}}
{{- .Release.Namespace -}}
{{- end -}}
{{- end -}}


{{- define "manager.image" -}}
{{- $image := .Values.manager.image -}}
{{- if .Values.manager.repository -}}
{{- $image = printf "%s/%s" .Values.manager.repository $image -}}
{{- end -}}
{{- $image }}:{{ required "manager.tag" .Values.manager.tag }}
{{- end -}}


{{- define "image-pull-secrets" -}}
{{- with .Values.imagePullSecrets }}
{{- toYaml . | nindent 8 }}
{{- end }}
{{- if (empty .Values.imagePullSecrets) -}}
[]
{{- end }}
{{- end -}}
