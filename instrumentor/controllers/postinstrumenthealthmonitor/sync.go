package postinstrumenthealthmonitor

import (
	"context"
	"time"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/k8sutils/pkg/utils"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	"github.com/odigos-io/odigos/status"
	postInstrumentHealthMonitor "github.com/odigos-io/odigos/status/instrumentationconfig/generated"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func postInstrumentHealthStateFromPods(pods []corev1.Pod, ic *odigosv1.InstrumentationConfig, rollbackConfig *autoRollbackConfig, now time.Time) (odigosv1.PostInstrumentHealthMonitor, status.Reason) {
	var firstInstrumentedPodStartTime *metav1.Time
	var unhealthyAt *metav1.Time
	if ic.Spec.PostInstrumentHealthMonitor != nil {
		firstInstrumentedPodStartTime = ic.Spec.PostInstrumentHealthMonitor.FirstInstrumentedPodStartTime
		unhealthyAt = ic.Spec.PostInstrumentHealthMonitor.UnhealthyAt
	}

	imagePullBackOff := false
	unhealthy := false
	anyApplicablePod := false
	anyInGracePeriodPod := false

	for i := range pods {
		if hasOdigosInitContainerInImagePullBackOff(&pods[i]) {
			imagePullBackOff = true
		}
		podState, agentStartTime := postInstrumentMonitorPodStateForPod(&pods[i], ic, rollbackConfig.graceTime, now)
		if shouldUpdateFirstInstrumentedPodStartTime(firstInstrumentedPodStartTime, agentStartTime) {
			firstInstrumentedPodStartTime = agentStartTime
		}
		switch podState {
		case postInstrumentMonitorPodApplicable:
			anyApplicablePod = true
			if hasAgentContainerInCrashLoopBackOff(&pods[i], ic) {
				unhealthy = true
			}
		case postInstrumentMonitorPodInGracePeriod:
			anyInGracePeriodPod = true
		}
	}

	if unhealthy {
		healthy := false
		return odigosv1.PostInstrumentHealthMonitor{
			HealthCheckResult:             &healthy,
			UnhealthyReason:               odigosv1.PostInstrumentHealthUnhealthyReasonUnhealthyAfterInstrumentation,
			UnhealthyAt:                   unhealthyAtTime(unhealthyAt, now),
			FirstInstrumentedPodStartTime: firstInstrumentedPodStartTime,
		}, postInstrumentHealthMonitor.PostInstrumentHealthMonitorUnhealthy
	}

	if imagePullBackOff {
		healthy := false
		return odigosv1.PostInstrumentHealthMonitor{
			HealthCheckResult:             &healthy,
			UnhealthyReason:               odigosv1.PostInstrumentHealthUnhealthyReasonOdigosInitContainerImagePullError,
			UnhealthyAt:                   unhealthyAtTime(unhealthyAt, now),
			FirstInstrumentedPodStartTime: firstInstrumentedPodStartTime,
		}, postInstrumentHealthMonitor.PostInstrumentHealthMonitorImagePullBackOff
	}

	if firstInstrumentedPodStartTime == nil {
		return odigosv1.PostInstrumentHealthMonitor{}, postInstrumentHealthMonitor.PostInstrumentHealthMonitorWaitingForInstrumentedPods
	}

	var reason status.Reason
	if anyApplicablePod {
		reason = postInstrumentHealthMonitor.PostInstrumentHealthMonitorEvaluating
	} else if anyInGracePeriodPod {
		reason = postInstrumentHealthMonitor.PostInstrumentHealthMonitorInGracePeriod
	} else {
		reason = postInstrumentHealthMonitor.PostInstrumentHealthMonitorNoInstrumentedPods
	}

	return odigosv1.PostInstrumentHealthMonitor{
		FirstInstrumentedPodStartTime: firstInstrumentedPodStartTime,
	}, reason
}

func desiredPostInstrumentHealthMonitorState(
	ctx context.Context,
	c client.Client,
	ic *odigosv1.InstrumentationConfig,
	pw k8sconsts.PodWorkload,
	rollbackConfig *autoRollbackConfig,
) (odigosv1.PostInstrumentHealthMonitor, status.Reason, error) {
	now := time.Now()

	pods, err := workload.ListNonCompletedPods(ctx, c, pw)
	if err != nil {
		return odigosv1.PostInstrumentHealthMonitor{}, status.Reason{}, err
	}
	if len(pods) == 0 {
		if ic.Spec.PostInstrumentHealthMonitor == nil {
			return odigosv1.PostInstrumentHealthMonitor{}, postInstrumentHealthMonitor.PostInstrumentHealthMonitorNoPods, nil
		}
		existing := *ic.Spec.PostInstrumentHealthMonitor
		if existing.FirstInstrumentedPodStartTime == nil {
			return existing, postInstrumentHealthMonitor.PostInstrumentHealthMonitorNoPods, nil
		}
		if existing.FirstInstrumentedPodStartTime.Add(rollbackConfig.graceTime).After(now) {
			return existing, postInstrumentHealthMonitor.PostInstrumentHealthMonitorInGracePeriod, nil
		}
		return existing, postInstrumentHealthMonitor.PostInstrumentHealthMonitorNoInstrumentedPods, nil
	}

	desired, reason := postInstrumentHealthStateFromPods(pods, ic, rollbackConfig, now)
	return desired, reason, nil
}

func calculatePostInstrumentHealthMonitor(
	ctx context.Context,
	c client.Client,
	ic *odigosv1.InstrumentationConfig,
	pw k8sconsts.PodWorkload,
	rollbackConfig *autoRollbackConfig,
) (*odigosv1.PostInstrumentHealthMonitor, *status.Reason, error) {
	if !isMonitoringSource(ic, *rollbackConfig) {
		// if the result is set to unhealthy, we keep it (to denote it)
		// otherwise we clear it (to denote it's not monitoring)
		if isUnhealthyHealthCheckResult(ic.Spec.PostInstrumentHealthMonitor) {
			reason := reasonFromHealthMonitorResult(ic.Spec.PostInstrumentHealthMonitor)
			return ic.Spec.PostInstrumentHealthMonitor, reason, nil
		}
		return nil, nil, nil
	}

	if isResultSet(ic.Spec.PostInstrumentHealthMonitor) {
		reason := reasonFromHealthMonitorResult(ic.Spec.PostInstrumentHealthMonitor)
		return ic.Spec.PostInstrumentHealthMonitor, reason, nil
	}

	now := time.Now()
	if isStabilityWindowEnded(ic.Spec.PostInstrumentHealthMonitor, rollbackConfig, now) {
		desired := ic.Spec.PostInstrumentHealthMonitor.DeepCopy()
		if desired == nil {
			desired = &odigosv1.PostInstrumentHealthMonitor{}
		}
		healthy := true
		desired.HealthCheckResult = &healthy
		desired.UnhealthyReason = ""
		desired.UnhealthyAt = nil
		reason := postInstrumentHealthMonitor.PostInstrumentHealthMonitorStable
		return desired, &reason, nil
	}

	desired, reason, err := desiredPostInstrumentHealthMonitorState(ctx, c, ic, pw, rollbackConfig)
	if err != nil {
		return nil, nil, err
	}
	return &desired, &reason, nil
}

func syncWorkload(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) (bool, error) {
	icName := workload.CalculateWorkloadRuntimeObjectName(pw.Name, pw.Kind)
	ic := odigosv1.InstrumentationConfig{}
	err := c.Get(ctx, types.NamespacedName{Namespace: pw.Namespace, Name: icName}, &ic)
	if err != nil {
		return false, client.IgnoreNotFound(err)
	}

	effectiveConfig, err := utils.GetCurrentOdigosConfiguration(ctx, c)
	if err != nil {
		return false, err
	}
	rollbackConfig, err := getAutoRollbackConfig(&effectiveConfig)
	if err != nil {
		return false, err
	}

	desiredMonitor, reason, err := calculatePostInstrumentHealthMonitor(ctx, c, &ic, pw, &rollbackConfig)
	if err != nil {
		return false, err
	}

	if !postInstrumentHealthMonitorEqual(ic.Spec.PostInstrumentHealthMonitor, desiredMonitor) {
		ic.Spec.PostInstrumentHealthMonitor = desiredMonitor
		if err := c.Update(ctx, &ic); err != nil {
			return false, err
		}
	}

	conditionChanged := false
	if reason == nil {
		conditionChanged = meta.RemoveStatusCondition(&ic.Status.Conditions, postInstrumentHealthMonitor.PostInstrumentHealthMonitorType)
	} else {
		message, _ := status.RenderMessage(*reason, postInstrumentHealthMonitor.PostInstrumentHealthMonitorMessageParams{
			GracePeriod:     rollbackConfig.graceTime.String(),
			StabilityWindow: rollbackConfig.stabilityWindow.String(),
		})
		conditionChanged = meta.SetStatusCondition(&ic.Status.Conditions, metav1.Condition{
			Type:    postInstrumentHealthMonitor.PostInstrumentHealthMonitorType,
			Status:  reason.K8sConditionStatus,
			Reason:  reason.Name,
			Message: message,
		})
	}

	if conditionChanged {
		if err := c.Status().Update(ctx, &ic); err != nil {
			return false, err
		}
	}

	return isActivelyMonitoring(&ic, rollbackConfig), nil
}
