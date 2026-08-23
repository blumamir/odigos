package rollout

import (
	"context"
	"errors"
	"fmt"
	"time"

	argorolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	commonlogger "github.com/odigos-io/odigos/common/logger"
	"github.com/odigos-io/odigos/distros"
	sourceutils "github.com/odigos-io/odigos/k8sutils/pkg/source"
	"github.com/odigos-io/odigos/k8sutils/pkg/utils"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	openshiftappsv1 "github.com/openshift/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const RequeueWaitingForWorkloadRollout = 10 * time.Second

type RolloutResult struct {
	StatusChanged bool
	// Result contains the controller result for requeue behavior.
	Result ctrl.Result
}

// WorkloadKey generates a unique key for rate limiting purposes
func WorkloadKey(pw k8sconsts.PodWorkload) string {
	return fmt.Sprintf("%s/%s/%s", pw.Namespace, pw.Kind, pw.Name)
}

// Do potentially triggers a rollout for the given workload based on the given instrumentation config.
// If the instrumentation config is nil, the workload is rolled out - this is used for un-instrumenting workloads.
// Otherwise, the rollout config hash is calculated and compared to the saved hash in the instrumentation config.
// If the hashes are different, the workload is rolled out.
// If the hashes are the same, this is a no-op.
//
// If a rollout is triggered the status of the instrumentation config is updated with the new rollout hash
// and a corresponding condition is set.
//
// Returns a RolloutResult and an error.
func Do(ctx context.Context, c client.Client, ic *odigosv1alpha1.InstrumentationConfig, pw k8sconsts.PodWorkload, conf *common.OdigosConfiguration, distroProvider *distros.Provider, rolloutConcurrencyLimiter *RolloutConcurrencyLimiter) (RolloutResult, error) {
	isAutomaticRolloutDisabled, rolloutOptions := getRolloutOptions(conf)
	logger := commonlogger.FromContext(ctx)
	workloadObj := workload.ClientObjectFromWorkloadKind(pw.Kind)
	getErr := c.Get(ctx, client.ObjectKey{Name: pw.Name, Namespace: pw.Namespace}, workloadObj)
	if getErr != nil {
		return RolloutResult{}, client.IgnoreNotFound(getErr)
	}

	// Don't allow rollout of static pods, cronjobs or jobs
	if pw.Kind == k8sconsts.WorkloadKindStaticPod {
		if ic == nil {
			return RolloutResult{}, nil
		}

		changed := false
		if ic.Spec.PodManifestInjectionOptional {
			changed = meta.SetStatusCondition(&ic.Status.Conditions, conditionRestartNotRequiredForDistro)
		} else {
			changed = meta.SetStatusCondition(&ic.Status.Conditions, conditionStaticPodsNotSupported)
		}
		return RolloutResult{StatusChanged: changed}, nil
	}

	if pw.Kind == k8sconsts.WorkloadKindCronJob || pw.Kind == k8sconsts.WorkloadKindJob {
		if ic == nil {
			return RolloutResult{}, nil
		}
		changed := meta.SetStatusCondition(&ic.Status.Conditions, conditionWaitingForJobTrigger)
		return RolloutResult{StatusChanged: changed}, nil
	}

	if ic == nil {
		// If ic is nil and the PodWorkload is missing the odigos.io/agents-meta-hash label,
		// it means it is a rolled back application that shouldn't be rolled out again.
		hasAgents, agentsErr := workloadHasOdigosAgents(ctx, c, workloadObj)
		if agentsErr != nil {
			logger.Error(agentsErr, "failed to check for odigos agent labels")
			return RolloutResult{}, agentsErr
		}
		if !hasAgents {
			logger.Info("skipping rollout - workload already runs without odigos agents",
				"workload", pw.Name, "namespace", pw.Namespace)
			return RolloutResult{}, nil
		}

		// Just because an IC is nil, it doesn't mean the workload is not instrumented.
		// The workload may be instrumented by a source (and the IC may just temporarily be missing)
		// So we need to check if the workload is still marked for instrumentation.
		// For example, in a racey scenario where a workload is deleted and quickly replaced (such as with ArgoCD),
		// the new workload may exist before the IC is actually garbage collected by the deletion of the old workload.
		// Without this check, it would look like the IC was intentionally deleted (ie, via sourceinstrumentation controller).
		// This is a safety check: ic==nil is the signal, but the Source is the source of truth.
		stillInstrumented, instrumentedErr := workloadStillMarkedForInstrumentation(ctx, c, pw)
		if instrumentedErr != nil {
			logger.Error(instrumentedErr, "failed to check if workload is still marked for instrumentation")
			return RolloutResult{}, instrumentedErr
		}
		if stillInstrumented {
			logger.Info("skipping uninstrumentation rollout - workload is still covered by an active source",
				"workload", pw.Name, "namespace", pw.Namespace)
			return RolloutResult{}, nil
		}

		if isAutomaticRolloutDisabled {
			logger.Info("skipping rollout to uninstrument workload source - automatic rollout is disabled",
				"workload", pw.Name, "namespace", pw.Namespace)
			return RolloutResult{}, nil
		}

		// instrumentation config is deleted, trigger a rollout for the associated workload
		// this should happen once per workload, as the instrumentation config is deleted
		// and we want to rollout the workload to remove the instrumentation
		// Note: uninstrumentation rollouts are not rate limited since we can't track completion
		// (the IC is deleted so we won't get subsequent reconciles)
		logger.Debug("proceeding with uninstrumentation rollout",
			"workload", pw.Name,
			"namespace", pw.Namespace)
		rolloutConcurrencyLimiter.ReleaseWorkloadRolloutSlot(WorkloadKey(pw))
		rolloutErr := rolloutRestartWorkload(ctx, workloadObj, c, time.Now())
		return RolloutResult{}, client.IgnoreNotFound(rolloutErr)
	}

	if ic.Spec.PodManifestInjectionOptional {
		// all distributions used by this workload do not require a restart
		// thus, no rollout is needed
		rolloutConcurrencyLimiter.ReleaseWorkloadRolloutSlot(WorkloadKey(pw))
		changed := meta.SetStatusCondition(&ic.Status.Conditions, conditionRestartNotRequiredForDistro)
		return RolloutResult{StatusChanged: changed}, nil
	}

	if isAutomaticRolloutDisabled {
		// TODO: add more fine grained status conditions for this case.
		// For example: if the workload has already been rolled out, we can set the status to true
		// and signal that the process is considered completed.
		// If manual rollout is required, we can mention this for better UX.
		rolloutConcurrencyLimiter.ReleaseWorkloadRolloutSlot(WorkloadKey(pw))
		changed := meta.SetStatusCondition(&ic.Status.Conditions, conditionRolloutDisabled)
		return RolloutResult{StatusChanged: changed}, nil
	}

	workloadKey := WorkloadKey(pw)
	savedRolloutHash := ic.Status.WorkloadRolloutHash
	newRolloutHash := ic.Spec.AgentsMetaHash
	// Scenario: successful instrumentation ("X"="X") or successful uninstrumentation (""="")
	if savedRolloutHash == newRolloutHash {
		rolloutDone := utils.IsWorkloadRolloutDone(workloadObj)

		if !rolloutDone {
			return RolloutResult{Result: ctrl.Result{RequeueAfter: RequeueWaitingForWorkloadRollout}}, nil
		}

		// at this point, we know that rollout happened.
		// make sure to cleanup any stale condition from previous rollout attempts.
		// example:
		// - workload is instrumented successfully.
		// - rollout regardless of odigos
		// - agent disabled -> status is changed to "waiting for previous rollout to finish"
		// - agent enabled -> no need for a rollout, but the status needs to be updated.
		statusChanged := false
		// This is the happy flow - the workload is rolled out successfully
		if rolloutDone {
			statusChanged = meta.SetStatusCondition(&ic.Status.Conditions, conditionRolloutFinished)
			// Rollout is complete - release the slot if we had one
			rolloutConcurrencyLimiter.ReleaseWorkloadRolloutSlot(workloadKey)
		}
		return RolloutResult{StatusChanged: statusChanged, Result: ctrl.Result{}}, nil
	}

	// if a rollout is ongoing, wait for it to finish, requeue
	statusChanged := false
	if !utils.IsWorkloadRolloutDone(workloadObj) {
		statusChanged = meta.SetStatusCondition(&ic.Status.Conditions, conditionPreviousRolloutOngoing)
		return RolloutResult{StatusChanged: statusChanged, Result: ctrl.Result{RequeueAfter: RequeueWaitingForWorkloadRollout}}, nil
	}

	if !rolloutConcurrencyLimiter.TryAcquire(workloadKey, rolloutOptions.MaxConcurrentRollouts) {
		logger.Debug("rate limited instrumentation rollout, requeuing",
			"workload", pw.Name,
			"namespace", pw.Namespace,
			"requeueAfter", RequeueWaitingForWorkloadRollout)
		statusChanged = meta.SetStatusCondition(&ic.Status.Conditions, conditionWaitingInQueue)
		return RolloutResult{StatusChanged: statusChanged, Result: ctrl.Result{RequeueAfter: RequeueWaitingForWorkloadRollout}}, nil
	}

	// use the AgentsMetaHashChangedTime if it exists,
	// so we are idempotent if the reconciler requeue for any reason.
	var t time.Time
	if ic.Spec.AgentsMetaHashChangedTime != nil {
		t = ic.Spec.AgentsMetaHashChangedTime.Time
	} else {
		t = time.Now()
	}

	rolloutErr := rolloutRestartWorkload(ctx, workloadObj, c, t)
	if rolloutErr != nil {
		logger.Error(rolloutErr, "error rolling out workload", "name", pw.Name, "namespace", pw.Namespace)
	}

	ic.Status.WorkloadRolloutHash = newRolloutHash

	// If we have new rollout hash and also, AgentInjectionEnabled is enabled, that means we're instrumenting a new app
	if ic.Spec.AgentInjectionEnabled {
		now := metav1.NewTime(time.Now())
		ic.Status.InstrumentationTime = &now
	}
	// Setting the condition for successful triggering the rollout, or a failed to patch condition with a specific error message.
	meta.SetStatusCondition(&ic.Status.Conditions, rolloutCondition(rolloutErr))

	// at this point, the hashes are different, notify the caller the status has changed
	// Requeue to try and catch a crashing app
	return RolloutResult{StatusChanged: true, Result: ctrl.Result{RequeueAfter: RequeueWaitingForWorkloadRollout}}, nil
}

// RolloutRestartWorkload restarts the given workload by patching its template annotations.
// this is bases on the kubectl implementation of restarting a workload
// https://github.com/kubernetes/kubectl/blob/master/pkg/polymorphichelpers/objectrestarter.go#L32
func rolloutRestartWorkload(ctx context.Context, workloadObj client.Object, c client.Client, ts time.Time) error {

	logger := commonlogger.FromContext(ctx)

	patch := []byte(fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, ts.Format(time.RFC3339)))
	switch obj := workloadObj.(type) {
	case *appsv1.Deployment:
		if obj.Spec.Paused {
			return errors.New("can't restart paused deployment")
		}
		logger.Info("auto rollout - restarting deployment", "name", obj.Name, "namespace", obj.Namespace)
		return c.Patch(ctx, obj, client.RawPatch(types.MergePatchType, patch))
	case *appsv1.StatefulSet:
		logger.Info("auto rollout - restarting stateful set", "name", obj.Name, "namespace", obj.Namespace)
		return c.Patch(ctx, obj, client.RawPatch(types.MergePatchType, patch))
	case *appsv1.DaemonSet:
		logger.Info("auto rollout - restarting daemon set", "name", obj.Name, "namespace", obj.Namespace)
		return c.Patch(ctx, obj, client.RawPatch(types.MergePatchType, patch))
	case *openshiftappsv1.DeploymentConfig:
		logger.Info("auto rollout - restarting deployment config", "name", obj.Name, "namespace", obj.Namespace)
		return c.Patch(ctx, obj, client.RawPatch(types.MergePatchType, patch))
	case *argorolloutsv1alpha1.Rollout:
		logger.Info("auto rollout - restarting argo rollout", "name", obj.Name, "namespace", obj.Namespace)
		// Rollouts use a different field (spec.restartAt) for restarting, so we need to patch it differently
		// https://github.com/argoproj/argo-rollouts/blob/cb1c33df7a2c2b1c2ed31b1ee0aa22621ef5577c/utils/replicaset/replicaset.go#L223-L232
		rolloutPatch := []byte(fmt.Sprintf(`{"spec":{"restartAt":"%s"}}`, ts.Format(time.RFC3339)))
		return c.Patch(ctx, obj, client.RawPatch(types.MergePatchType, rolloutPatch))
	case *corev1.Pod:
		if workload.IsStaticPod(obj) {
			return errors.New("can't restart static pods")
		}
		logger.Info("auto rollout - restarting standalone pod", "name", obj.Name, "namespace", obj.Namespace)
		return errors.New("currently not supporting standalone pods as workloads for rollout")
	default:
		return errors.New("unknown kind")
	}
}

func rolloutCondition(rolloutErr error) metav1.Condition {
	if rolloutErr == nil {
		return conditionTriggeredSuccessfully
	}
	return newConditionFailedToPatch(rolloutErr)
}

// workloadLabelSelector returns the label selector for a workload object
func workloadLabelSelector(obj client.Object) (*metav1.LabelSelector, error) {
	var selector *metav1.LabelSelector

	switch o := obj.(type) {
	case *appsv1.Deployment:
		selector = o.Spec.Selector
	case *appsv1.StatefulSet:
		selector = o.Spec.Selector
	case *appsv1.DaemonSet:
		selector = o.Spec.Selector
	case *openshiftappsv1.DeploymentConfig:
		// DeploymentConfig selector is map[string]string, convert to *metav1.LabelSelector
		selector = &metav1.LabelSelector{
			MatchLabels: o.Spec.Selector,
		}
	case *argorolloutsv1alpha1.Rollout:
		selector = &metav1.LabelSelector{
			MatchLabels: o.Spec.Selector.MatchLabels,
		}
	default:
		return nil, fmt.Errorf("workloadLabelSelector: unsupported workload kind %T", obj)
	}

	if selector == nil {
		return nil, fmt.Errorf("workloadLabelSelector: workload has nil selector")
	}

	return selector, nil
}

// instrumentedPodsSelector returns a selector for all the instrumented pods that are associated with the workload object
func instrumentedPodsSelector(obj client.Object) (labels.Selector, error) {
	labelSelector, err := workloadLabelSelector(obj)
	if err != nil {
		return nil, err
	}

	// Create a deep copy of the selector to avoid mutating the original
	selectorCopy := labelSelector.DeepCopy()
	selectorCopy.MatchExpressions = append(selectorCopy.MatchExpressions, metav1.LabelSelectorRequirement{
		Key:      k8sconsts.OdigosAgentsMetaHashLabel,
		Operator: metav1.LabelSelectorOpExists,
	})

	sel, err := metav1.LabelSelectorAsSelector(selectorCopy)
	if err != nil {
		return nil, fmt.Errorf("instrumentedPodsSelector: invalid selector: %w", err)
	}

	return sel, nil
}

// workloadHasOdigosAgents returns true if the workload still has *any* pod present in the instrumented-pod.
func workloadHasOdigosAgents(ctx context.Context, c client.Client, obj client.Object) (bool, error) {
	sel, err := instrumentedPodsSelector(obj)
	if err != nil {
		return false, fmt.Errorf("workloadHasOdigosAgents: invalid selector: %w", err)
	}

	var pods corev1.PodList
	if err := c.List(
		ctx, &pods,
		client.InNamespace(obj.GetNamespace()),
		client.MatchingLabelsSelector{Selector: sel},
	); err != nil {
		return false, fmt.Errorf("workloadHasOdigosAgents: listing pods failed: %w", err)
	}

	return len(pods.Items) > 0, nil
}

func workloadStillMarkedForInstrumentation(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) (bool, error) {
	sources, err := odigosv1alpha1.GetSources(ctx, c, pw)
	if err != nil {
		return false, err
	}
	enabled, _, err := sourceutils.IsObjectInstrumentedBySource(ctx, sources, err)
	if err != nil {
		return false, err
	}
	return enabled, nil
}
