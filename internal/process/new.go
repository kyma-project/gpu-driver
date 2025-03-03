package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/actions"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func New() composed.Action {
	return func(ctx context.Context) (context.Context, error) {
		ctx = NewStateToCtx(ctx)

		return composed.ComposeActions(
			actions.LoadObjStopIfNotFound,
			stopIfNodeNotSelected,
			stopIfNodeDeletingOrUnschedulable,
			nodeLabelWithKernelVersion,
			logCtx,
			composed.If(
				composed.Not(nodeHasLabeledPredicate(LabelDriverInstalled)),
				driverPodLoad,
				driverPodCreate,
				driverPodWaitComplete,
				nodeLabelWithDriverVersion,
				driverPodDelete,
			),
			composed.If(
				nodeHasLabeledPredicate(LabelGpuName),
				nodeLabelAddFabricManager,
			),
		)(ctx)
	}
}
