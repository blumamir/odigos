package agentenabled

import (
	"testing"

	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	agentInjectionEnabled "github.com/odigos-io/odigos/status/instrumentationconfig/generated"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPostInstrumentHealthMonitorStepDownCondition_UnhealthyAfterInstrumentation(t *testing.T) {
	healthy := false
	ic := &odigosv1alpha1.InstrumentationConfig{
		Spec: odigosv1alpha1.InstrumentationConfigSpec{
			PostInstrumentHealthMonitor: &odigosv1alpha1.PostInstrumentHealthMonitor{
				HealthCheckResult: &healthy,
				UnhealthyReason:   odigosv1alpha1.PostInstrumentHealthUnhealthyReasonUnhealthyAfterInstrumentation,
			},
		},
	}

	condition, steppedDown := postInstrumentHealthMonitorStepDownCondition(ic)

	assert.True(t, steppedDown)
	assert.NotNil(t, condition)
	assert.Equal(t, metav1.ConditionFalse, condition.Status)
	assert.Equal(t, odigosv1alpha1.AgentEnabledReason(agentInjectionEnabled.AgentEnabledReasonUnhealthyAfterInstrumentation), condition.Reason)
	assert.Equal(t, agentInjectionEnabled.AgentEnabledUnhealthyAfterInstrumentation.Message, condition.Message)
}

func TestPostInstrumentHealthMonitorStepDownCondition_InitContainerImagePullError(t *testing.T) {
	healthy := false
	ic := &odigosv1alpha1.InstrumentationConfig{
		Spec: odigosv1alpha1.InstrumentationConfigSpec{
			PostInstrumentHealthMonitor: &odigosv1alpha1.PostInstrumentHealthMonitor{
				HealthCheckResult: &healthy,
				UnhealthyReason:   odigosv1alpha1.PostInstrumentHealthUnhealthyReasonOdigosInitContainerImagePullError,
			},
		},
	}

	condition, steppedDown := postInstrumentHealthMonitorStepDownCondition(ic)

	assert.True(t, steppedDown)
	assert.NotNil(t, condition)
	assert.Equal(t, odigosv1alpha1.AgentEnabledReason(agentInjectionEnabled.AgentEnabledReasonInitContainerImagePullError), condition.Reason)
	assert.Equal(t, agentInjectionEnabled.AgentEnabledInitContainerImagePullError.Message, condition.Message)
}

func TestPostInstrumentHealthMonitorStepDownCondition_Healthy(t *testing.T) {
	healthy := true
	ic := &odigosv1alpha1.InstrumentationConfig{
		Spec: odigosv1alpha1.InstrumentationConfigSpec{
			PostInstrumentHealthMonitor: &odigosv1alpha1.PostInstrumentHealthMonitor{
				HealthCheckResult: &healthy,
			},
		},
	}

	condition, steppedDown := postInstrumentHealthMonitorStepDownCondition(ic)

	assert.False(t, steppedDown)
	assert.Nil(t, condition)
}
