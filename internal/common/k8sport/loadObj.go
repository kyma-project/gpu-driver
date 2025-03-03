package k8sport

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sLoadPort interface {
	LoadStateObj(ctx context.Context) error
	LoadObj(ctx context.Context, name types.NamespacedName, obj client.Object) error
}

func NewK8sLoadPort(clusterID string) K8sLoadPort {
	return &k8sLoadPortImpl{clusterID: clusterID}
}

func NewK8sLoadPortOnDefaultCluster() K8sLoadPort {
	return NewK8sLoadPort(composed.DefaultClusterID)
}

type k8sLoadPortImpl struct {
	clusterID string
}

func (p *k8sLoadPortImpl) LoadStateObj(ctx context.Context) error {
	state := composed.StateFromCtx[composed.State](ctx)
	return p.LoadObj(ctx, state.Name(), state.Obj())
}

func (p *k8sLoadPortImpl) LoadObj(ctx context.Context, name types.NamespacedName, obj client.Object) error {
	cluster := composed.ClusterFromCtx(ctx, p.clusterID)
	return cluster.K8sClient().Get(ctx, name, obj)
}
