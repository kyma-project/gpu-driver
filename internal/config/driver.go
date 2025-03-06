package config

import (
	"github.com/elliotchance/pie/v2"
	"github.com/kyma-project/gpu-driver/api/v1beta1"
	"github.com/kyma-project/gpu-driver/internal/util"
	v1 "k8s.io/api/core/v1"
	"sync"
)

type driverData struct {
	m sync.Mutex

	data []*v1beta1.GpuDriver
}

var instance *driverData

func init() {
	instance = &driverData{}
}

func Remove(name string) {
	instance.m.Lock()
	defer instance.m.Unlock()

	instance.data = pie.FilterNot(instance.data, func(x *v1beta1.GpuDriver) bool {
		return x.Name == name
	})
}

func Sync(in *v1beta1.GpuDriver) {
	instance.m.Lock()
	defer instance.m.Unlock()

	if in.DeletionTimestamp.IsZero() {
		instance.data = append(instance.data, in)
	} else {
		instance.data = pie.FilterNot(instance.data, func(x *v1beta1.GpuDriver) bool {
			return x.Name == in.Name
		})
	}
}

func FindNodeConfigs(node *v1.Node) []*v1beta1.GpuDriver {
	return pie.Filter(instance.data, func(x *v1beta1.GpuDriver) bool {
		return util.MatchesLabels(node, x.Spec.NodeSelector)
	})
}

func DriverVersionForConfig(in *v1beta1.GpuDriver) string {
	if in.Spec.DriverVersion != "" {
		return in.Spec.DriverVersion
	}
	return DefaultDriverVersion()
}
