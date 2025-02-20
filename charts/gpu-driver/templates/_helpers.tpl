
{{- define "gpu-driver.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 53 | trimSuffix "-" -}}-installer
{{- end -}}

{{- define "device-plugin.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 53 | trimSuffix "-" -}}-device-plugin
{{- end -}}

{{- define "node-labeler.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 53 | trimSuffix "-" -}}-node-labeler
{{- end -}}

{{- define "rbac.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 53 | trimSuffix "-" -}}
{{- end -}}


{{- define "gardenlinux.version" -}}
{{- if .Values.kernelVersion -}}
{{- if not (hasKey .Values.kernelVersions .Values.kernelVersion) }}{{ fail (printf "Unknown kernel version '%s'" .Values.kernelVersion) }}{{ end -}}
{{- get .Values.kernelVersions .Values.kernelVersion -}}
{{- else if .Values.gardenlinux.versionOverride -}}
{{- else -}}
{{- fail ".kernelVersion or .gardenlinux.versionOverride must be set" -}}
{{- end -}}
{{- end -}}


{{- define "image-pull-secrets" -}}
{{- with .Values.imagePullSecrets }}
{{- toYaml . | nindent 8 }}
{{- end }}
{{- if empty .Values.imagePullSecrets -}}
[]
{{- end }}
{{- end -}}

{{- define "node-selector" -}}
{{- $dict := .Values.nodeSelector }}
{{- if .Values.nodePool -}}
{{- $dict = set $dict "worker.gardener.cloud/pool" .Values.nodePool -}}
{{- end -}}
{{- if $dict -}}
      nodeSelector:
{{ toYaml $dict | indent 8 }}
{{- end -}}
{{- end -}}
