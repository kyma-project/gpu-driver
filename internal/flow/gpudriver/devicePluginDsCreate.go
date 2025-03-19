package gpudriver

import (
	"context"
	"github.com/elliotchance/pie/v2"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	_ "github.com/kyma-project/gpu-driver/internal/config"
	"github.com/kyma-project/gpu-driver/internal/flow"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func devicePluginDsCreate(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.DevicePluginDS != nil {
		return ctx, nil
	}

	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: config.GetNamespace(),
			Name:      state.DevicePluginDSName(),
			Labels: map[string]string{
				flow.LabelGpuDriverConfig: state.ObjAsGpuDriver().Name,
				flow.LabelSignature:       state.ObjAsGpuDriver().DevicePluginHash(),
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					flow.LabelDevicePlugin:    "true",
					flow.LabelGpuDriverConfig: state.ObjAsGpuDriver().Name,
				},
			},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						flow.LabelDevicePlugin:    "true",
						flow.LabelGpuDriverConfig: state.ObjAsGpuDriver().Name,
					},
				},
				Spec: corev1.PodSpec{
					PriorityClassName: "system-node-critical",
					Volumes: []corev1.Volume{
						{
							Name: "device-plugin",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet/device-plugins",
								},
							},
						},
						{
							Name: "dev",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev",
								},
							},
						},
						{
							Name: "nvidia",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									// expecting that once driver is installed in
									// /opt/nvidia-installer/cache/nvidia/$DRIVER_VERSION
									// it is symlinked to /opt/drivers/nvidia as that was set in $INSTALL_DIR
									Path: "/opt/drivers/nvidia",
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Key:      "",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoExecute,
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: corev1.TolerationOpExists,
						},
					},
					ImagePullSecrets: pie.Map(state.ObjAsGpuDriver().Spec.DevicePlugin.ImagePullSecrets, func(x string) corev1.LocalObjectReference {
						return corev1.LocalObjectReference{Name: x}
					}),
					NodeSelector: state.ObjAsGpuDriver().Spec.NodeSelector,
					Containers: []corev1.Container{
						{
							Name:            "nvidia-gpu-device-plugin",
							Image:           state.ObjAsGpuDriver().DevicePluginImage(),
							ImagePullPolicy: state.ObjAsGpuDriver().Spec.DevicePlugin.ImagePullPolicy,
							Command: []string{
								"/usr/bin/nvidia-gpu-device-plugin",
								"-logtostderr",
								"-host-path=/opt/drivers/nvidia",
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: ptr.To(true),
							},
							Env: []corev1.EnvVar{
								{
									Name:  "LD_LIBRARY_PATH",
									Value: "/usr/local/nvidia/lib",
								},
								{
									Name:  "GOMAXPROCS",
									Value: "1",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "device-plugin",
									MountPath: "/device-plugin",
								},
								{
									Name:      "dev",
									MountPath: "/dev",
								},
								{
									Name:      "nvidia",
									MountPath: "/usr/local/nvidia",
								},
							},
						},
					},
				},
			},
		},
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)

	err := k8s.Create(ctx, ds)
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error creating device plugin daemonset", composed.StopWithRequeue, ctx)
	}

	logger := composed.LoggerFromCtx(ctx)
	logger.WithValues(
		"device-plugin-ds", ds.Name,
	).Info("Device plugin daemonset created")

	k8s.Event(ctx, state.ObjAsGpuDriver(), "Normal", "DevicePluginDaemonsetCreated", "Device plugin daemonset created")

	return ctx, nil
}
