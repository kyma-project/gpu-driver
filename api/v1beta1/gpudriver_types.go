/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GpuDriverInstaller struct {
	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// +kubebuilder:default="IfNotPresent"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +kubebuilder:default="ghcr.io"
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:default="gardenlinux/gardenlinux/kmodbuild"
	Image string `json:"image,omitempty"`
}

//type GpuDriverFabricManager struct {
//	// +kubebuilder:default=true
//	Enabled bool `json:"enabled,omitempty"`
//
//	// +optional
//	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
//
//	// +kubebuilder:default="IfNotPresent"
//	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
//
//	// +kubebuilder:default="ghcr.io"
//	Repository string `json:"repository,omitempty"`
//
//	// +kubebuilder:default="gardenlinux/gardenlinux/kmodbuild"
//	Image string `json:"image,omitempty"`
//
//	// +kubebuilder:default="bookworm"
//	Version string `json:"version,omitempty"`
//}

type GpuDriverDevicePlugin struct {
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// +kubebuilder:default="IfNotPresent"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +kubebuilder:default="gcr.io"
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:default="gke-release/nvidia-gpu-device-plugin"
	Image string `json:"image,omitempty"`

	// +kubebuilder:default="1.0.25-gke.56"
	Version string `json:"version,omitempty"`
}

// GpuDriverSpec defines the desired state of GpuDriver.
type GpuDriverSpec struct {
	// +kubebuilder:validation:Required
	NodeSelector map[string]string `json:"nodeSelector"`

	// +optional
	DriverVersion string `json:"driverVersion,omitempty"`

	// +optional
	Installer GpuDriverInstaller `json:"installer"`

	//// +optional
	//FabricManager GpuDriverFabricManager `json:"fabricManager"`

	// +optional
	DevicePlugin GpuDriverDevicePlugin `json:"devicePlugin"`
}

// GpuDriverStatus defines the observed state of GpuDriver.
type GpuDriverStatus struct {
	State string `json:"state,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// GpuDriver is the Schema for the gpudrivers API.
type GpuDriver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GpuDriverSpec   `json:"spec,omitempty"`
	Status GpuDriverStatus `json:"status,omitempty"`
}

func (in *GpuDriver) InstallerImage() string {
	repo := in.Spec.Installer.Repository
	if repo == "" {
		repo = "ghcr.io"
	}
	img := in.Spec.Installer.Image
	if img == "" {
		img = "gardenlinux/gardenlinux/kmodbuild"
	}
	return fmt.Sprintf("%s/%s", repo, img)
}

func (in *GpuDriver) DevicePluginImage() string {
	repo := in.Spec.DevicePlugin.Repository
	if repo == "" {
		repo = "gcr.io"
	}
	img := in.Spec.DevicePlugin.Image
	if img == "" {
		img = "gke-release/nvidia-gpu-device-plugin"
	}
	version := in.Spec.DevicePlugin.Version
	if version == "" {
		version = "1.0.25-gke.56"
	}
	return fmt.Sprintf("%s/%s:%s", repo, img, version)
}

// +kubebuilder:object:root=true

// GpuDriverList contains a list of GpuDriver.
type GpuDriverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GpuDriver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GpuDriver{}, &GpuDriverList{})
}
