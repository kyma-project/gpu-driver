package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/actions"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func New() composed.Action {
	return func(ctx context.Context) (context.Context, error) {
		ctx = NewStateToCtx(ctx)

		return composed.ComposeActions(
			actions.LoadObj,
			devicePluginDsLoad,
			composed.IfElse(
				daemonsetShouldExist,
				// DevicePlugin should exist
				composed.ComposeActions(
					nodeNock,
					devicePluginDsSignatureCheck,
					devicePluginDsCreate,
				),
				// DevicePlugin should not exist
				composed.ComposeActions(
					devicePluginDsDelete,
				),
			),
		)(ctx)
	}
}

func daemonsetShouldExist(ctx context.Context) bool {
	return composed.All(
		composed.IsLoaded,
		composed.Not(composed.MarkForDeletionPredicate),
		func(ctx context.Context) bool {
			state := composed.StateFromCtx[*State](ctx)
			if state.Obj() != nil && state.Obj().GetName() != "" && state.Obj().GetGeneration() > 0 {
				return !state.ObjAsGpuDriver().Spec.DevicePlugin.Disabled
			}
			return false
		},
	)(ctx)
}
