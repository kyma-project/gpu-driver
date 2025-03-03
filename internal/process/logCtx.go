package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func logCtx(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)
	logger := composed.LoggerFromCtx(ctx)
	logger.WithValues(
		"nodeID", state.ID,
		"kernelVersion", state.ObjAsNode().Status.NodeInfo.KernelVersion,
	)

	ctx = composed.LoggerIntoCtx(ctx, logger)

	return ctx, nil
}
