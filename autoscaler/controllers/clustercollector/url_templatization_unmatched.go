package clustercollector

import (
	"github.com/odigos-io/odigos/api/k8sconsts"
	commonconf "github.com/odigos-io/odigos/autoscaler/controllers/common"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/config"
	pipelinegen "github.com/odigos-io/odigos/common/pipelinegen"
)

// effectiveCardinalityControl returns the cardinality-control config the gateway
// should be wired for. Auto-compute URL templatization rules is enterprise-only —
// helm already gates the helm value on an on-prem token, and community tier must
// not install the side-channel pipeline. Same reasoning as effectiveInsightsConfig.
func effectiveCardinalityControl(cc *common.CardinalityControlConfiguration, tier common.OdigosTier) *common.CardinalityControlConfiguration {
	if !tier.IsEnterprise() {
		return nil
	}
	return cc
}

// addUrlTemplatizationUnmatchedExporter appends the enterprise
// odigos_url_templatization exporter to the gateway root traces pipeline so
// every span fans out after root processors (including odigosurltemplate when
// that processor is configured). The exporter counts unmatched HTTP paths in
// cacheDb Redis. Noop when disabled or when no root traces pipeline exists.
func addUrlTemplatizationUnmatchedExporter(c *config.Config, odigosNs string, cardinalityControl *common.CardinalityControlConfiguration) error {
	if !common.UrlTemplatizationAutoComputeRulesActive(cardinalityControl) {
		return nil
	}

	rootPipelineName := pipelinegen.GetTelemetryRootPipelineName(common.TracesObservabilitySignal)
	rootPipeline, hasRoot := c.Service.Pipelines[rootPipelineName]
	if !hasRoot {
		return nil
	}

	if c.Exporters == nil {
		c.Exporters = config.GenericMap{}
	}

	c.Exporters[commonconf.UrlTemplatizationExporter] = config.GenericMap{
		"odigos_config_extension": k8sconsts.OdigosConfigK8sExtensionType,
		"endpoint":                k8sconsts.OdigosCacheEndpoint(odigosNs),
	}

	rootPipeline.Exporters = append(rootPipeline.Exporters, commonconf.UrlTemplatizationExporter)
	c.Service.Pipelines[rootPipelineName] = rootPipeline

	return nil
}
