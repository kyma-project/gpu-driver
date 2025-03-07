package actions

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func StopIfNotLoaded(ctx context.Context, st composed.State) (context.Context, error) {
	if composed.IsLoaded(ctx) {
		return ctx, nil
	}

	return ctx, composed.StopAndForget
}
