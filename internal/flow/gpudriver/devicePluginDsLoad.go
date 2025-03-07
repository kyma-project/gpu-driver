package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func devicePluginDsLoad(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	ds := &appsv1.DaemonSet{}
	err := k8s.LoadObj(ctx, types.NamespacedName{
		Namespace: config.GetNamespace(),
		Name:      state.ObjAsGpuDriver().Name,
	}, ds)
	if client.IgnoreNotFound(err) != nil {
		return composed.LogErrorAndReturn(err, "Error loading device plugin daemonset", composed.StopWithRequeue, ctx)
	}

	if err == nil {
		state.DevicePluginDS = ds
	}

	return ctx, nil
}
