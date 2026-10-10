package clustercollector

import (
	"testing"

	commonconf "github.com/odigos-io/odigos/autoscaler/controllers/common"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/config"
	pipelinegen "github.com/odigos-io/odigos/common/pipelinegen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func liveTrafficLearning(enabled bool) *common.UrlTemplatizationCardinalityControlConfiguration {
	return &common.UrlTemplatizationCardinalityControlConfiguration{
		LiveTrafficLearning: &common.LiveTrafficLearningConfiguration{Enabled: &enabled},
	}
}

func TestAddUrlTemplatizationUnmatchedExporter_Disabled(t *testing.T) {
	t.Run("nil_config_noop", func(t *testing.T) {
		c := configWithTracesIn()
		require.NoError(t, addUrlTemplatizationUnmatchedExporter(c, "odigos-system", nil))

		_, hasExp := c.Exporters[commonconf.UrlTemplatizationExporter]
		assert.False(t, hasExp)

		rootPipe := c.Service.Pipelines[pipelinegen.GetTelemetryRootPipelineName(common.TracesObservabilitySignal)]
		assert.Equal(t, []string{"odigosrouterconnector/traces"}, rootPipe.Exporters)
	})

	t.Run("explicit_false_noop", func(t *testing.T) {
		c := configWithTracesIn()
		require.NoError(t, addUrlTemplatizationUnmatchedExporter(c, "odigos-system", &common.CardinalityControlConfiguration{
			UrlTemplatization: liveTrafficLearning(false),
		}))

		_, hasExp := c.Exporters[commonconf.UrlTemplatizationExporter]
		assert.False(t, hasExp)
	})
}

func TestAddUrlTemplatizationUnmatchedExporter_NoTracesInPipelineNoop(t *testing.T) {
	c := &config.Config{Service: config.Service{Pipelines: map[string]config.Pipeline{}}}
	require.NoError(t, addUrlTemplatizationUnmatchedExporter(c, "odigos-system", &common.CardinalityControlConfiguration{
		UrlTemplatization: liveTrafficLearning(true),
	}))

	_, hasExp := c.Exporters[commonconf.UrlTemplatizationExporter]
	assert.False(t, hasExp)
}

func TestAddUrlTemplatizationUnmatchedExporter_Enabled(t *testing.T) {
	c := configWithTracesIn()
	require.NoError(t, addUrlTemplatizationUnmatchedExporter(c, "odigos-system", &common.CardinalityControlConfiguration{
		UrlTemplatization: liveTrafficLearning(true),
	}))

	exp, ok := c.Exporters[commonconf.UrlTemplatizationExporter].(config.GenericMap)
	require.True(t, ok)
	assert.Equal(t, "odigos_config_k8s", exp["odigos_config_extension"])
	assert.Equal(t, "odigos-cache.odigos-system:6379", exp["endpoint"])

	rootPipe := c.Service.Pipelines[pipelinegen.GetTelemetryRootPipelineName(common.TracesObservabilitySignal)]
	assert.Equal(t, []string{"resource/odigos-version", "transform/url-template"}, rootPipe.Processors,
		"root pipeline processors must be preserved so the exporter observes post-odigosurltemplate spans")
	assert.Equal(t, []string{"odigosrouterconnector/traces", commonconf.UrlTemplatizationExporter}, rootPipe.Exporters)

	assert.NotContains(t, c.Service.Pipelines, "traces/url-templatization-unmatched")
	assert.NotContains(t, c.Connectors, "forward/url-templatization-unmatched")
	assert.NotContains(t, c.Processors, "filter/url-templatization-unmatched")
}

func TestEffectiveCardinalityControl_CommunityTierNil(t *testing.T) {
	cc := &common.CardinalityControlConfiguration{
		UrlTemplatization: liveTrafficLearning(true),
	}
	assert.Nil(t, effectiveCardinalityControl(cc, common.CommunityOdigosTier))
	assert.Equal(t, cc, effectiveCardinalityControl(cc, common.OnPremOdigosTier))
}
