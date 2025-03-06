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
	"github.com/kyma-project/gpu-driver/internal/config"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
)

// GpuDriverReconciler reconciles a GpuDriver object
type GpuDriverReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudrivers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudrivers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gpu.kyma-project.io,resources=gpudrivers/finalizers,verbs=update

func (r *GpuDriverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	gpuDriver := &gpuv1beta1.GpuDriver{}
	err := r.Get(ctx, req.NamespacedName, gpuDriver)
	if apierrors.IsNotFound(err) {
		config.Remove(req.Name)
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	config.Sync(gpuDriver)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GpuDriverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1beta1.GpuDriver{}).
		Named("gpudriver").
		Complete(r)
}
