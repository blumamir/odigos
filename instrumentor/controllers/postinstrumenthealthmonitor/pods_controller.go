package postinstrumenthealthmonitor

import (
	"context"
	"time"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/k8sutils/pkg/utils"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodsController struct {
	client.Client
	monitorEnabled *monitorEnabledState
}

func (r *PodsController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if !r.monitorEnabled.enabled.Load() {
		return ctrl.Result{}, nil
	}

	var pod corev1.Pod
	err := r.Client.Get(ctx, req.NamespacedName, &pod)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	pw, err := workload.PodWorkloadObject(&pod)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Pod does not have a workload owner (e.g., system pods), skip processing.
	if pw == nil {
		return ctrl.Result{}, nil
	}

	activelyMonitoring, err := syncWorkload(ctx, r.Client, *pw)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !activelyMonitoring {
		return ctrl.Result{}, nil
	}

	requeueAfter, err := r.requeueAfterGracePeriod(ctx, &pod, *pw)
	if err != nil {
		return ctrl.Result{}, err
	}
	if requeueAfter > 0 {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

// requeueAfterGracePeriod schedules a follow-up reconcile when the triggering pod is still
// within the post-instrumentation grace period.
//
// This lives on the pods controller because grace period is scoped to a single pod: agent
// start time and grace expiry depend on that pod (and PodManifestInjectionOptional on the IC),
// not on the workload-wide state that syncWorkload evaluates.
//
// syncWorkload reports whether post-instrument health monitoring is still active for the
// source. IC and rollback config are re-fetched here only when monitoring is active, so we
// can compute a pod-specific requeue delay without pushing that cost into syncWorkload.
func (r *PodsController) requeueAfterGracePeriod(ctx context.Context, pod *corev1.Pod, pw k8sconsts.PodWorkload) (time.Duration, error) {
	icName := workload.CalculateWorkloadRuntimeObjectName(pw.Name, pw.Kind)
	ic := odigosv1.InstrumentationConfig{}
	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: pw.Namespace, Name: icName}, &ic); err != nil {
		return 0, client.IgnoreNotFound(err)
	}

	effectiveConfig, err := utils.GetCurrentOdigosConfiguration(ctx, r.Client)
	if err != nil {
		return 0, err
	}
	rollbackConfig, err := getAutoRollbackConfig(&effectiveConfig)
	if err != nil {
		return 0, err
	}

	return podGracePeriodRequeueAfter(pod, &ic, rollbackConfig), nil
}

func podGracePeriodRequeueAfter(
	pod *corev1.Pod,
	ic *odigosv1.InstrumentationConfig,
	rollbackConfig autoRollbackConfig,
) time.Duration {
	now := time.Now()
	podState, agentStartTime := postInstrumentMonitorPodStateForPod(pod, ic, rollbackConfig.graceTime, now)
	if podState != postInstrumentMonitorPodInGracePeriod || agentStartTime == nil {
		return 0
	}

	remaining := agentStartTime.Add(rollbackConfig.graceTime).Sub(now)
	if remaining <= 0 {
		return 0
	}
	return remaining
}
