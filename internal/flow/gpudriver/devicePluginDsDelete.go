package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func devicePluginDsDelete(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.DevicePluginDS == nil {
		return ctx, nil
	}

	logger := composed.LoggerFromCtx(ctx)
	logger.Info("Deleting device plugin daemonset")

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	err := k8s.Delete(ctx, state.DevicePluginDS)
	if client.IgnoreNotFound(err) != nil {
		return composed.LogErrorAndReturn(err, "Error deleting device plugin daemonset", composed.StopWithRequeue, ctx)
	}

	k8s.Event(ctx, state.ObjAsGpuDriver(), "Normal", "DevicePluginDeleted", "Device plugin deleted")

	return ctx, nil
}
