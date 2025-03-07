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
			composed.If(
				shouldDeleteDaemonset,
				devicePluginDsDelete,
			),
			composed.If(
				shouldCreateDaemonset,
				nodeNock,
				devicePluginDsCreate,
			),
		)(ctx)
	}
}

func shouldCreateDaemonset(ctx context.Context) bool {
	return composed.All(
		composed.IsLoaded,
		composed.Not(composed.MarkForDeletionPredicate),
	)(ctx)
}

func shouldDeleteDaemonset(ctx context.Context) bool {
	return composed.Any(
		composed.Not(composed.IsLoaded),
		composed.MarkForDeletionPredicate,
	)(ctx)
}
