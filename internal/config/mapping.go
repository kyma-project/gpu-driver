package config

import (
	"github.com/spf13/viper"
	"sync"
)

const (
	keyKmodbuildVersions    = "kmodbuildVersions"
	keyDriverVersions       = "driverVersions"
	keyDefaultDriverVersion = "defaultDriverVersion"
)

var cfg *viper.Viper

type mappingData struct {
	m sync.Mutex

	kmodbuildVersions    map[string]string
	defaultDriverVersion string
	driverVersions       map[string]string
}

var defaultKmodbuildVersions = map[string]string{
	"6.6.62-cloud-amd64": "amd64-1592.3",
	"6.6.63-cloud-amd64": "amd64-1592.4",
	"6.6.71-cloud-amd64": "amd64-1592.5",
	"6.6.78-cloud-amd64": "amd64-1592.6",
}

const defaultDriverVersion = "550.127.08"

var defaultDriverVersions = map[string]string{}

func Initialize(dir string) error {
	cfg = viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath(dir)
	cfg.AddConfigPath("/opt/gpu-driver/config")

	cfg.SetDefault(keyKmodbuildVersions, defaultKmodbuildVersions)
	cfg.SetDefault(keyDefaultDriverVersion, defaultDriverVersions)
	cfg.SetDefault(keyDefaultDriverVersion, defaultDriverVersion)

	err := cfg.ReadInConfig()
	if err != nil {
		cfg.WatchConfig()
	}

	return err
}

func DefaultDriverVersion() string {
	return cfg.GetString(keyDefaultDriverVersion)
}

func KernelToKmodbuild(kernel string) string {
	mapping := cfg.GetStringMapString(keyKmodbuildVersions)
	val, ok := mapping[kernel]
	if !ok {
		return ""
	}
	return val
}

func KernelToDriver(kernel string) string {
	mapping := cfg.GetStringMapString(keyDriverVersions)
	val, ok := mapping[kernel]
	if !ok || val == "" {
		return defaultDriverVersion
	}
	return val
}
