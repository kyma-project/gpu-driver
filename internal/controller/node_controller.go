/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow"
	"github.com/kyma-project/gpu-driver/internal/flow/node"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeReconciler reconciles a Guestbook object
type NodeReconciler struct {
	Cluster composed.StateCluster
}

// +kubebuilder:rbac:groups=v1,resources=nodes,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=v1,resources=nodes/status,verbs=get
// +kubebuilder:rbac:groups=v1,resources=pods,verbs=get;list;watch;update;patch;create;delete
// +kubebuilder:rbac:groups=v1,resources=pods/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guestbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = composed.InitState(ctx, req.NamespacedName, &corev1.Node{})
	ctx = composed.ClusterToCtx(ctx, r.Cluster)
	ctx = k8sport.ToCtx(ctx, k8sport.NewK8sPortOnDefaultCluster())
	act := node.New()
	return composed.Handle(act(ctx))
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&corev1.Node{}).
		Named("driver").
		Watches(
			&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				pod, ok := obj.(*corev1.Pod)
				if !ok {
					return []reconcile.Request{}
				}
				if pod.Labels[flow.LabelCompiler] == "true" {
					return []reconcile.Request{}
				}
				return []reconcile.Request{{NamespacedName: types.NamespacedName{Name: pod.Spec.NodeName}}}
			}),
		).
		Complete(r)
}
