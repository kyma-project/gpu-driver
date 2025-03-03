package k8sport

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sCreatePort interface {
	Create(ctx context.Context, obj client.Object) error
}

func NewK8sCreatePort(clusterID string) K8sCreatePort {
	return &k8sCreateObjPortImpl{clusterID: clusterID}
}

func NewK8sCreatePortOnDefaultCluster() K8sCreatePort {
	return NewK8sCreatePort(composed.DefaultClusterID)
}

var _ K8sCreatePort = &k8sCreateObjPortImpl{}

type k8sCreateObjPortImpl struct {
	clusterID string
}

func (p *k8sCreateObjPortImpl) Create(ctx context.Context, obj client.Object) error {
	cluster := composed.ClusterFromCtx(ctx, p.clusterID)
	return cluster.K8sClient().Create(ctx, obj)
}
