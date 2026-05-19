package node

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow"
)

var (
	re = regexp.MustCompile(`Garden Linux ([a-zA-Z0-9.]+)$`)
)

func nodeDetectKernelVersion(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	state.KernelVersion = state.ObjAsNode().Status.NodeInfo.KernelVersion
	state.ID = state.ObjAsNode().Labels[flow.LabelId]
	if state.ID == "" {
		state.ID = uuid.NewString()
	}

	//Parse the OS version.
	match := re.FindStringSubmatch(state.ObjAsNode().Status.NodeInfo.OSImage)
	if len(match) > 1 {
		state.OsImageVersion = match[1]
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	// When the kernel version changes, the previously installed driver and its
	// resolved version are no longer valid. Clear both labels so the operator
	// re-resolves the correct driver version for the new kernel and reinstalls.
	previousKernel := state.ObjAsNode().Labels[flow.LabelKernelVersion]
	if previousKernel != "" && previousKernel != state.KernelVersion {
		if err := k8s.PatchDeleteLabels(ctx, state.ObjAsNode(), []string{
			flow.LabelDriverInstalled,
			flow.LabelDriverVersion,
		}); err != nil {
			return composed.LogErrorAndReturn(err, "Error clearing stale driver labels on kernel change", composed.StopWithRequeue, ctx)
		}
		k8s.Event(ctx, state.ObjAsNode(), "Normal", "KernelVersionChanged",
			fmt.Sprintf("Kernel changed from %s to %s, driver will be reinstalled", previousKernel, state.KernelVersion))
	}

	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		flow.LabelKernelVersion: state.KernelVersion,
		flow.LabelId:            state.ID,
	})
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error patch merging kernel version label to node", composed.StopWithRequeue, ctx)
	}

	logger := composed.LoggerFromCtx(ctx)
	logger = logger.WithValues(
		"nodeID", state.ID,
		"kernelVersion", state.KernelVersion,
	)

	if changed {
		k8s.Event(ctx, state.ObjAsNode(), "Normal", "NodeKernelVersionDetected", "Nade kernel version is detected")
	}

	return ctx, nil
}
