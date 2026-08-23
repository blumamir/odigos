package postinstrumenthealthmonitor

import (
	"sync/atomic"

	"sigs.k8s.io/controller-runtime/pkg/event"
	cr_predicate "sigs.k8s.io/controller-runtime/pkg/predicate"
)

// monitorEnabledState is shared across reconcilers. EffectiveConfigReconciler
// updates enabled from the odigos effective config; PodsController and
// InstrumentationConfigController check it to skip work when monitoring is off.
type monitorEnabledState struct {
	enabled atomic.Bool
}

// MonitorEnabledPredicate drops events when post-instrument health monitoring is disabled.
type MonitorEnabledPredicate struct {
	monitorEnabled *monitorEnabledState
}

func (p MonitorEnabledPredicate) Create(e event.CreateEvent) bool {
	return p.monitorEnabled.enabled.Load()
}

func (p MonitorEnabledPredicate) Update(e event.UpdateEvent) bool {
	return p.monitorEnabled.enabled.Load()
}

func (p MonitorEnabledPredicate) Delete(e event.DeleteEvent) bool {
	return p.monitorEnabled.enabled.Load()
}

func (p MonitorEnabledPredicate) Generic(e event.GenericEvent) bool {
	return false
}

var _ cr_predicate.Predicate = MonitorEnabledPredicate{}
