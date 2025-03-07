package flow

import "regexp"

const (
	LabelGpuDriverConfig = "gpu.kyma-project.io/gpu-driver-config"
	LabelKernelVersion   = "gpu.kyma-project.io/kernel-version"
	LabelId              = "gpu.kyma-project.io/id"
	LabelDriverInstalled = "gpu.kyma-project.io/driver-installed"
	LabelDriverVersion   = "gpu.kyma-project.io/driver-version"
	LabelCompiler        = "gpu.kyma-project.io/compiler"
	LabelGpuName         = "gpu.kyma-project.io/gpu-name"
	LabelDevicePlugin    = "gpu.kyma-project.io/device-plugin"
)

const (
	AnnotationNodeNock = "gpu.kyma-project.io/nock"
)

var (
	fabricManagerGpuRegex = regexp.MustCompile("(A100|H100|H200|B100|B200)")
)
