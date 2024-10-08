package container

import (
	"strings"

	"github.com/odigos-io/odigos/common"
	v1 "k8s.io/api/core/v1"
)

func GetLanguageAndOtelSdk(container v1.Container) (common.ProgrammingLanguage, common.OtelSdk, bool) {
	deviceName := PodContainerDeviceName(container)
	if deviceName == nil {
		return common.UnknownProgrammingLanguage, common.OtelSdk{}, false
	}

	language, sdk := common.InstrumentationDeviceNameToComponents(*deviceName)
	return language, sdk, true
}

func PodContainerDeviceName(container v1.Container) *string {
	if container.Resources.Limits == nil {
		return nil
	}

	for resourceName := range container.Resources.Limits {
		resourceNameStr := string(resourceName)
		if strings.HasPrefix(resourceNameStr, common.OdigosResourceNamespace) {
			return &resourceNameStr
		}
	}

	return nil
}
