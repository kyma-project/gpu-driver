package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/actions"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/flow"
)

func New() composed.Action {
	return func(ctx context.Context) (context.Context, error) {
		ctx = NewStateToCtx(ctx)

		return composed.ComposeActions(
			actions.LoadObjStopIfNotFound,
			stopIfNodeDeletingOrUnschedulable,
			nodeDetectKernelVersion,
			nodeFindAllMatchingConfigs,
			nodeFindGpuDriverConfig,
			stopIfNodeNotMatchingSelector,
			composed.If(
				composed.Not(nodeHasLabeledPredicate(flow.LabelDriverInstalled)),
				driverPodLoad,
				driverPodCreate,
				driverPodWaitComplete,
				driverPodDelete,
				nodeLabelDriverInstalled,
			),
		)(ctx)
	}
}
