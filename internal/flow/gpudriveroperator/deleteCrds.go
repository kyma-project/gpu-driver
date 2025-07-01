package gpudriveroperator

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/util"
)

func deleteCrds(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)
	logger := composed.LoggerFromCtx(ctx)

	if state.ObjAsGpuDriverOperator().Status.Served != v1beta1.ServedTrue {
		return nil, nil
	}
	k8s := k8sport.FromCtxDefaultCluster(ctx)

	list := util.NewCrdListUnstructured()
	err := k8s.List(ctx, list)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error listing CRDs", composed.StopWithRequeue, ctx)
	}

	logger.Info("Checking CRDs to uninstall")

	suffix := ".gpu.kyma-project.io"
	skip := "gpudriveroperator.gpu.kyma-project.io"
	for _, crd := range list.Items {
		if !strings.HasSuffix(crd.GetName(), suffix) {
			continue
		}
		if crd.GetName() == skip {
			continue
		}

		logger.
			WithValues("crd", crd.GetName()).
			Info("Deleting CRD")

		u := util.NewCrdUnstructured()
		u.SetName(crd.GetName())
		err = k8s.Delete(ctx, u)
		if err != nil {
			return composed.LogErrorAndReturn(err, fmt.Sprintf("Error deleting CRD %s", crd.GetName()), composed.StopWithRequeue, ctx)
		}
	}

	return nil, nil
}
