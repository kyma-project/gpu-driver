package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func stopIfNodeDeletingOrUnschedulable(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if composed.IsMarkedForDeletion(state.ObjAsNode()) {
		return ctx, composed.StopAndForget
	}
	if state.ObjAsNode().Spec.Unschedulable {
		return ctx, composed.StopAndForget
	}

	return ctx, nil
}
