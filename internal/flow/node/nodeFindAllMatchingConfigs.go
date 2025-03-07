package node

import (
	"context"
	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
)

func nodeFindAllMatchingConfigs(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	var list gpuv1beta1.GpuDriverList
	err := k8s.List(ctx, &list)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error listing all GpuDriver configs in node reconciler", composed.StopWithRequeue, ctx)
	}

	for _, x := range list.Items {
		state.AllMatchingConfigs = append(state.AllMatchingConfigs, &x)
	}

	return ctx, nil
}
