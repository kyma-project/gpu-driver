package process

import "regexp"

const (
	LabelKernelVersion   = "gpu.kyma-project.io/kernel-version"
	LabelId              = "gpu.kyma-project.io/id"
	LabelDriverInstalled = "gpu.kyma-project.io/driver-installed"
	LabelCompiler        = "gpu.kyma-project.io/compiler"
	LabelGpuName         = "gpu.kyma-project.io/gpu-name"
	LabelFabricManager   = "gpu.kyma-project.io/fabric-manager"
)

var (
	fabricManagerGpuRegex = regexp.MustCompile("(A100|H100|H200|B100|B200)")
)
