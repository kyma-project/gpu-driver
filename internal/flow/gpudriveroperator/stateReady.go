package gpudriveroperator

import (
	"context"

	"github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func statusReady(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)
	logger := composed.LoggerFromCtx(ctx)
	obj := state.ObjAsGpuDriverOperator()
	k8s := k8sport.FromCtxDefaultCluster(ctx)

	if composed.MarkForDeletionPredicate(ctx) {
		return nil, nil
	}

	condReady := meta.FindStatusCondition(obj.Status.Conditions, v1beta1.ConditionTypeReady)
	if condReady != nil && obj.Status.State == v1beta1.ModuleStateReady {
		return nil, nil
	}

	logger.Info("Setting Ready state and condition to CloudResources CR")
	obj.Status.State = v1beta1.ModuleStateReady
	obj.Status.Conditions = []metav1.Condition{
		{
			Type:    v1beta1.ConditionTypeReady,
			Status:  metav1.ConditionTrue,
			Reason:  v1beta1.ConditionReasonReady,
			Message: "Ready",
		},
	}

	err := k8s.UpdateStatus(ctx, obj)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error updating GpuDriverOperator with Ready condition", composed.StopWithRequeue, ctx)
	}

	return ctx, composed.StopAndForget

}
