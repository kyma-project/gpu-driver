package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	"github.com/kyma-project/gpu-driver/internal/flow"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func devicePluginDsLoad(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	list := &appsv1.DaemonSetList{}
	err := k8s.List(ctx, list, client.InNamespace(config.GetNamespace()), client.MatchingLabels{
		flow.LabelGpuDriverConfig: state.ObjAsGpuDriver().Name,
	})
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error listing device plugin daemonsets", composed.StopWithRequeue, ctx)
	}

	for _, ds := range list.Items {
		if ds.Name == state.DevicePluginDSName() {
			state.DevicePluginDS = &ds
		} else if !composed.IsMarkedForDeletion(&ds) {
			if err := k8s.Delete(ctx, &ds); client.IgnoreNotFound(err) != nil {
				ctx = composed.LoggerIntoCtx(ctx, composed.LoggerFromCtx(ctx).WithValues(
					"device-plugin-ds-being-deleted", ds.Name,
				))
				return composed.LogErrorAndReturn(err, "Error deleting device plugin daemonset with old name", composed.StopWithRequeue, ctx)
			}
		}
	}

	return ctx, nil
}
