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

	ctrl "sigs.k8s.io/controller-runtime"

	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/flow/gpudriver"
)

// GpuDriverOperatorReconciler reconciles a GpuDriverOperator object
type GpuDriverOperatorReconciler struct {
	Cluster composed.StateCluster
}

//+kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudriveroperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudriveroperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudriveroperators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GpuDriverOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *GpuDriverOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = composed.InitState(ctx, req.NamespacedName, &gpuv1beta1.GpuDriver{})
	ctx = composed.ClusterToCtx(ctx, r.Cluster)
	ctx = k8sport.ToCtx(ctx, k8sport.NewK8sPortOnDefaultCluster())
	act := gpudriver.New()
	return composed.Handle(act(ctx))
}

// SetupWithManager sets up the controller with the Manager.
func (r *GpuDriverOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1beta1.GpuDriverOperator{}).
		Complete(r)
}
