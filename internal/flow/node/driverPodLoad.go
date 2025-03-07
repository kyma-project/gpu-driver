package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func driverPodLoad(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)
	k8s := k8sport.FromCtxDefaultCluster(ctx)

	pod := &corev1.Pod{}
	err := k8s.LoadObj(ctx, types.NamespacedName{
		Namespace: config.GetNamespace(),
		Name:      state.JobName(),
	}, pod)

	if client.IgnoreNotFound(err) != nil {
		return composed.LogErrorAndReturn(err, "Error loading pod", composed.StopWithRequeue, ctx)
	}

	if err == nil {
		state.Pod = pod
	}

	return ctx, nil
}
