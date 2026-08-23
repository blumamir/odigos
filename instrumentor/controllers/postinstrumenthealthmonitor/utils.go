package postinstrumenthealthmonitor

import (
	"reflect"
	"time"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	containerutils "github.com/odigos-io/odigos/k8sutils/pkg/container"
	"github.com/odigos-io/odigos/status"
	postInstrumentHealthMonitor "github.com/odigos-io/odigos/status/instrumentationconfig/generated"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type postInstrumentMonitorPodState int

const (
	postInstrumentMonitorPodNotApplicable postInstrumentMonitorPodState = iota
	postInstrumentMonitorPodInGracePeriod
	postInstrumentMonitorPodApplicable
)

// postInstrumentMonitorPodStateForPod classifies a pod for the post-instrument health monitor.
// It returns the agent start time for instrumented pods (including those still in grace).
func postInstrumentMonitorPodStateForPod(pod *corev1.Pod, ic *odigosv1.InstrumentationConfig, graceTime time.Duration, now time.Time) (postInstrumentMonitorPodState, *metav1.Time) {
	if pod.Status.StartTime == nil {
		return postInstrumentMonitorPodNotApplicable, nil
	}

	podManifestInjectionRequired := !ic.Spec.PodManifestInjectionOptional
	if podManifestInjectionRequired {
		odigosAgentsMetaHash, hasOdigosMetaHashLabel := pod.Labels[k8sconsts.OdigosAgentsMetaHashLabel]
		if !hasOdigosMetaHashLabel {
			return postInstrumentMonitorPodNotApplicable, nil
		}
		if odigosAgentsMetaHash != ic.Spec.AgentsMetaHash {
			return postInstrumentMonitorPodNotApplicable, nil
		}
	}

	var agentStartTime *metav1.Time
	if ic.Spec.PodManifestInjectionOptional {
		agentStartTime = ic.Spec.AgentsMetaHashChangedTime
	} else {
		agentStartTime = pod.Status.StartTime
	}
	if agentStartTime == nil {
		return postInstrumentMonitorPodNotApplicable, nil
	}

	if agentStartTime.Add(graceTime).After(now) {
		return postInstrumentMonitorPodInGracePeriod, agentStartTime
	}

	return postInstrumentMonitorPodApplicable, agentStartTime
}

func hasAgentContainerInCrashLoopBackOff(pod *corev1.Pod, ic *odigosv1.InstrumentationConfig) bool {
	agentContainers := make(map[string]struct{})
	for _, container := range ic.Spec.Containers {
		if container.AgentEnabled {
			agentContainers[container.ContainerName] = struct{}{}
		}
	}
	if len(agentContainers) == 0 {
		return false
	}

	allStatuses := append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...)
	for i := range allStatuses {
		cs := &allStatuses[i]
		if _, ok := agentContainers[cs.Name]; !ok {
			continue
		}
		if containerutils.IsContainerInCrashLoopBackOff(cs) {
			return true
		}
	}
	return false
}

func hasOdigosInitContainerInImagePullBackOff(pod *corev1.Pod) bool {
	for i := range pod.Status.InitContainerStatuses {
		cs := &pod.Status.InitContainerStatuses[i]
		if cs.Name != k8sconsts.OdigosInitContainerName {
			continue
		}
		return containerutils.IsContainerInImagePullBackOff(cs)
	}
	return false
}

func hasEnabledAgent(ic *odigosv1.InstrumentationConfig) bool {
	for _, container := range ic.Spec.Containers {
		if container.AgentEnabled {
			return true
		}
	}
	return false
}

func isMonitoringSource(ic *odigosv1.InstrumentationConfig, rollbackConfig autoRollbackConfig) bool {
	return !rollbackConfig.disabled && hasEnabledAgent(ic)
}

func isResultSet(postInstrumentHealthMonitor *odigosv1.PostInstrumentHealthMonitor) bool {
	return postInstrumentHealthMonitor != nil && postInstrumentHealthMonitor.HealthCheckResult != nil
}

func isStabilityWindowEnded(postInstrumentHealthMonitor *odigosv1.PostInstrumentHealthMonitor, rollbackConfig *autoRollbackConfig, now time.Time) bool {
	if postInstrumentHealthMonitor == nil {
		return false
	}
	firstInstrumentedPodStartTime := postInstrumentHealthMonitor.FirstInstrumentedPodStartTime
	if firstInstrumentedPodStartTime == nil {
		return false
	}

	endTime := firstInstrumentedPodStartTime.Add(rollbackConfig.stabilityWindow)
	return now.After(endTime)
}

func shouldUpdateFirstInstrumentedPodStartTime(firstInstrumentedPodStartTime *metav1.Time, agentStartTime *metav1.Time) bool {
	if agentStartTime == nil {
		return false
	}
	if firstInstrumentedPodStartTime == nil {
		return true
	}
	return agentStartTime.Before(firstInstrumentedPodStartTime)
}

func postInstrumentHealthMonitorEqual(a, b *odigosv1.PostInstrumentHealthMonitor) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return reflect.DeepEqual(*a, *b)
}

func isActivelyMonitoring(ic *odigosv1.InstrumentationConfig, rollbackConfig autoRollbackConfig) bool {
	return isMonitoringSource(ic, rollbackConfig) && !isResultSet(ic.Spec.PostInstrumentHealthMonitor)
}

func reasonFromHealthMonitorResult(monitor *odigosv1.PostInstrumentHealthMonitor) status.Reason {
	if monitor.HealthCheckResult != nil && *monitor.HealthCheckResult {
		return postInstrumentHealthMonitor.PostInstrumentHealthMonitorStable
	}
	switch monitor.UnhealthyReason {
	case odigosv1.PostInstrumentHealthUnhealthyReasonOdigosInitContainerImagePullError:
		return postInstrumentHealthMonitor.PostInstrumentHealthMonitorImagePullBackOff
	case odigosv1.PostInstrumentHealthUnhealthyReasonUnhealthyAfterInstrumentation:
		return postInstrumentHealthMonitor.PostInstrumentHealthMonitorUnhealthy
	default:
		return postInstrumentHealthMonitor.PostInstrumentHealthMonitorEvaluating
	}
}
