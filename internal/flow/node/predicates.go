package node

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
)

func nodeHasLabeledPredicate(label string) composed.Predicate {
	return func(ctx context.Context) bool {
		state := composed.StateFromCtx[*State](ctx)
		val, ok := state.ObjAsNode().Labels[label]
		if !ok || val == "" {
			return false
		}

		return true
	}
}
