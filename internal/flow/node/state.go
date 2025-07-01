package node

import (
	"context"
	"fmt"

	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func NewState(baseState composed.State) *State {
	return &State{
		State: baseState,
	}
}

func NewStateToCtx(ctx context.Context) context.Context {
	baseState := composed.StateFromCtx[composed.State](ctx)
	return composed.StateToCtx(ctx, NewState(baseState))
}

type State struct {
	composed.State

	KernelVersion      string
	AllMatchingConfigs []*gpuv1beta1.GpuDriver
	GpuDriverConfig    *gpuv1beta1.GpuDriver
	DriverVersion      string
	OsImageVersion     string

	ID  string
	Job *batchv1.Job
	Pod *corev1.Pod
}

func (s *State) ObjAsNode() *corev1.Node {
	return s.Obj().(*corev1.Node)
}

func (s *State) JobName() string {
	configName := fmt.Sprintf("%.*s", 200, s.GpuDriverConfig.Name)
	return fmt.Sprintf("installer-%s-%s", configName, s.ID)
}
