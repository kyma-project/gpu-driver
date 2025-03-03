package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func driverPodDelete(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.Pod == nil {
		return ctx, nil
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	err := k8s.Delete(ctx, state.Pod)
	if client.IgnoreNotFound(err) != nil {
		return composed.LogErrorAndReturn(err, "Error deleting driver install pod", composed.StopWithRequeue, ctx)
	}

	return ctx, nil
}
