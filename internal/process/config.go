package process

import "github.com/kyma-project/gpu-driver/internal/config"

type Config struct {
	Namespace               string `json:"namespace"`
	ConfigMapName           string `json:"configMapName"`
	InstallerServiceAccount string `json:"installerServiceAccount"`

	DriverVersions []DriverVersionConfig `json:"driverVersions"`

	DefaultDriverVersion string `json:"defaultDriverVersion"`

	KernelVersions   map[string]string `json:"kernelVersions"`
	CompileImage     string            `yaml:"compileImage"`
	ImagePullSecrets []string          `yaml:"imagePullSecrets"`
}

type DriverVersionConfig struct {
	DriverVersion string            `yaml:"driverVersion"`
	NodeSelector  map[string]string `yaml:"nodeSelector"`
}

var ProcessConfig = &Config{}

func InitConfig(cfg config.Config) {
	cfg.Path(
		"gpu.process",
		config.Bind(ProcessConfig),
		config.SourceFile("config.yaml"),
		config.Path(
			"namespace",
			config.DefaultScalar("default"),
		),
		config.Path(
			"installerServiceAccount",
			config.DefaultScalar("gpu-driver-installer"),
		),
		config.Path(
			"configMapName",
			config.DefaultScalar("gpu-driver-scripts"),
		),
		config.Path(
			"defaultDriverVersion",
			config.DefaultScalar("550.127.08"),
		),
		config.Path(
			"compileImage",
			config.DefaultScalar("ghcr.io/gardenlinux/gardenlinux/kmodbuild"),
		),
		config.Path(
			"kernelVersions",
			config.DefaultObj(map[string]string{
				"6.6.62-cloud-amd64": "amd64-1592.3",
				"6.6.63-cloud-amd64": "amd64-1592.4",
				"6.6.71-cloud-amd64": "amd64-1592.5",
				"6.6.78-cloud-amd64": "amd64-1592.6",
			}),
		),
		config.Path(
			"namespace",
			config.SourceEnv("NAMESPACE"),
		),
		config.Path(
			"configMapName",
			config.SourceEnv("CONFIGMAP_NAME"),
		),
	)
}
