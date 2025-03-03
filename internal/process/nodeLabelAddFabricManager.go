package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func nodeLabelAddFabricManager(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	_, hasFabricManagerLabel := state.ObjAsNode().Labels[LabelFabricManager]
	if hasFabricManagerLabel {
		return ctx, nil
	}

	gpu := state.ObjAsNode().Labels[LabelGpuName]
	matchesRegex := fabricManagerGpuRegex.MatchString(gpu)

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	if !matchesRegex {
		return ctx, nil
	}

	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{LabelFabricManager: "true"})
	if apierrors.IsNotFound(err) {
		return ctx, composed.StopAndForget
	}
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error patch merge add node label fabric manager", composed.StopWithRequeue, ctx)
	}

	if changed {
		k8s.Event(ctx, state.ObjAsNode(), "Normal", "NodeMarkedForFabricManager", "Node marked for Fabric manager")
	}

	return ctx, nil
}
