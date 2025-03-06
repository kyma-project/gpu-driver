package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func nodeLabelWithDriverVersion(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		flow.LabelDriverVersion: state.DriverVersion,
	})
	if client.IgnoreNotFound(err) != nil {
		return composed.LogErrorAndReturn(err, "Error labeling node with installed driver version", composed.StopWithRequeue, ctx)
	}
	if !changed {
		return ctx, nil
	}

	logger := composed.LoggerFromCtx(ctx)
	logger.Info("GPU driver installed")

	k8s.AnnotatedEventf(
		ctx, state.ObjAsNode(), map[string]string{
			"driverVersion": state.DriverVersion,
		}, "Warning", "GpuDriverInstalled",
		"GPU Driver installed version '%s'", state.DriverVersion,
	)

	return ctx, nil
}
