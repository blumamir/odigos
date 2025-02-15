package updatedpods

import (
	"context"
	"fmt"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodInstrumentationStatus struct {

	// the hash of the agents meta data or empty string if not found
	// this is extracted from the pod label which the webhook sets
	PodAgentMetaHash string
}

func isContainerAgentEnabled(icContainers []odigosv1.ContainerAgentConfig, conatinerName string) bool {
	for i := range icContainers {
		if icContainers[i].ContainerName == conatinerName {
			return true
		}
	}
	return false
}

// func calculateContainersHealth(podsContainers []corev1.ContainerStatus, icContainers []odigosv1.ContainerAgentConfig) (containersHealth int, totalContainers int) {

// 	odigosReadyContainers := 0
// 	odigosStartedContainers := 0

// 	for i := range podsContainers {
// 		containerStatus := &podsContainers[i]
// 		if isContainerAgentEnabled(icContainers, containerStatus.Name) {
// 			if containerStatus.Ready {
// 				odigosReadyContainers++
// 			} else if containerStatus.Started != nil && *containerStatus.Started {
// 				// if container is not ready, but has been started, we consider it as healthy
// 				odigosStartedContainers++
// 			} else {

// 			}
// 		}
// 	}
// }

func calculateInstrumentationStatusForPod(pod *corev1.Pod, ic *odigosv1.InstrumentationConfig) *PodInstrumentationStatus {

	// if agents hash not found on the pod, it will be empty string to denote that
	actualAgentsMetaHash := pod.Labels[k8sconsts.OdigosAgentsMetaHashLabel]
	//the desired hash will be empty if no agents should be injected to this pod
	desiredAgentsMetaHash := ic.Spec.AgentsMetaHash

	runningCorrectAgentsHash := actualAgentsMetaHash == desiredAgentsMetaHash

	if runningCorrectAgentsHash && ic.Spec.AgentInjectionEnabled {
		// calculate containers health for instrumented containers, e.g. only if agents hash is correct and not empty
		// containersHealth, totalContainers := calculateContainersHealth(pod.Status.ContainerStatuses, ic.Spec.Containers)
	}

	return &PodInstrumentationStatus{
		PodAgentMetaHash: actualAgentsMetaHash,
	}
}

func reconcileWorkload(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) (ctrl.Result, error) {

	icName := workload.InstrumentationConfigNameFromPodWorkload(pw)

	var ic odigosv1.InstrumentationConfig
	err := c.Get(ctx, client.ObjectKey{Name: icName, Namespace: pw.Namespace}, &ic)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err) // if ic not found, we have nothing to update
	}

	workloadObj := workload.ClientObjectFromWorkloadKind(pw.Kind)
	err = c.Get(ctx, client.ObjectKey{Name: pw.Name, Namespace: pw.Namespace}, workloadObj)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wl, err := workload.ObjectToWorkload(workloadObj)
	if err != nil {
		return ctrl.Result{}, reconcile.TerminalError(err)
	}

	selector := wl.GetPodsSelector()
	selectedPods := &corev1.PodList{}
	err = c.List(ctx, selectedPods, client.MatchingLabels(selector.MatchLabels), client.InNamespace(pw.Namespace))
	if err != nil {
		return ctrl.Result{}, err
	}

	desiredAgentsMetaHash := ic.Spec.AgentsMetaHash

	updatedReadyPods := 0
	updatedStartingPods := 0
	instrumentedNotUpToDatePods := 0

	for i := range selectedPods.Items {
		pod := &selectedPods.Items[i]
		ps := calculateInstrumentationStatusForPod(pod, &ic)
		if ps.PodAgentMetaHash == desiredAgentsMetaHash {
			updatedReadyPods++
		} else {
			instrumentedNotUpToDatePods++
		}
	}

	reason := "FooBar"
	message := "foo bar"
	if instrumentedNotUpToDatePods == 0 {
		reason = string(odigosv1.UpdatedPodsReasonPodsUpToDateReady)
		message = fmt.Sprintf("All %d pods are up to date and ready", updatedReadyPods)
	}

	cond := metav1.Condition{
		Type:    odigosv1.UpdatedPodsStatusConditionType,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	}
	changed := meta.SetStatusCondition(&ic.Status.Conditions, cond)
	if changed {
		err = c.Status().Update(ctx, &ic)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
