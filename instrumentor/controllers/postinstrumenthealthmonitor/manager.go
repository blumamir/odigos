package postinstrumenthealthmonitor

import (
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	odigospredicate "github.com/odigos-io/odigos/k8sutils/pkg/predicate"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PodsPredicate limits pod events to creates, deletes, and container waiting-reason updates.
// Create covers agent labels and leaving the "no pods" state.
// Delete covers entering the "no pods" state when the last pod is removed.
// Update fires only when State.Waiting.Reason changes (the field checked for CrashLoopBackOff / ImagePullBackOff).
type PodsPredicate struct{}

func (p PodsPredicate) Create(e event.CreateEvent) bool {
	return true
}

func (p PodsPredicate) Update(e event.UpdateEvent) bool {
	oldPod, oldOk := e.ObjectOld.(*corev1.Pod)
	newPod, newOk := e.ObjectNew.(*corev1.Pod)
	if !oldOk || !newOk {
		return false
	}
	return containerWaitingReasonsChanged(oldPod.Status.ContainerStatuses, newPod.Status.ContainerStatuses) ||
		containerWaitingReasonsChanged(oldPod.Status.InitContainerStatuses, newPod.Status.InitContainerStatuses)
}

func (p PodsPredicate) Delete(e event.DeleteEvent) bool {
	return true
}

func (p PodsPredicate) Generic(e event.GenericEvent) bool {
	return false
}

func containerWaitingReasonsChanged(oldStatuses, newStatuses []corev1.ContainerStatus) bool {
	if len(oldStatuses) != len(newStatuses) {
		return true
	}
	for i := range oldStatuses {
		if containerWaitingReason(&oldStatuses[i]) != containerWaitingReason(&newStatuses[i]) {
			return true
		}
	}
	return false
}

func containerWaitingReason(cs *corev1.ContainerStatus) string {
	if cs.State.Waiting == nil {
		return ""
	}
	return cs.State.Waiting.Reason
}

// InstrumentationConfigPredicate limits InstrumentationConfig events to creates and AgentsMetaHash changes.
// Agent enablement, container identity, and rollout-relevant config all flow through AgentsMetaHash.
type InstrumentationConfigPredicate struct{}

func (p InstrumentationConfigPredicate) Create(e event.CreateEvent) bool {
	return true
}

func (p InstrumentationConfigPredicate) Update(e event.UpdateEvent) bool {
	oldIC, oldOk := e.ObjectOld.(*odigosv1.InstrumentationConfig)
	newIC, newOk := e.ObjectNew.(*odigosv1.InstrumentationConfig)
	if !oldOk || !newOk {
		return false
	}
	return oldIC.Spec.AgentsMetaHash != newIC.Spec.AgentsMetaHash
}

func (p InstrumentationConfigPredicate) Delete(e event.DeleteEvent) bool {
	return false
}

func (p InstrumentationConfigPredicate) Generic(e event.GenericEvent) bool {
	return false
}

func SetupWithManager(mgr ctrl.Manager) error {
	monitorEnabled := &monitorEnabledState{}
	// Default on until the first effective-config reconcile, matching auto-rollback defaults.
	monitorEnabled.enabled.Store(true)

	err := builder.
		ControllerManagedBy(mgr).
		Named("postinstrumenthealthmonitor-pods").
		For(&corev1.Pod{}).
		WithEventFilter(predicate.And(
			MonitorEnabledPredicate{monitorEnabled: monitorEnabled},
			&PodsPredicate{},
		)).
		Complete(&PodsController{
			Client:         mgr.GetClient(),
			monitorEnabled: monitorEnabled,
		})
	if err != nil {
		return err
	}

	err = builder.
		ControllerManagedBy(mgr).
		Named("postinstrumenthealthmonitor-instrumentationconfig").
		For(&odigosv1.InstrumentationConfig{}).
		WithEventFilter(predicate.And(
			MonitorEnabledPredicate{monitorEnabled: monitorEnabled},
			&InstrumentationConfigPredicate{},
		)).
		Complete(&InstrumentationConfigController{
			Client:         mgr.GetClient(),
			monitorEnabled: monitorEnabled,
		})
	if err != nil {
		return err
	}

	err = builder.
		ControllerManagedBy(mgr).
		Named("postinstrumenthealthmonitor-effectiveconfig").
		For(&corev1.ConfigMap{}).
		WithEventFilter(odigospredicate.OdigosEffectiveConfigMapPredicate).
		Complete(&EffectiveConfigReconciler{
			Client:         mgr.GetClient(),
			monitorEnabled: monitorEnabled,
		})
	if err != nil {
		return err
	}

	return nil
}
