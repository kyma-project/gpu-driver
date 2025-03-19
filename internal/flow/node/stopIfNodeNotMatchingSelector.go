package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"k8s.io/apimachinery/pkg/labels"
)

func stopIfNodeNotMatchingSelector(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	sel := labels.SelectorFromSet(state.GpuDriverConfig.Spec.NodeSelector)

	if sel.Matches(labels.Set(state.ObjAsNode().Labels)) {
		return ctx, nil
	}

	return composed.ComposeActions(
		driverPodLoad,
		driverPodDelete,
		composed.StopAndForgetAction,
	)(ctx)
}
