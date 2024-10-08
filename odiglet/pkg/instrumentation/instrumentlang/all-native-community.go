package instrumentlang

import (
	"fmt"

	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/envOverwrite"
	"github.com/odigos-io/odigos/odiglet/pkg/env"
	"github.com/odigos-io/odigos/odiglet/pkg/instrumentation/consts"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	envOtlpEndpoint             = "OTEL_EXPORTER_OTLP_ENDPOINT"
	envOtelExporterOtlpProtocol = "OTEL_EXPORTER_OTLP_PROTOCOL"
	envOtelServiceName          = "OTEL_SERVICE_NAME"
	envOpampServerHost          = "ODIGOS_OPAMP_SERVER_HOST"
	envInstrumentationDeviceId  = "ODIGOS_INSTRUMENTATION_DEVICE_ID"
	envOtelTracesSampler        = "OTEL_TRACES_SAMPLER"
	envOtelResourceAttributes   = "OTEL_RESOURCE_ATTRIBUTES"
)

func AllNativeCommunity(deviceId string, uniqueDestinationSignals map[common.ObservabilitySignal]struct{}) *v1beta1.ContainerAllocateResponse {

	// general environment variables for all languages
	otlpEndpoint := fmt.Sprintf("http://%s:%d", env.Current.NodeIP, consts.OTLPHttpPort)
	opampServerHost := fmt.Sprintf("%s:%d", env.Current.NodeIP, consts.OpAMPPort)

	metricsExporter := "none"
	tracesExporter := "none"
	if _, ok := uniqueDestinationSignals[common.MetricsObservabilitySignal]; ok {
		metricsExporter = "otlp"
	}
	if _, ok := uniqueDestinationSignals[common.TracesObservabilitySignal]; ok {
		tracesExporter = "otlp"
	}

	// nodejs
	nodeOptionsVal, _ := envOverwrite.ValToAppend(nodeEnvNodeOptions, common.OtelSdkNativeCommunity)

	// python
	pythonpathVal, _ := envOverwrite.ValToAppend(envPythonPath, common.OtelSdkNativeCommunity)

	// java
	javaOptsVal, _ := envOverwrite.ValToAppend(javaOptsEnvVar, common.OtelSdkNativeCommunity)
	javaToolOptionsVal, _ := envOverwrite.ValToAppend(javaToolOptionsEnvVar, common.OtelSdkNativeCommunity)

	return &v1beta1.ContainerAllocateResponse{
		Envs: map[string]string{

			// general environment variables for all languages
			envOtlpEndpoint:             otlpEndpoint,
			envOtelServiceName:          deviceId,
			envOpampServerHost:          opampServerHost,
			envInstrumentationDeviceId:  deviceId,
			envOtelTracesSampler:        "always_on",
			envOtelLogsExporter:         "none", // Logs are collected by the node collector with the file receiver
			envOtelTracesExporter:       tracesExporter,
			envOtelMetricsExporter:      metricsExporter,
			envOtelExporterOtlpProtocol: "http/protobuf",

			// go - taken care of by eBPF, no need for envs

			// nodejs
			nodeEnvNodeOptions: nodeOptionsVal,

			// python
			envPythonPath:     pythonpathVal,
			envLogCorrelation: "true",

			// java
			javaToolOptionsEnvVar: javaToolOptionsVal,
			javaOptsEnvVar:        javaOptsVal,

			// dotnet
			enableProfilingEnvVar: "1",
			profilerEndVar:        profilerId,
			profilerPathEnv:       fmt.Sprintf(profilerPath, getArch()), // TODO(edenfed): Support both musl and glibc. Requires improved language detection
			tracerHomeEnv:         tracerHome,
			startupHookEnv:        startupHook,
			additonalDepsEnv:      additonalDeps,
			sharedStoreEnv:        sharedStore,
		},
		Mounts: []*v1beta1.Mount{
			{
				ContainerPath: "/var/odigos",
				HostPath:      "/var/odigos",
				ReadOnly:      true,
			},
		},
	}
}
