package fake

import (
	"context"
	"fmt"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sync"
)

type K8sFakePort interface {
	k8sport.K8sPort
	Find(name types.NamespacedName, obj client.Object) (client.Object, schema.GroupVersionKind, error)
}

var _ K8sFakePort = &k8sFakePort{}

type k8sFakePort struct {
	m sync.Mutex

	scheme  *runtime.Scheme
	objects map[string][]client.Object
}

func (f *k8sFakePort) Find(name types.NamespacedName, obj client.Object) (client.Object, schema.GroupVersionKind, error) {
	f.m.Lock()
	defer f.m.Unlock()
	return f.findNoLock(name, obj)
}

func (f *k8sFakePort) findNoLock(name types.NamespacedName, obj client.Object) (client.Object, schema.GroupVersionKind, error) {
	gvk, err := apiutil.GVKForObject(obj, f.scheme)
	if err != nil {
		return nil, schema.GroupVersionKind{}, err
	}
	key := fmt.Sprintf("%s/%s/%s", gvk.Group, gvk.Kind, name)
	arr := f.objects[key]
	for _, o := range arr {
		if o.GetName() == name.Name && o.GetNamespace() == name.Namespace {
			return o, gvk, nil
		}
	}
	return nil, gvk, nil
}

func (f *k8sFakePort) LoadStateObj(ctx context.Context) error {
	state := composed.StateFromCtx[composed.State](ctx)
	return f.LoadObj(ctx, state.Name(), state.Obj())
}

func (f *k8sFakePort) PatchMergeLabels(ctx context.Context, obj client.Object, labels map[string]string) (bool, error) {
	f.m.Lock()
	defer f.m.Unlock()
	if isContextCanceled(ctx) {
		return false, context.Canceled
	}

	o, gvk, err := f.findNoLock(client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return false, err
	}
	if o == nil {
		return false, apierrors.NewNotFound(schema.GroupResource{
			Group:    gvk.Group,
			Resource: gvk.Kind,
		}, obj.GetName())
	}

	changed := false
	if o.GetLabels() == nil {
		o.SetLabels(map[string]string{})
	}
	for k, v := range labels {
		if o.GetLabels()[k] != v {
			o.GetLabels()[k] = v
			changed = true
		}
	}

	return changed, nil
}

func (f *k8sFakePort) PatchMergeAnnotations(ctx context.Context, obj client.Object, annotations map[string]string) (bool, error) {
	f.m.Lock()
	defer f.m.Unlock()
	if isContextCanceled(ctx) {
		return false, context.Canceled
	}

	o, gvk, err := f.findNoLock(client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return false, err
	}
	if o == nil {
		return false, apierrors.NewNotFound(schema.GroupResource{
			Group:    gvk.Group,
			Resource: gvk.Kind,
		}, obj.GetName())
	}

	changed := false
	if o.GetAnnotations() == nil {
		o.SetAnnotations(map[string]string{})
	}
	for k, v := range annotations {
		if o.GetAnnotations()[k] != v {
			o.GetAnnotations()[k] = v
			changed = true
		}
	}

	return changed, nil
}

func (f *k8sFakePort) LoadObj(ctx context.Context, name types.NamespacedName, obj client.Object) error {
	f.m.Lock()
	defer f.m.Unlock()
	if isContextCanceled(ctx) {
		return context.Canceled
	}

	o, gvk, err := f.findNoLock(name, obj)
	if err != nil {
		return err
	}

	if o != nil {
		obj = o.DeepCopyObject().(client.Object)
		return nil
	}

	return apierrors.NewNotFound(schema.GroupResource{
		Group:    gvk.Group,
		Resource: gvk.Kind,
	}, name.Name)
}
