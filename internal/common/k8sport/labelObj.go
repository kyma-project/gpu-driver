package k8sport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sLabelObjPort interface {
	PatchMergeLabels(ctx context.Context, obj client.Object, labels map[string]string) (bool, error)
	PatchDeleteLabels(ctx context.Context, obj client.Object, labelKeys []string) error
}

func NewK8sLabelObjPort(clusterID string) K8sLabelObjPort {
	return &k8sLabelObjPortImpl{clusterID: clusterID}
}

func NewK8sLabelObjPortOnDefaultCluster() K8sLabelObjPort {
	return NewK8sLabelObjPort(composed.DefaultClusterID)
}

type k8sLabelObjPortImpl struct {
	clusterID string
}

func (p *k8sLabelObjPortImpl) PatchDeleteLabels(ctx context.Context, obj client.Object, labelKeys []string) error {
	nullLabels := make(map[string]interface{}, len(labelKeys))
	for _, k := range labelKeys {
		nullLabels[k] = nil
	}
	data := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": nullLabels,
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal label delete patch: %w", err)
	}
	cluster := composed.ClusterFromCtx(ctx, p.clusterID)
	if err := cluster.K8sClient().Patch(ctx, obj, client.RawPatch(types.MergePatchType, b)); err != nil {
		return err
	}
	for _, k := range labelKeys {
		delete(obj.GetLabels(), k)
	}
	return nil
}

func (p *k8sLabelObjPortImpl) PatchMergeLabels(ctx context.Context, obj client.Object, labels map[string]string) (bool, error) {
	if obj.GetLabels() == nil {
		obj.SetLabels(make(map[string]string))
	}
	payload := map[string]string{}
	for k, v := range labels {
		if obj.GetLabels()[k] != v {
			payload[k] = v
		}
	}
	if len(payload) == 0 {
		return false, nil
	}
	data := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": payload,
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("failed to marshal labels for merge patch: %w", err)
	}
	cluster := composed.ClusterFromCtx(ctx, p.clusterID)
	err = cluster.K8sClient().Patch(ctx, obj, client.RawPatch(types.MergePatchType, b))
	if err == nil {
		for k, v := range labels {
			obj.GetLabels()[k] = v
		}
	}
	return true, err
}
