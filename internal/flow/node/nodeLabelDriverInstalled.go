package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow"
)

func nodeLabelDriverInstalled(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	_, hasLabel := state.ObjAsNode().Labels[flow.LabelDriverInstalled]
	if hasLabel {
		return ctx, nil
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		flow.LabelDriverInstalled: "true",
	})
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error patch merging node with driver installed label", composed.StopWithRequeue, ctx)
	}

	if changed {
		logger := composed.LoggerFromCtx(ctx)
		logger.Info("GPU driver installed")
		k8s.Event(ctx, state.ObjAsNode(), "Normal", "GpuDriverInstalled", "GPU driver installed")
	}

	return ctx, nil
}
