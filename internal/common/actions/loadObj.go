package actions

import (
	"context"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func LoadObj(ctx context.Context) (context.Context, error) {
	k8s := k8sport.FromCtxDefaultCluster(ctx)
	err := k8s.LoadStateObj(ctx)
	if client.IgnoreNotFound(err) != nil {
		return ctx, err
	}

	return ctx, nil
}

func LoadObjStopIfNotFound(ctx context.Context) (context.Context, error) {
	k8s := k8sport.FromCtxDefaultCluster(ctx)
	err := k8s.LoadStateObj(ctx)
	if client.IgnoreNotFound(err) != nil {
		return ctx, err
	}
	if apierrors.IsNotFound(err) {
		return ctx, composed.StopAndForget
	}

	return ctx, nil
}
