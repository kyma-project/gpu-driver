package gpudriver

import (
	"context"
	"fmt"
	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	appsv1 "k8s.io/api/apps/v1"
)

type State struct {
	composed.State

	DevicePluginDS *appsv1.DaemonSet
}

func NewState(baseState composed.State) *State {
	return &State{
		State: baseState,
	}
}

func NewStateToCtx(ctx context.Context) context.Context {
	baseState := composed.StateFromCtx[composed.State](ctx)
	return composed.StateToCtx(ctx, NewState(baseState))
}

func (s *State) ObjAsGpuDriver() *gpuv1beta1.GpuDriver {
	return s.Obj().(*gpuv1beta1.GpuDriver)
}

func (s *State) DevicePluginDSName() string {
	return fmt.Sprintf("device-plugin-%s", s.ObjAsGpuDriver().Name)
}
