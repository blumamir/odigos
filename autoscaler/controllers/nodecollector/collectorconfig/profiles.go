package collectorconfig

import (
	"fmt"

	"github.com/odigos-io/odigos/api/k8sconsts"
	commonconf "github.com/odigos-io/odigos/autoscaler/controllers/common"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/config"
	odigosconsts "github.com/odigos-io/odigos/common/consts"
)

// ProfilingPipelineConfig builds the node collector profiles domain when profiling is enabled.
// When profilesLoadBalancingNeeded is true (interrogation), profiles are exported via
// odigos_profiles_loadbalancing so TraceID-linked samples land on the same gateway as spans.
func ProfilingPipelineConfig(odigosNamespace string, profiling *common.ProfilingConfiguration, manifestProcessorNames []string, profilesLoadBalancingNeeded bool) config.Config {
	if !common.ProfilingPipelineActive(profiling) {
		return config.Config{}
	}

	exporters, exporterName := profilesExporters(odigosNamespace, profiling, profilesLoadBalancingNeeded)

	// memory_limiter itself is defined once, globally, by commonProcessors() (see
	// common.go) — every pipeline just references its name, never redefines it.
	processors := config.GenericMap{
		commonconf.ProfilingNodeFilterProcessor:         commonconf.ProfilingFilterProcessorConfig(),
		commonconf.ProfilingNodeK8sAttributesProcessor:  commonconf.K8sAttributesProfilesProcessorConfig(),
		commonconf.ProfilingNodeOdigosProfilesProcessor: commonconf.OdigosProfilesProcessorConfig(),
		commonconf.ProfilingNodeServiceNameProcessor:    commonconf.ProfilingServiceNameTransformConfig(),
	}
	pipelineProcessors := []string{
		memoryLimiterProcessorName,
		commonconf.ProfilingNodeFilterProcessor,
		commonconf.ProfilingNodeK8sAttributesProcessor,
		commonconf.ProfilingNodeOdigosProfilesProcessor,
	}
	// Native symbolization is opt-in (profiling.symbolization.native). When on, the
	// symbolize processor runs after the keep-filter (only retained profiles are
	// symbolized) and before service-name enrichment.
	if profiling.NativeSymbolizationEnabled() {
		processors[commonconf.ProfilingNodeSymbolizeProcessor] = commonconf.OdigosSymbolizeProcessorConfig()
		pipelineProcessors = append(pipelineProcessors, commonconf.ProfilingNodeSymbolizeProcessor)
	}
	pipelineProcessors = append(pipelineProcessors, commonconf.ProfilingNodeServiceNameProcessor)
	pipelineProcessors = append(pipelineProcessors, manifestProcessorNames...)
	pipelineProcessors = append(pipelineProcessors, odigosTrafficMetricsProcessorName) // keep traffic metrics last for most accurate tracking

	return config.Config{
		Receivers: config.GenericMap{
			commonconf.ProfilingReceiver: config.GenericMap{
				"obi_process_ctx": true,
				"probe_links": []string{
					"kprobe:uprobe_notify_resume",
				},
			},
		},
		Processors: processors,
		Exporters:  exporters,
		Service: config.Service{
			Pipelines: map[string]config.Pipeline{
				"profiles": {
					Receivers:  []string{commonconf.ProfilingReceiver},
					Processors: pipelineProcessors,
					Exporters:  []string{exporterName},
				},
			},
		},
	}
}

func profilesExporters(odigosNamespace string, profiling *common.ProfilingConfiguration, profilesLoadBalancingNeeded bool) (config.GenericMap, string) {
	if profilesLoadBalancingNeeded {
		service := fmt.Sprintf("%s.%s", k8sconsts.OdigosClusterCollectorServiceName, odigosNamespace)
		otlpConfig := commonconf.MergeProfilingOtlpExporter(config.GenericMap{
			"tls":         config.GenericMap{"insecure": true},
			"compression": "none",
		}, profiling.Exporter)
		return config.GenericMap{
			commonconf.ProfilingNodeLoadbalancingExporter: config.GenericMap{
				"protocol": config.GenericMap{
					"otlp": otlpConfig,
				},
				"resolver": config.GenericMap{
					"k8s": config.GenericMap{
						"service": service,
					},
				},
			},
		}, commonconf.ProfilingNodeLoadbalancingExporter
	}

	endpoint := k8sconsts.OtlpGrpcDNSEndpoint(k8sconsts.OdigosClusterCollectorServiceName, odigosNamespace, odigosconsts.OTLPPort)
	exp := commonconf.MergeProfilingOtlpExporter(config.GenericMap{
		"endpoint":    endpoint,
		"tls":         config.GenericMap{"insecure": true},
		"compression": "none",
	}, profiling.Exporter)
	return config.GenericMap{
		commonconf.ProfilingNodeToGatewayExporter: exp,
	}, commonconf.ProfilingNodeToGatewayExporter
}
