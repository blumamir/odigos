package status

import (
	"testing"

	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/frontend/graph/model"
	postInstrumentHealthMonitorStatus "github.com/odigos-io/odigos/status/instrumentationconfig/generated"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func reasonOf(t *testing.T, s *model.DesiredConditionStatus) string {
	t.Helper()
	if s == nil || s.ReasonEnum == nil {
		return ""
	}
	return *s.ReasonEnum
}

func TestCalculateAutoRollbackStatus_Unhealthy(t *testing.T) {
	ic := &odigosv1alpha1.InstrumentationConfig{
		Status: odigosv1alpha1.InstrumentationConfigStatus{
			Conditions: []metav1.Condition{
				{
					Type:    postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorType,
					Status:  metav1.ConditionFalse,
					Reason:  string(postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorReasonUnhealthy),
					Message: "pod unhealthy after instrumentation",
				},
			},
		},
	}

	got := CalculateAutoRollbackStatus(ic)

	if got == nil {
		t.Fatal("expected a status, got nil")
	}
	if reasonOf(t, got) != postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorUnhealthy.Title {
		t.Fatalf("expected reason %q, got %q", postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorUnhealthy.Title, reasonOf(t, got))
	}
	if got.Status != model.DesiredStateProgressNotice {
		t.Fatalf("expected status %q, got %q", model.DesiredStateProgressNotice, got.Status)
	}
	if len(got.ActionItems) != 1 {
		t.Fatalf("expected recovery action item, got %d", len(got.ActionItems))
	}
}

func TestCalculateAutoRollbackStatus_Stable(t *testing.T) {
	ic := &odigosv1alpha1.InstrumentationConfig{
		Status: odigosv1alpha1.InstrumentationConfigStatus{
			Conditions: []metav1.Condition{
				{
					Type:    postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorType,
					Status:  metav1.ConditionTrue,
					Reason:  string(postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorReasonStable),
					Message: "pods are stable after instrumentation",
				},
			},
		},
	}

	got := CalculateAutoRollbackStatus(ic)

	if reasonOf(t, got) != postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorStable.Title {
		t.Fatalf("expected reason %q, got %q", postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorStable.Title, reasonOf(t, got))
	}
	if got.Status != model.DesiredStateProgressSuccess {
		t.Fatalf("expected status %q, got %q", model.DesiredStateProgressSuccess, got.Status)
	}
}

func TestCalculateAutoRollbackStatus_NoMonitorCondition(t *testing.T) {
	ic := &odigosv1alpha1.InstrumentationConfig{}

	got := CalculateAutoRollbackStatus(ic)

	if reasonOf(t, got) != postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorNotApplicable.Title {
		t.Fatalf("expected reason %q, got %q", postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorNotApplicable.Title, reasonOf(t, got))
	}
	if got.Status != model.DesiredStateProgressIrrelevant {
		t.Fatalf("expected status %q, got %q", model.DesiredStateProgressIrrelevant, got.Status)
	}
	if got.Message != postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorNotApplicable.Message {
		t.Fatalf("expected message %q, got %q", postInstrumentHealthMonitorStatus.PostInstrumentHealthMonitorNotApplicable.Message, got.Message)
	}
}
