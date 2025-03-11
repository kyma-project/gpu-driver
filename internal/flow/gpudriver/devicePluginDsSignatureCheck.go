package gpudriver

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/flow"
	"github.com/kyma-project/gpu-driver/internal/util"
)

func devicePluginDsSignatureCheck(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.DevicePluginDS == nil {
		return ctx, nil
	}

	desiredSignature := state.ObjAsGpuDriver().DevicePluginHash()
	if state.DevicePluginDS.Labels[flow.LabelSignature] == desiredSignature {
		return ctx, nil
	}

	// signature mismatch, have to delete the old ds

	if composed.IsMarkedForDeletion(state.DevicePluginDS) {
		// ds is still being deleted, requeue until it's gone
		return ctx, composed.StopWithRequeueDelay(util.Timing.T1000ms())
	}

	var err error
	ctx, err = devicePluginDsDelete(ctx)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error deleting device plugin ds due to signature mismatch", composed.StopWithRequeue, ctx)
	}

	return ctx, composed.StopWithRequeueDelay(util.Timing.T1000ms())
}
