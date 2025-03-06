package config

var namespace = "gpu-driver-system"

var scriptsConfigMapName = "gpu-driver-scripts"

var installerServiceAccountName = "gpu-driver-installer"

func SetNamespace(ns string) {
	namespace = ns
}

func GetNamespace() string {
	return namespace
}

func SetScriptsConfigMapName(name string) {
	scriptsConfigMapName = name
}

func GetScriptsConfigMapName() string {
	return scriptsConfigMapName
}

func SetInstallerServiceAccountName(name string) {
	installerServiceAccountName = name
}

func GetInstallerServiceAccountName() string {
	return installerServiceAccountName
}
