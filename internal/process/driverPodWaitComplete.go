package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/util"
	corev1 "k8s.io/api/core/v1"
)

func driverPodWaitComplete(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.Pod == nil {
		return ctx, composed.StopWithRequeue
	}

	completed := false
	if state.Pod.Status.Phase == corev1.PodSucceeded || state.Pod.Status.Phase == corev1.PodFailed {
		completed = true
	}

	if !completed {
		return ctx, composed.StopWithRequeueDelay(util.Timing.T10000ms())
	}

	logger := composed.LoggerFromCtx(ctx)
	logger.Info("Driver install pod has completed")

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	k8s.AnnotatedEventf(ctx, state.ObjAsNode(), map[string]string{
		"pod-name":      state.Pod.Name,
		"pod-namespace": state.Pod.Namespace,
	}, "Normal", "DriverIntallationPodCompleted", "Driver installation pod is completed")

	return ctx, nil
}
