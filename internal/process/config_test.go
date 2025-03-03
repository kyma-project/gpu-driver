package process

import (
	"github.com/kyma-project/gpu-driver/internal/config"
	"github.com/kyma-project/gpu-driver/internal/util"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("mixed file and env", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "cloud-manager-config")
		assert.NoError(t, err, "error creating tmp dir")
		defer func() {
			_ = os.RemoveAll(dir)
		}()
		err = os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(`
driverVersion: 10.20.30
nodeSelector:
  node-type: gpu
`), 0644)
		assert.NoError(t, err, "error creating key file")

		cfg := config.NewConfig(util.NewMockedEnvironment(map[string]string{
			"NAMESPACE":      "the-namespace",
			"CONFIGMAP_NAME": "cm-gpu",
		}))
		cfg.BaseDir(dir)
		InitConfig(cfg)
		cfg.Read()

		// from env
		assert.Equal(t, "cm-gpu", ProcessConfig.ConfigMapName)
		assert.Equal(t, "the-namespace", ProcessConfig.Namespace)
		// from file
		assert.Equal(t, "10.20.30", ProcessConfig.DriverVersion)
		assert.Equal(t, map[string]string{"node-type": "gpu"}, ProcessConfig.NodeSelector)
		// defaulted
		assert.Equal(t, "ghcr.io/gardenlinux/gardenlinux/kmodbuild", ProcessConfig.CompileImage)
	})
}
