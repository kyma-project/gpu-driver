package process

import (
	"context"
	"github.com/google/uuid"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
)

func nodeLabelWithKernelVersion(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	kernelVersion := state.ObjAsNode().Status.NodeInfo.KernelVersion
	id := state.ObjAsNode().Labels[LabelId]
	if id == "" {
		id = uuid.NewString()
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		LabelKernelVersion: kernelVersion,
		LabelId:            id,
	})
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error patch merging kernel version label to node", composed.StopWithRequeue, ctx)
	}

	if !changed {
		state.ID = state.ObjAsNode().Labels[LabelId]
		return ctx, nil
	}

	k8s.Event(ctx, state.ObjAsNode(), "Normal", "LabeledWithKernelVersion", "Nade has been labeled with kernel version")

	state.ID = id

	return ctx, nil
}
