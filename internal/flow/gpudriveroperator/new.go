package gpudriveroperator

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
			composed.IfElse(
				composed.MarkForDeletionPredicate,
				// GpuDriverOperator deletion
				composed.ComposeActions(
					deleteCrds,
				),
				// GpuDriverOperator - reconcile served
				composed.ComposeActions(
					handleServed,
					statusReady,
				),
			),
		)(ctx)
	}
}
