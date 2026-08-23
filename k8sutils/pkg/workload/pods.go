package workload

import (
	"context"
	"fmt"
	"strings"

	"github.com/odigos-io/odigos/api/k8sconsts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListNonCompletedPods returns pods belonging to the given workload,
// excluding pods in Succeeded or Failed phase.
func ListNonCompletedPods(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) ([]corev1.Pod, error) {
	workloadObj := ClientObjectFromWorkloadKind(pw.Kind)
	if workloadObj == nil {
		return nil, fmt.Errorf("unsupported workload kind %q", pw.Kind)
	}
	err := c.Get(ctx, types.NamespacedName{Namespace: pw.Namespace, Name: pw.Name}, workloadObj)
	if err != nil {
		return nil, err
	}

	var pods []corev1.Pod
	switch pw.Kind {
	case k8sconsts.WorkloadKindStaticPod:
		// Static pods are the workload themselves and have no label selector.
		pods = []corev1.Pod{*workloadObj.(*corev1.Pod)}
	case k8sconsts.WorkloadKindCronJob:
		// CronJobs have no label selector. Their pods are owned by Jobs named
		// <cronjob-name>-<timestamp>; resolve ownership the same way as ownerreference.go.
		podList := &corev1.PodList{}
		err = c.List(ctx, podList, client.InNamespace(pw.Namespace))
		if err != nil {
			return nil, err
		}
		for i := range podList.Items {
			pod := &podList.Items[i]
			for _, owner := range pod.OwnerReferences {
				if owner.Kind != string(k8sconsts.WorkloadKindJob) {
					continue
				}
				workloadName, workloadKind, err := GetWorkloadNameAndKind(owner.Name, owner.Kind, pod)
				if err != nil {
					continue
				}
				if workloadKind == k8sconsts.WorkloadKindCronJob && workloadName == pw.Name {
					pods = append(pods, *pod)
					break
				}
			}
		}
	default:
		genericWorkload, err := ObjectToWorkload(workloadObj)
		if err != nil {
			return nil, err
		}

		labelSelector := genericWorkload.LabelSelector()
		if labelSelector == nil {
			return nil, fmt.Errorf("unexpected nil label selector for workload %s/%s kind %s", pw.Namespace, pw.Name, pw.Kind)
		}

		podList := &corev1.PodList{}
		err = c.List(ctx, podList, client.InNamespace(pw.Namespace), client.MatchingLabels(labelSelector.MatchLabels))
		if err != nil {
			return nil, err
		}
		pods = podList.Items
	}

	nonCompleted := make([]corev1.Pod, 0, len(pods))
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		nonCompleted = append(nonCompleted, pod)
	}
	return nonCompleted, nil
}

// IsStaticPod return true whether the pod is static or not
// https://kubernetes.io/docs/tasks/configure-pod-container/static-pod/
func IsStaticPod(p *corev1.Pod) bool {
	var nodeOwner *metav1.OwnerReference
	for _, owner := range p.OwnerReferences {
		if owner.Kind == "Node" {
			nodeOwner = &owner
			break
		}
	}

	// static pods are owned by nodes
	if nodeOwner == nil {
		return false
	}

	// https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetes-io-config-source
	// This annotation is added by the kubelet to indicate where the Pod comes from.
	// For static Pods, the annotation value could be one of file or http depending on where the Pod manifest is located.
	// For a Pod created on the API server and then scheduled to the current node, the annotation value is api.
	if p.Annotations == nil {
		return false
	}
	configSource, ok := p.Annotations["kubernetes.io/config.source"]
	if !ok {
		return false
	}
	return configSource == "file" || configSource == "http"
}

func PodUID(p *corev1.Pod) string {
	if IsStaticPod(p) {
		// https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetes-io-config-hash
		// When the kubelet creates a static Pod based on a given manifest,
		// it attaches this annotation to the static Pod. The value of the annotation is the UID of the Pod
		return p.Annotations["kubernetes.io/config.hash"]
	}

	return string(p.UID)
}

// StaticPodName returns the value of the static pod name without its node name
// since static pods name are of the form
// <static-pod-name>-<node-name>
// if the pod is not static, or its name does not match the expected pattern,
// an empty string will be returned
func StaticPodName(p *corev1.Pod) string {
	if p == nil {
		return ""
	}

	if !IsStaticPod(p) {
		return ""
	}

	nodeName := p.Spec.NodeName
	if nodeName == "" {
		return ""
	}

	staticPodName, found := strings.CutSuffix(p.Name, "-"+nodeName)
	if !found {
		return ""
	}

	return staticPodName
}
