package gpudriveroperator

import (
	"context"
	"fmt"

	"github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleServed(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)
	logger := composed.LoggerFromCtx(ctx)
	obj := state.ObjAsGpuDriverOperator()
	k8s := k8sport.FromCtxDefaultCluster(ctx)

	served, err := func() (*v1beta1.GpuDriverOperator, error) {
		list := &v1beta1.GpuDriverOperatorList{}
		err := k8s.List(ctx, list)
		if err != nil {
			return nil, err
		}

		var servedObj *v1beta1.GpuDriverOperator
		for _, item := range list.Items {
			if item.Status.Served == v1beta1.ServedTrue {
				servedObj = &item
				break
			}
		}

		return servedObj, nil
	}()

	if err != nil {
		return composed.LogErrorAndReturn(err, "Error listing GpuDriverOperator", composed.StopWithRequeue, ctx)
	}

	if served != nil && obj.GetName() == served.Name && obj.GetNamespace() == served.Namespace {
		// nothing to do, the obj is the served one
		return ctx, nil
	}

	if obj.Status.Served == v1beta1.ServedFalse {
		logger.Info("Ignoring not served GpuDriverOperator")
		// we're not handling not served objects
		return ctx, composed.StopAndForget
	}

	if served == nil {

		err = k8s.UpdateStatus(ctx, obj)
		if err != nil {
			return composed.LogErrorAndReturn(err, "Error updating GpuDriverOperator served=true", composed.StopWithRequeue, ctx)
		} else {
			return ctx, composed.StopWithRequeue
		}
	}

	if served == nil {
		// none is served so far, this obj will be the one
		logger.Info("Setting GpuDriverOperator Served to True")
		obj.Status.Served = v1beta1.ServedTrue
	} else {
		logger.Info("Setting CloudResources Served to False")
		obj.Status.Served = v1beta1.ServedFalse
		obj.Status.State = v1beta1.ModuleStateWarning
		obj.Status.Conditions = []metav1.Condition{
			{
				Type:   v1beta1.ConditionTypeError,
				Status: metav1.ConditionTrue,
				Reason: v1beta1.ReasonOtherIsServed,
				Message: fmt.Sprintf("only one instance of GpuDriverOperator is allowed (current served instance: %s",
					served.Name),
			},
		}
	}
	err = k8s.UpdateStatus(ctx, obj)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error updating GpuDriverOperator served=true", composed.StopWithRequeue, ctx)
	}

	return ctx, composed.StopWithRequeue

}
