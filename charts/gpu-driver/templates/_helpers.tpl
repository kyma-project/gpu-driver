
{{- define "gpu-driver.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}-installer
{{- end -}}

{{- define "device-plugin.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 53 | trimSuffix "-" -}}-device-plugin
{{- end -}}

{{- define "image-pull-secrets" -}}
{{- with .Values.imagePullSecrets }}
{{- toYaml . | nindent 8 }}
{{- end }}
{{- if empty .Values.imagePullSecrets -}}
[]
{{- end }}
{{- end -}}

