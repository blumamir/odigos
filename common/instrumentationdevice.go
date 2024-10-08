package common

import (
	"strings"
)

type OdigosInstrumentationDevice string

// This is the resource namespace of the lister in k8s device plugin manager.
// from the "github.com/kubevirt/device-plugin-manager" package source:
// GetResourceNamespace must return namespace (vendor ID) of implemented Lister. e.g. for
// resources in format "color.example.com/<color>" that would be "color.example.com".
const OdigosResourceNamespace = "instrumentation.odigos.io"

// Special experimental devices to apply multiple languages instrumentations at once.
// all-native-community - all supported languages with native instrumentations and community otel SDK
const OdigosInstrumentationPluginAllNativeCommunity = "instrumentation.odigos.io/all-native-community"

// all-native-enterprise - all supported languages with native instrumentations and enterprise eBPF otel SDK
const OdigosInstrumentationPluginAllNativeEnterprise = "instrumentation.odigos.io/all-native-enterprise"

// all-ebpf-enterprise - all supported languages with eBPF instrumentations (where supported) and enterprise otel SDK
const OdigosInstrumentationPluginAllEbpfEnterprise = "instrumentation.odigos.io/all-ebpf-enterprise"

// The plugin name is also part of the device-plugin-manager.
// It is used to control which environment variables and fs mounts are available to the pod.
// Each OtelSdkType has its own plugin name, which is used to control instrumentation for that SDK.
//
// Odigos convention for plugin name is as follows:
// <runtime-name>-<sdk-type>-<sdk-tier>
//
// for example:
// the native Java SDK will be named "java-native-community".
// the ebpf Java enterprise sdk will be named "java-ebpf-enterprise".

func InstrumentationPluginName(language ProgrammingLanguage, otelSdk OtelSdk) string {
	if language == "" { // all languages
		if otelSdk == OtelSdkNativeEnterprise {
			return "all-native-enterprise"
		} else if otelSdk == OtelSdkEbpfEnterprise {
			return "all-ebpf-enterprise"
		} else {
			return "all-native-community"
		}
	} else {
		return string(language) + "-" + string(otelSdk.SdkType) + "-" + string(otelSdk.SdkTier)
	}
}

func InstrumentationPluginNameToComponents(pluginName string) (ProgrammingLanguage, OtelSdk) {
	components := strings.Split(pluginName, "-")
	otelSdk := OtelSdk{SdkType: OtelSdkType(components[1]), SdkTier: OtelSdkTier(components[2])}
	return ProgrammingLanguage(components[0]), otelSdk
}

func InstrumentationDeviceName(language ProgrammingLanguage, otelSdk OtelSdk) OdigosInstrumentationDevice {
	pluginName := InstrumentationPluginName(language, otelSdk)
	return OdigosInstrumentationDevice(OdigosResourceNamespace + "/" + pluginName)
}

func InstrumentationDeviceNameToComponents(deviceName string) (ProgrammingLanguage, OtelSdk) {
	pluginName := strings.Split(deviceName, "/")[1]
	return InstrumentationPluginNameToComponents(pluginName)
}

func IsResourceNameOdigosInstrumentation(resourceName string) bool {
	return strings.HasPrefix(resourceName, OdigosResourceNamespace)
}
