package process

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/util"
)

func stopIfNodeNotSelected(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if len(ProcessConfig.DriverVersions) == 0 {
		state.DesiredDriverVersion = ProcessConfig.DefaultDriverVersion
		return ctx, nil
	}

	for _, driverVersion := range ProcessConfig.DriverVersions {
		if util.MatchesLabels(state.ObjAsNode(), driverVersion.NodeSelector) {
			state.DesiredDriverVersion = driverVersion.DriverVersion

			logger := composed.LoggerFromCtx(ctx).WithValues("driverVersion", state.DesiredDriverVersion)
			ctx = composed.LoggerIntoCtx(ctx, logger)

			return ctx, nil
		}
	}

	return ctx, composed.StopAndForget
}
