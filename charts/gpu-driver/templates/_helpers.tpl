
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


{{- define "image-pull-secrets" -}}
{{- with .Values.imagePullSecrets }}
{{- toYaml . | nindent 8 }}
{{- end }}
{{- if empty .Values.imagePullSecrets -}}
[]
{{- end }}
{{- end -}}

{{- define "node-selector.no-kernel" -}}
{{- $dict := merge .Values.nodeSelector dict }}
{{- if .Values.nodePool -}}
{{- $dict = set $dict "worker.gardener.cloud/pool" .Values.nodePool -}}
{{- end -}}
{{- $dict = unset $dict "gpu.kyma-project.io/kernel-version" -}}
{{- if $dict -}}
      nodeSelector:
{{ toYaml $dict | indent 8 }}
{{- else }}
# no node selector
{{- end -}}
{{- end -}}


{{- define "node-selector.with-kernel" -}}
{{- $dict := merge .Values.nodeSelector dict }}
{{- if .Values.nodePool -}}
{{- $dict = set $dict "worker.gardener.cloud/pool" .Values.nodePool -}}
{{- end -}}
{{- if .Values.kernel -}}
{{- $dict = set $dict "gpu.kyma-project.io/kernel-version" .Values.kernel -}}
{{- end -}}
{{- if $dict -}}
      nodeSelector:
{{ toYaml $dict | indent 8 }}
{{- else }}
# no node selector
{{- end -}}
{{- end -}}
