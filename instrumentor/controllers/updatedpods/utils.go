package updatedpods

import (
	"github.com/odigos-io/odigos/api/k8sconsts"
	corev1 "k8s.io/api/core/v1"
)

func getDistroNameFromContainerName(container *corev1.Container) (string, bool) {
	for _, env := range container.Env {
		if env.Name == k8sconsts.OdigosEnvVarDistroName {
			return env.Value, true
		}
	}
	return "", false
}
