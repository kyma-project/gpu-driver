package node

import (
	"context"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/kyma-project/gpu-driver/internal/common/composed"
	"github.com/kyma-project/gpu-driver/internal/common/k8sport"
	"github.com/kyma-project/gpu-driver/internal/config"
	"github.com/kyma-project/gpu-driver/internal/flow"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func driverPodCreate(ctx context.Context) (context.Context, error) {
	state := composed.StateFromCtx[*State](ctx)

	if state.Pod != nil {
		return ctx, nil
	}

	tag := config.KernelToKmodbuild(state.KernelVersion)
	if tag == "" {
		k8s := k8sport.FromCtxDefaultCluster(ctx)
		k8s.Event(ctx, state.ObjAsNode(), "Warning", "UnknownKernelVersion", fmt.Sprintf("Unknown kernel version '%s", state.ObjAsNode().Status.NodeInfo.KernelVersion))

		return composed.LogErrorAndReturn(
			fmt.Errorf("unknown kernel version '%s", state.ObjAsNode().Status.NodeInfo.KernelVersion),
			"Configuration error",
			composed.StopAndForget,
			ctx,
		)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: config.GetNamespace(),
			Name:      state.JobName(),
			Labels: map[string]string{
				flow.LabelCompiler: "true",
			},
		},
		Spec: corev1.PodSpec{
			PriorityClassName: "system-node-critical",
			HostPID:           true,
			ImagePullSecrets: pie.Map(state.GpuDriverConfig.Spec.Installer.ImagePullSecrets, func(s string) corev1.LocalObjectReference {
				return corev1.LocalObjectReference{Name: s}
			}),
			NodeName:      state.ObjAsNode().Name,
			RestartPolicy: corev1.RestartPolicyOnFailure,
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
			Volumes: []corev1.Volume{
				{
					Name: "scripts",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: config.GetScriptsConfigMapName(),
							},
							DefaultMode: ptr.To(int32(0744)),
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
					Name: "ld-root",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/",
						},
					},
				},
				{
					Name: "module-cache",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/opt/nvidia-installer/cache",
						},
					},
				},
				{
					Name: "module-install-dir-base",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/opt/drivers",
						},
					},
				},
				{
					Name: "usr-local",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/usr/local",
						},
					},
				},
			},
			ServiceAccountName: config.GetInstallerServiceAccountName(),
			Containers: []corev1.Container{
				{
					Name: "gpu-driver-installer",
					Image: fmt.Sprintf(
						"%s:%s",
						state.GpuDriverConfig.InstallerImage(),
						tag,
					),
					WorkingDir: "/work",
					Command:    []string{"/work/installer_entrypoint.sh"},
					SecurityContext: &corev1.SecurityContext{
						Privileged: ptr.To(true),
					},
					Env: []corev1.EnvVar{
						{
							Name:  "DEBUG",
							Value: "false",
						},
						{
							Name:  "NODE_NAME",
							Value: state.ObjAsNode().Name,
						},
						{
							Name:  "LABEL_GPU_NAME",
							Value: flow.LabelGpuName,
						},
						{
							Name:  "KERNEL_TYPE",
							Value: "cloud",
						},
						{
							Name:  "TARGET_ARCH",
							Value: "amd64",
						},
						{
							Name:  "KERNEL_NAME",
							Value: state.ObjAsNode().Status.NodeInfo.KernelVersion,
						},
						{
							Name:  "DRIVER_VERSION",
							Value: state.DriverVersion,
						},
						{
							Name:  "LD_ROOT",
							Value: "/root",
						},
						{
							Name:  "INSTALL_DIR",
							Value: "/opt/drivers",
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "scripts",
							MountPath: "/work",
						},
						{
							Name:      "dev",
							MountPath: "/dev",
						},
						{
							Name:      "ld-root",
							MountPath: "/root",
						},
						{
							Name:      "module-cache",
							MountPath: "/opt/nvidia-installer/cache",
						},
						{
							Name:      "module-install-dir-base",
							MountPath: "/opt/drivers",
						},
						{
							Name:      "usr-local",
							MountPath: "/usr/local",
						},
					},
				},
			},
		},
	}

	k8s := k8sport.FromCtxDefaultCluster(ctx)
	err := k8s.Create(ctx, pod)
	if apierrors.IsAlreadyExists(err) {
		return ctx, composed.StopWithRequeue
	}
	if err != nil {
		return composed.LogErrorAndReturn(err, "Error creating pod", composed.StopWithRequeue, ctx)
	}

	k8s.AnnotatedEventf(ctx, state.ObjAsNode(), map[string]string{
		"pod-name":      pod.Name,
		"pod-namespace": pod.Namespace,
	}, "Normal", "DriverIntallationPodStarted", "Driver installation pod is started")

	logger := composed.LoggerFromCtx(ctx)
	logger.Info("GPU install driver pod created")

	return ctx, composed.StopWithRequeue
}
