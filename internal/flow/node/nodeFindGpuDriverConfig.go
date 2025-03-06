package node

import (
	"context"
	"fmt"
	"github.com/elliotchance/pie/v2"
	gpuv1beta1 "github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	"github.com/kyma-project/gpu-driver/internal/flow"
	"strings"
)

func nodeFindGpuDriverConfig(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	allMatching := config.FindNodeConfigs(state.ObjAsNode())
	assignedGpuDriverName := state.ObjAsNode().Labels[flow.LabelGpuDriverConfig]

	if len(allMatching) == 0 {
		// ignore the node, it matches to none existing GpuDriver
		return ctx, composed.StopAndForget
	}

	// assume GpuDriver selectors do not overlap, and node matches to only one
	gpuDriver := allMatching[0]

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	if len(allMatching) > 1 && assignedGpuDriverName == "" {
		// The node matches to multiple GpuDriver configs,
		// but it has not already been assigned to some
		// * choose assumed first (maybe later change to oldest)
		// * emit an event
		matchingNames := pie.Map(allMatching, func(x *gpuv1beta1.GpuDriver) string {
			return x.Name
		})
		k8s.AnnotatedEventf(
			ctx, state.ObjAsNode(), map[string]string{
				"all-gpu-driver-configs":     strings.Join(matchingNames, ", "),
				"assigned-gpu-driver-config": gpuDriver.Name,
			}, "Warning", "OverlapingGpuDriverConfigs",
			fmt.Sprintf("Multiple GpuDriver configurations match single node"),
		)
	}
	if len(allMatching) > 1 && assignedGpuDriverName != "" {
		// The node matches to multiple GpuDriver configs,
		// but it already was assigned to some

		// Find if already assigned GpuDriver config still exists
		gpuDriver = nil
		for _, gpu := range allMatching {
			if gpu.Namespace == assignedGpuDriverName {
				gpuDriver = gpu
			}
		}

		if gpuDriver == nil {
			// Already assigned GpuConfig to the node does not exist anymore
			// * emit an event
			// * ignore the node
			k8s.AnnotatedEventf(
				ctx, state.ObjAsNode(), map[string]string{
					"assigned-gpu-driver-config": assignedGpuDriverName,
				}, "Warning", "GpuDriverConfigDoesNotExist",
				"GpuDriver config %s assigned to node no longer exists or matches to the node",
				assignedGpuDriverName,
			)
			return ctx, composed.StopAndForget
		}
	}

	// label the node with assigned GpuDriver config name
	changed, err := k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		flow.LabelGpuDriverConfig: gpuDriver.Name,
	})
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error labeling node with chosen GpuDriver config name", composed.StopWithRequeue, ctx)
	}

	if changed {
		k8s.AnnotatedEventf(ctx, state.ObjAsNode(), map[string]string{
			"assigned-gpu-driver-config": gpuDriver.Name,
		}, "Normal", "GpuDriverAssigned", "GpuDriver config %s assigned", gpuDriver.Name)
	}

	state.GpuDriverConfig = gpuDriver

	state.DriverVersion = gpuDriver.Spec.DriverVersion
	if state.DriverVersion == "" {
		state.DriverVersion = config.KernelToDriver(state.KernelVersion)
	}

	changed, err = k8s.PatchMergeLabels(ctx, state.ObjAsNode(), map[string]string{
		flow.LabelDriverVersion: state.DriverVersion,
	})

	logger := composed.LoggerFromCtx(ctx).WithValues(
		"driverVersion", state.DriverVersion,
		"gpuDriver", gpuDriver.Name,
	)
	ctx = composed.LoggerIntoCtx(ctx, logger)

	return ctx, nil
}
