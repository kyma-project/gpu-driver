package k8sport

import (
	"context"

	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sUpdateStatusPort interface {
	UpdateStatus(ctx context.Context, obj client.Object) error
}

func NewK8sUpdateStatusPort(clusterID string) K8sUpdateStatusPort {
	return &k8sUpdateStatusPortImpl{clusterID: clusterID}
}

func NewK8sUpdateStatusPortOnDefaultCluster() K8sUpdateStatusPort {
	return NewK8sUpdateStatusPort(composed.DefaultClusterID)
}

var _ K8sUpdateStatusPort = &k8sUpdateStatusPortImpl{}

type k8sUpdateStatusPortImpl struct {
	clusterID string
}

func (p *k8sUpdateStatusPortImpl) UpdateStatus(ctx context.Context, obj client.Object) error {
	cluster := composed.ClusterFromCtx(ctx, p.clusterID)
	return cluster.K8sClient().Status().Update(ctx, obj)
}
