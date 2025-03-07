package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// nodeNock annotates all matching nodes so they get changed to trigger their new reconciliation loop
// in order to notify them there's a new GpuConfig affecting them. If a node already has the annotation
// it is not changed
func nodeNock(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	var opts []client.ListOption
	if len(state.ObjAsGpuDriver().Spec.NodeSelector) > 0 {
		opts = append(opts, client.MatchingLabels(state.ObjAsGpuDriver().Spec.NodeSelector))
	}

	list := corev1.NodeList{}
	if err := k8s.List(ctx, &list, opts...); err != nil {
		return composed.LogErrorAndReturn(err, "Error listing nodes for GpuDriver nock", composed.StopWithRequeue, ctx)
	}

	for _, node := range list.Items {
		if _, ok := node.Annotations[flow.AnnotationNodeNock]; ok {
			continue
		}

		_, err := k8s.PatchMergeAnnotations(ctx, &node, map[string]string{flow.AnnotationNodeNock: node.Name})
		if client.IgnoreNotFound(err) != nil {
			return composed.LogErrorAndReturn(err, "Error patching node nock annotation", composed.StopWithRequeue, ctx)
		}
	}

	return ctx, nil
}
